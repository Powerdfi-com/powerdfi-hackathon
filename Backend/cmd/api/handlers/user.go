package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/external/hederaUtils"
	"github.com/Powerdfi-com/Backend/external/shufti"
	utils "github.com/Powerdfi-com/Backend/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	app internal.Application
}

func NewUserHandler(app internal.Application) userHandler {
	return userHandler{app: app}
}

func (h userHandler) ActivateUser(ctx echo.Context) error {
	reqBody := struct {
		UserName string `json:"userName" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
	}{}

	if err := ctx.Bind(&reqBody); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	fetchedUser, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	// prevent multiple activations of the same user
	if fetchedUser.IsActive {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "user is already activated",
		}
	}

	account, err := hederaUtils.CreateNewAccount(h.app.HederaClient.Client)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	tokenId, err := hedera.TokenIDFromString(h.app.Config.Hedera.TokenIdUSDC)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	_, err = account.AssociateAccountToToken(tokenId)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	fetchedUser.AccountID = account.GetAccountID().String()

	encryptedPrivateKey, err := utils.EncryptPlainText(account.GetAccountPrivateKey().BytesDer(), h.app.Config.MasterKey)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	fetchedUser.Username = &reqBody.UserName

	fetchedUser.Email = &reqBody.Email

	err = h.app.Repositories.User.SetPrivateKey(authenticatedUser.Id, encryptedPrivateKey)

	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	// TODO: check for duplicate email to avoid waste of hedera resources
	updatedUser, err := h.app.Repositories.User.Update(fetchedUser)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateDetails):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "email or username already exists",
			}

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	// activate user
	err = h.app.Repositories.User.Activate(authenticatedUser.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	updatedUser.IsActive = true
	updatedUser.EncryptedPrivateKey = encryptedPrivateKey
	// generate new access and refresh tokens for updated user with new activation status
	accessToken, refreshToken, expiration, err := helpers.GenerateTokens(h.app, updatedUser)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := map[string]interface{}{
		"tokens": response.TokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiration,
		},
		"user": response.UserResponseFromModel(updatedUser),
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) Get(ctx echo.Context) error {

	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed user id")
	}
	user, err := h.app.Repositories.User.GetById(parsedId.String())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	resp := response.UserResponseFromModel(user)

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetOwnProfile(ctx echo.Context) error {

	authenticatedUser := helpers.ContextGetUser(ctx)

	user, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	resp := response.UserResponseFromModel(user)

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) UpdateProfile(ctx echo.Context) error {
	// request binding
	req := struct {
		Username string `json:"username" `
		Avatar   string `json:"avatar"`
		Bio      string `json:"bio"`
		Website  string `json:"website"`
		Twitter  string `json:"twitter"`
		Discord  string `json:"discord"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	user, err := h.app.Repositories.User.GetById(authenticatedUser.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	user.Bio = req.Bio

	user.Website = req.Website
	user.Twitter = req.Twitter
	user.Discord = req.Discord

	if strings.TrimSpace(req.Avatar) != "" {
		_, err := url.ParseRequestURI(req.Avatar)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"avatar must be a valid url",
			)
		}
		user.Avatar = req.Avatar
	}

	// set empty username and email to nil
	if strings.TrimSpace(req.Username) == "" {
		user.Username = nil
	} else {
		// validate provided username
		if username := req.Username; common.IsHexAddress(username) || utils.IsValidUuid(username) {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"username cannot be an address or a UUID",
			)
		}

		user.Username = &req.Username
	}
	// if strings.TrimSpace(req.Email) == "" {
	// 	user.Email = nil
	// } else {
	// 	// validate provided email
	// 	if !utils.IsValidEmail(req.Email) {
	// 		return echo.NewHTTPError(
	// 			http.StatusBadRequest,
	// 			"invalid email address",
	// 		)
	// 	}

	// 	user.Email = &req.Email
	// }

	// update user details
	updatedUser, err := h.app.Repositories.User.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateDetails):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "username or email already exists",
			}

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	user.UpdatedAt = updatedUser.UpdatedAt

	return ctx.JSON(http.StatusOK, response.UserResponseFromModel(user))
}

