package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/external/hederaUtils"
	utils "github.com/Powerdfi-com/Backend/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/labstack/echo/v4"
)

type assetHandler struct {
	app internal.Application
}

func NewAssetsHandler(app internal.Application) assetHandler {
	return assetHandler{app: app}
}

func (h assetHandler) Create(ctx echo.Context) error {
	u := helpers.ContextGetUser(ctx)

	req := struct {
		Name         string   `json:"name" validate:"min=3"`
		Symbol       string   `json:"symbol" validate:"min=3"`
		BlockchainId string   `json:"blockchainId"`
		Description  string   `json:"description"`
		Properties   string   `json:"properties"`
		CategoryId   *int     `json:"categoryId"`
		URLs         []string `json:"urls"`

		LegalDocumentURLs    []string `json:"legalDocs"`
		IssuanceDocumentURLs []string `json:"issuanceDocs"`
		Supply               int64    `json:"totalSupply" validate:"min=1,max=100000"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Properties == "" {
		req.Properties = "[]"
	}

	_, ok := internal.ChainNameMap[models.Chain(req.BlockchainId)]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid blockchain")
	}
	if req.CategoryId != nil {
		_, err := h.app.Repositories.Category.GetById(*req.CategoryId)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryId")

			default:
				log.Printf("error: category getById %v", err)
				return helpers.ErrInternalServer(ctx, err)
			}

		}
	}
	properties := models.AssetProperties{}
	err := json.NewDecoder(strings.NewReader(req.Properties)).Decode(&properties)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"malformed asset property; it must be a valid JSON list of pairs",
		)
	}

	user, err := h.app.Repositories.User.GetById(u.Id)
	if err != nil {
		log.Printf("error: get userId %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	// if !user.IsVerified {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "kyc not fufilled")
	// }

	nftCreate, err := hedera.NewTokenCreateTransaction().
		SetTokenName(req.Name).
		SetTokenSymbol(req.Symbol).
		SetTokenType(hedera.TokenTypeNonFungibleUnique).
		SetDecimals(0).
		SetInitialSupply(0).
		SetTreasuryAccountID(*h.app.HederaClient.AccountID).
		SetSupplyType(hedera.TokenSupplyTypeFinite).
		SetMaxSupply(req.Supply).
		SetSupplyKey(*h.app.HederaClient.PrivateKey).
		FreezeWith(h.app.HederaClient.Client)

	if err != nil {
		log.Printf("error: tokenCreate Txn %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	nftCreateTxSign := nftCreate.Sign(*h.app.HederaClient.PrivateKey)

	// Submit the transaction to a Hedera network
	nftCreateSubmit, err := nftCreateTxSign.Execute(h.app.HederaClient.Client)
	if err != nil {
		log.Printf("error: create TxSign %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	// Get the transaction receipt
	nftCreateRx, err := nftCreateSubmit.GetReceipt(h.app.HederaClient.Client)
	if err != nil {
		log.Printf("error: get recpt %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	// var serialNumber int64 = nftCreateRx.

	asset := models.Asset{
		// set item ID here in order to use it for immutable metadata and image URLs
		Id:            uuid.NewString(),
		Name:          req.Name,
		Symbol:        req.Symbol,
		CategoryId:    req.CategoryId,
		Description:   req.Description,
		TotalSupply:   req.Supply,
		Properties:    properties,
		CreatorUserID: user.Id,
		TokenId:       nftCreateRx.TokenID.String(),
		Blockchain:    models.Chain(req.BlockchainId),
		// MetadataUrl:req.MetadataUrl,

		URLs:                 req.URLs,
		LegalDocumentURLs:    req.LegalDocumentURLs,
		IssuanceDocumentURLs: req.IssuanceDocumentURLs,
	}

	createdAsset, err := h.app.Repositories.Asset.Create(asset)
	if err != nil {
		log.Printf("error: create assets %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	data, err := utils.DecryptPlainText(user.EncryptedPrivateKey, h.app.Config.MasterKey)
	if err != nil {
		log.Printf("error: decrypt private Key %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	key, err := hedera.PrivateKeyFromBytesECDSA(data)
	if err != nil {
		log.Printf("error: private key to ECDSA%v", err)
		return helpers.ErrInternalServer(ctx, err)
	}
	accountID, err := hedera.AccountIDFromString(user.AccountID)
	if err != nil {
		log.Printf("error: accountId %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	account := hederaUtils.InitializeAccount(h.app.HederaClient.Client, &accountID)

	account.SetPrivateKey(&key)

	_, err = account.AssociateAccountToToken(*nftCreateRx.TokenID)
	if err != nil {
		log.Printf("error: asscoiate account to tokenId %v", err)
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response.AssetResponseFromModel(createdAsset))

}

func (h assetHandler) RecommendedList(ctx echo.Context) error {
	req := struct {
		Id       string `param:"id" validate:"required,uuid"`
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
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
	asset, err := h.app.Repositories.Asset.GetById(req.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}
	assets, totalCount, err := h.app.Repositories.Asset.ListRecommended(asset, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.AssetStatsResponse{}
	for _, asset := range assets {
		logo := internal.ChainLogoMap[models.Chain(asset.Blockchain)]
		asset.BlockchainLogo = logo
		resp = append(resp, response.AssetStatsResponseFromModel(asset))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": resp,
	})

}

func (h assetHandler) GetChains(ctx echo.Context) error {
	chains := []response.ChainResponse{}
	for id, name := range internal.ChainNameMap {
		logo := internal.ChainLogoMap[id]
		isEnabled := internal.ChainDisabledMap[id]
		chains = append(chains, response.ChainResponse{
			Id:        id,
			Name:      name,
			Logo:      logo,
			IsEnabled: isEnabled,
		})
	}
	return ctx.JSON(http.StatusOK, chains)
}

func (h assetHandler) GetCategories(ctx echo.Context) error {
	categories, err := h.app.Repositories.Category.List()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	resp := []response.CategoryResponse{}
	for _, category := range categories {
		catResp := response.CategoryResponse{
			Id:   category.Id,
			Name: category.Name,
			Slug: category.UrlSlug,
		}
		resp = append(resp, catResp)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h assetHandler) Get(ctx echo.Context) error {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "asset not found")
	}

	authenticatedUser := helpers.ContextGetUser(ctx)

	asset, err := h.app.Repositories.Asset.GetById(parsedId.String())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}

	if utils.IsValidUuid(authenticatedUser.Id) {

		// check if item is a favourite of the user and set the response accordingly
		isFavourite, err := h.app.Repositories.Asset.IsFavourite(
			authenticatedUser.Id,
			asset.Id,
		)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}

		asset.IsFavourite = &isFavourite

		isListedByUser, err := h.app.Repositories.User.HasOpenSellOrder(authenticatedUser.Id, asset.Id)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		asset.IsListedByUser = isListedByUser
		// increment the view count, allowing the request to proceed if there is an error during the process
		// as that's not cause for a fatal response
		isViewed, _ := h.app.Repositories.Asset.IsViewed(
			authenticatedUser.Id,
			asset.Id,
		)

		// only increment view count for users who haven't viewed the item
		if !isViewed {
			err = h.app.Repositories.Asset.AddView(
				authenticatedUser.Id,
				asset.Id,
			)

			// increment previously fetched view count of item to new state
			// only if no error occurs
			if nil == err {
				asset.Views++
			}
		}
	}

	return ctx.JSON(http.StatusOK, response.AssetResponseFromModel(asset))
}

func (h assetHandler) GetListings(ctx echo.Context) error {
	req := struct {
		AssetId  string `param:"asset-id"`
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
	}{}

	err := ctx.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
	}

	if filter.Page < 1 {
		filter.Page = models.DefaultPage
	}
	if filter.Limit < 1 || filter.Limit > models.MaxPageLimit {
		filter.Limit = models.PageLimit8
	}

	if strings.TrimSpace(req.AssetId) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	listings, err := h.app.Repositories.Asset.GetListings(req.AssetId, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.ListingResponse{}
	for _, listing := range listings {
		resp = append(resp, response.ListingResponseFromModel(listing))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h assetHandler) Favourite(ctx echo.Context) error {
	id := ctx.Param("asset-id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "invalid assetId")
	}
	assetId := parsedId.String()

	authenticatedUser := helpers.ContextGetUser(ctx)

	isFavourite, err := h.app.Repositories.Asset.IsFavourite(authenticatedUser.Id, assetId)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	// toggle the favourite status, deleting it if it's currently true
	// and adding it if it's false
	if isFavourite {
		err := h.app.Repositories.Asset.RemoveFavourite(authenticatedUser.Id, assetId)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
	} else {
		err := h.app.Repositories.Asset.AddFavourite(authenticatedUser.Id, assetId)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"isFavourite": !isFavourite})
}

func (h assetHandler) GetAssetOrderBook(ctx echo.Context) error {
	req := struct {
		Id       string              `param:"asset-id" validate:"required,uuid"`
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
	asset, err := h.app.Repositories.Asset.GetById(req.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}
	orders, totalCount, err := h.app.Repositories.Asset.GetOrders(asset.Id, req.Status, req.Type, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.OrderResponse{}
	for _, order := range orders {

		resp = append(resp, response.OrderResponseFromModel(order))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"orders": resp,
	})
}
func (h assetHandler) GetAssetActivities(ctx echo.Context) error {
	req := struct {
		Id       string `param:"asset-id" validate:"required,uuid"`
		Page     int    `query:"page"`
		PageSize int    `query:"size"`
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
	asset, err := h.app.Repositories.Asset.GetById(req.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}
	activities, err := h.app.Repositories.Activity.ListByAssetID(asset.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.ActivityResponse{}
	for _, activity := range activities {

		resp = append(resp, response.ActivityResponseFromModel(activity))
	}

	return ctx.JSON(http.StatusOK, resp)
}
