package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/external/hederaUtils"
	utils "github.com/Powerdfi-com/Backend/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/labstack/echo/v4"
)

type portfolioHandler struct {
	app internal.Application
}

func NewPortfolioHandler(app internal.Application) portfolioHandler {
	return portfolioHandler{app: app}
}

func (h portfolioHandler) Update(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	req := struct {
		ListingId string `json:"listingId" validate:"required"`
		Quantity  int64  `json:"quantity" validate:"required,min=1,max=100000"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	listing, err := h.app.Repositories.Listing.GetById(req.ListingId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset listing not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if req.Quantity > listing.Quantity {
		return echo.NewHTTPError(http.StatusBadRequest, "quantity exceeds listing quantity")
	}

	asset, err := h.app.Repositories.Asset.GetById(listing.AssetId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {

			return echo.NewHTTPError(http.StatusNotFound, "asset not found")
		}

		return helpers.ErrInternalServer(ctx, err)
	}

	user, err := h.app.Repositories.User.GetByAddress(authenticatedUser.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}
	creator, err := h.app.Repositories.User.GetById(listing.UserId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	data, err := utils.DecryptPlainText(user.EncryptedPrivateKey, h.app.Config.MasterKey)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	key, err := hedera.PrivateKeyFromBytesECDSA(data)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	cdata, err := utils.DecryptPlainText(creator.EncryptedPrivateKey, h.app.Config.MasterKey)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	ckey, err := hedera.PrivateKeyFromBytesECDSA(cdata)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	accountID, err := hedera.AccountIDFromString(user.AccountID)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	creatorID, err := hedera.AccountIDFromString(creator.AccountID)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	assetTokenId, err := hedera.TokenIDFromString(asset.TokenId)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	usdcTokenId, err := hedera.TokenIDFromString(h.app.Config.Hedera.TokenIdUSDC)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	account := hederaUtils.InitializeAccount(h.app.HederaClient.Client, &accountID)
	account.SetPrivateKey(&key)

	balance, err := account.GetTokenBalanceWithID(usdcTokenId)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	parsedBalance := hedera.HbarFrom(float64(balance), hedera.HbarUnits.Microbar)
	if parsedBalance.As(hedera.HbarUnits.Hbar) < listing.PriceUSD {
		return echo.NewHTTPError(http.StatusBadRequest, "insufficient balance")
	}

	creatorAccount := hederaUtils.InitializeAccount(h.app.HederaClient.Client, &creatorID)
	creatorAccount.SetPrivateKey(&ckey)

	assetOwnershipCreator, err := h.app.Repositories.AssetOwner.GetOwnerAsset(
		asset.Id,
		asset.CreatorUserID,
	)

	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	assetOwnershipUser, err := h.app.Repositories.AssetOwner.GetOwnerAsset(
		asset.Id,
		user.Id,
	)
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if len(assetOwnershipCreator.SerialNumbers) < int(req.Quantity) {
		return helpers.ErrInternalServer(ctx, fmt.Errorf("quantitiy exceeds creator's balance"))
	}

	// return ctx.JSON(http.StatusOK, "igwe wake up")
	_, err = account.TransferFungibleToken(
		uint64(hedera.HbarFrom(listing.PriceUSD, hedera.HbarUnits.Hbar).As(hedera.HbarUnits.Microbar)),
		&creatorID,
		&usdcTokenId)

	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	_, err = account.AssociateAccountToToken(assetTokenId)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)

	}

	for range req.Quantity {
		serialNumber := assetOwnershipCreator.SerialNumbers[0]
		_, _, err = creatorAccount.TransferNonFungibleToken(&accountID, &hedera.NftID{
			TokenID:      assetTokenId,
			SerialNumber: serialNumber,
		})
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}

		assetOwnershipUser.SerialNumbers = append(assetOwnershipUser.SerialNumbers, serialNumber)

		assetOwnershipCreator.SerialNumbers = assetOwnershipCreator.SerialNumbers[1:]

		err = h.app.Repositories.AssetOwner.Update(assetOwnershipCreator)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}

		err = h.app.Repositories.AssetOwner.Update(assetOwnershipUser)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, "success")
}

func (h portfolioHandler) Get(ctx echo.Context) error {
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