func (h userHandler) GetWalletDeets(ctx echo.Context) error {

	authenticatedUser := helpers.ContextGetUser(ctx)

	user, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	accountId, err := hedera.AccountIDFromString(user.AccountID)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	tokenId, err := hedera.TokenIDFromString(h.app.Config.Hedera.TokenIdUSDC)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	query := hedera.NewAccountBalanceQuery().SetAccountID(accountId)

	balanceCheck, err := query.Execute(h.app.HederaClient.Client)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	rawBalance := balanceCheck.Tokens.Get(tokenId)

	parsedBalance := hedera.HbarFrom(float64(rawBalance), hedera.HbarUnits.Microbar)
	// parsedBalance := hedera.HbarFromTinybar(int64(rawBalance))

	resp := response.WalletResponse{
		AccountId: user.AccountID,
		Address:   accountId.ToSolidityAddress(),
		Balance:   parsedBalance.As(hedera.HbarUnits.Hbar),
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetCreatedAssets(ctx echo.Context) error {
	req := struct {
		UserId   string `param:"userId" validate:"required,uuid"`
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
		Search   string `query:"search"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:   req.Page,
		Limit:  req.PageSize,
		Search: req.Search,
	}

	// check if given address belongs to a registered user
	user, err := h.app.Repositories.User.GetById(req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	assets, err := h.app.Repositories.User.ListCreatedAssets(user.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	resp := []response.AssetResponse{}
	for _, asset := range assets {
		isListedByUser, err := h.app.Repositories.Asset.IsListed(authenticatedUser.Id, asset.Id)
		if err != nil {
			ctx.Logger().Errorf("err: getting listed status%w", err)
		}
		asset.IsListedByUser = isListedByUser
		resp = append(resp, response.AssetResponseFromModel(asset))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetKycLink(ctx echo.Context) error {

	authenticatedUser := helpers.ContextGetUser(ctx)

	user, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if user.IsVerified {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "user is already verified",
		}
	}

	userKyc, err := h.app.Repositories.UserKyc.GetByUserId(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			referenceId := fmt.Sprintf("%s%d", user.Address, time.Now().Unix())
			var email string

			if user.Email != nil {
				email = *user.Email
			}
			shuftikycResp, err := h.app.ShuftiClient.GetJourneyVerificationLink(referenceId, email)

			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}

			userKyc = models.UserKyc{
				UserId:      user.Id,
				URL:         shuftikycResp.VerificationURL,
				Platform:    internal.KYC_PLATFORM,
				ReferenceId: referenceId,
			}

			err = h.app.Repositories.UserKyc.Create(userKyc)
			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}
		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	} else {
		resp, err := h.app.ShuftiClient.GetKycStatus(userKyc.ReferenceId)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}

		switch resp.Event {
		case shufti.EVENT_TYPE_INVALID, shufti.EVENT_TYPE_CANCELLED, shufti.EVENT_TYPE_TIMEOUT:
			fmt.Println("renewing")
			// err = h.app.Repositories.UserKyc.Delete(user.Id)
			// if err != nil {
			// 	if !errors.Is(err, repository.ErrRecordNotFound) {
			// 		return helpers.ErrInternalServer(ctx, err)
			// 	}

			// }

			referenceId := fmt.Sprintf("%s%d", user.Address, time.Now().Unix())
			var email string

			if user.Email != nil {
				email = *user.Email
			}
			shuftikycResp, err := h.app.ShuftiClient.GetJourneyVerificationLink(referenceId, email)

			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}

			userKyc = models.UserKyc{
				UserId:      user.Id,
				URL:         shuftikycResp.VerificationURL,
				Platform:    internal.KYC_PLATFORM,
				ReferenceId: referenceId,
			}

			err = h.app.Repositories.UserKyc.Update(userKyc)
			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}

		}

	}

	resp := map[string]interface{}{
		"link": userKyc.URL,
	}
	return ctx.JSON(http.StatusOK, resp)

}

func (h userHandler) InitiateKYB(ctx echo.Context) error {

	reqBody := struct {
		CompanyName      string `json:"companyName" validate:"required"`
		CompanyRegNumber string `json:"companyRegNumber" validate:"required"`
		CompanyLocation  string `json:"companyLocation" validate:"required"`
		CompanyAddress   string `json:"companyAddress" validate:"required"`
		CertificateOfInc string `json:"certificateOfInc" validate:"required"`
	}{}

	if err := ctx.Bind(&reqBody); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	user, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if user.IsVerified {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "user is already verified",
		}
	}

	_, err = h.app.Repositories.UserKYB.GetByUserId(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			referenceId := fmt.Sprintf("%s%d", user.Address, time.Now().Unix())

			err := h.app.ShuftiClient.InitiateKYBVerification(referenceId, reqBody.CompanyRegNumber, reqBody.CompanyLocation)

			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}

			userKyb := models.UserKyB{
				UserId:           user.Id,
				CertificateOfInc: reqBody.CertificateOfInc,
				CompanyName:      reqBody.CompanyName,
				CompanyLocation:  reqBody.CompanyLocation,
				CompanyAddress:   reqBody.CompanyAddress,
				CompanyRegNo:     reqBody.CompanyRegNumber,
				Status:           models.KYB_STATUS_PENDING,
				Platform:         internal.KYC_PLATFORM,
				ReferenceId:      referenceId,
			}

			err = h.app.Repositories.UserKYB.Create(userKyb)
			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}
		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})

}
func (h userHandler) GetKYBStatus(ctx echo.Context) error {

	authenticatedUser := helpers.ContextGetUser(ctx)

	user, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	var status string

	_, err = h.app.Repositories.UserKYB.GetByUserId(user.Id)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			status = "unverified"
		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if user.IsVerified {
		status = models.KYC_STATUS_SUCCESS
	} else {
		status = models.KYB_STATUS_PENDING
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":     status,
		"isVerified": user.IsVerified,
	})

}

func (h userHandler) GetKYBCountries(ctx echo.Context) error {

	countries, err := h.app.ShuftiClient.GetCountries()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	resp := []shufti.Country{}
	for _, country := range countries {
		catResp := shufti.Country{
			Id:   country.Id,
			Name: country.Name,
		}
		resp = append(resp, catResp)
	}

	return ctx.JSON(http.StatusOK, resp)

}

func (h userHandler) GetListedItems(ctx echo.Context) error {
	req := struct {
		UserId   string `param:"userId" validate:"required,uuid"`
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
		Search   string `query:"search"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 || req.PageSize > models.MaxPageLimit {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:   req.Page,
		Limit:  req.PageSize,
		Search: req.Search,
	}

	// check if given address belongs to a registered user
	user, err := h.app.Repositories.User.GetById(req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	items, err := h.app.Repositories.User.GetListedAssets(user.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.AssetResponse{}
	for _, item := range items {
		resp = append(resp, response.AssetResponseFromModel(item))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetUnListedItems(ctx echo.Context) error {
	req := struct {
		UserId   string `param:"userId" validate:"required,uuid"`
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
		Search   string `query:"search"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 || req.PageSize > models.MaxPageLimit {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:   req.Page,
		Limit:  req.PageSize,
		Search: req.Search,
	}

	// check if given address belongs to a registered user
	user, err := h.app.Repositories.User.GetById(req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	items, err := h.app.Repositories.User.GetUnListedAssets(user.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.AssetResponse{}
	for _, item := range items {
		resp = append(resp, response.AssetResponseFromModel(item))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetPortfolio(ctx echo.Context) error {
	req := struct {
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
		Search   string `query:"search"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:   req.Page,
		Limit:  req.PageSize,
		Search: req.Search,
	}

	// check if given address belongs to a registered user
	authenticatedUser := helpers.ContextGetUser(ctx)

	assets, err := h.app.Repositories.User.ListOwnedAssets(authenticatedUser.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.AssetResponse{}
	for _, asset := range assets {

		resp = append(resp, response.AssetResponseFromModel(asset))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetUserOrders(ctx echo.Context) error {
	req := struct {
		Page     int                 `query:"page"`
		PageSize int                 `query:"size"`
		Status   *models.OrderStatus `query:"status"`
		Type     *models.OrderType   `query:"type"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 || req.PageSize > models.MaxPageLimit {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
	}
	authenticatedUser := helpers.ContextGetUser(ctx)

	orders, err := h.app.Repositories.User.GetOrders(authenticatedUser.Id, req.Status, req.Type, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.OrderResponse{}
	for _, order := range orders {

		resp = append(resp, response.OrderResponseFromModel(order))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h userHandler) GetUserActivities(ctx echo.Context) error {
	req := struct {
		Page     int `query:"page"`
		PageSize int `query:"size"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 || req.PageSize > models.MaxPageLimit {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
	}
	authenticatedUser := helpers.ContextGetUser(ctx)

	activities, totalCount, err := h.app.Repositories.Activity.ListForUser(authenticatedUser.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.ActivityResponse{}
	for _, activity := range activities {

		resp = append(resp, response.ActivityResponseFromModel(activity))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":      totalCount,
		"activities": resp,
	})

}

func (h userHandler) TopAssetsPerfs(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.PageLimit9
	}
	if req.Range == 0 {
		req.Range = models.Range30
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
		Range: req.Range,
	}
	authenticatedUser := helpers.ContextGetUser(ctx)

	topAssets, totalCount, err := h.app.Repositories.Stats.TopUserAssetPerfs(authenticatedUser.Id, filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	topAssetsResponse := []response.AssetStatsResponse{}
	for _, topAsset := range topAssets {
		logo := internal.ChainLogoMap[models.Chain(topAsset.Blockchain)]
		topAsset.BlockchainLogo = logo
		priceChanges, err := h.app.Repositories.Stats.GetAssetPriceData(topAsset.AssetId, filter)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		topAsset.PriceChanges = priceChanges
		topAssetsResponse = append(topAssetsResponse, response.AssetStatsResponseFromModel(topAsset))
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": topAssetsResponse,
	})
}

func (h userHandler) TrendingAssetsPerfs(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.PageLimit9
	}
	if req.Range == 0 {
		req.Range = models.Range30
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
		Range: req.Range,
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	trendingAssets, totalCount, err := h.app.Repositories.Stats.TrendingUserAssetPerfs(authenticatedUser.Id, filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	trendingAssetsResponse := []response.AssetStatsResponse{}

	for _, trendingAsset := range trendingAssets {
		logo := internal.ChainLogoMap[models.Chain(trendingAsset.Blockchain)]
		trendingAsset.BlockchainLogo = logo
		priceChanges, err := h.app.Repositories.Stats.GetAssetPriceData(trendingAsset.AssetId, filter)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		trendingAsset.PriceChanges = priceChanges
		trendingAssetsResponse = append(trendingAssetsResponse, response.AssetStatsResponseFromModel(trendingAsset))
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": trendingAssetsResponse,
	})
}

func (h userHandler) SalesTrend(ctx echo.Context) error {
	req := struct {
		Range uint `query:"range"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if req.Range == 0 {
		req.Range = models.Range30
	}

	filter := models.Filter{
		Range: req.Range,
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	saleTrends, err := h.app.Repositories.Stats.SalesTrend(authenticatedUser.Id, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	saleTrendResponse := []response.SalesTrendResponse{}

	for _, saleTrend := range saleTrends {

		saleTrendResponse = append(saleTrendResponse, response.SalesTrendResponse{
			Period: saleTrend.Period,
		})
	}
	return ctx.JSON(http.StatusOK, saleTrendResponse)
}
