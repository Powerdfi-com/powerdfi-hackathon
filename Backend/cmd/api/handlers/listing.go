package handlers

import (
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type listingHandler struct {
	app internal.Application
}

func NewListingHandler(app internal.Application) listingHandler {
	return listingHandler{app: app}
}

func (h listingHandler) Create(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	req := struct {
		AssetId string `json:"assetId" validate:"required"`
		// Type            string   `json:"type" validate:"required,oneof=fixed auction"`
		Price           float64  `json:"price" validate:"required"`
		MinInvestAmount *float64 `json:"min_investment_amount"`
		MaxInvestAmount *float64 `json:"max_investment_amount"`
		MaxRiseAmount   *float64 `json:"max_raise_amount"`
		MinRaiseAmount  *float64 `json:"min_raise_amount"`
		Currency        []string `json:"currency" validate:"required"`
		Quantity        int64    `json:"quantity" validate:"required,min=1,max=100000"`

		StartDate *time.Time `json:"startAt" validate:"required"`
		EndDate   *time.Time `json:"endAt" validate:"required"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	var err error
	asset, err := h.app.Repositories.Asset.GetById(req.AssetId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {

			return echo.NewHTTPError(http.StatusNotFound, "asset not found")
		}

		return helpers.ErrInternalServer(ctx, err)
	}

	currency := make([]models.ListingCurrency, 0)
	blockchainTokens, ok := internal.ChainCurrencyListMap[asset.Blockchain]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "asset blockchain not supported")

	}
	// TODO: validate if it part of supported currency
	for _, c := range req.Currency {
		if !slices.Contains(blockchainTokens, models.ListingCurrency(c)) {
			return echo.NewHTTPError(http.StatusBadRequest, "currency not supported for blockchain")
		}

		currency = append(currency, models.ListingCurrency(c))
	}

	if authenticatedUser.Id != asset.CreatorUserID {
		// check balance of item
		return echo.NewHTTPError(http.StatusBadRequest, "user is not creator")
	}

	listing := models.Listing{
		// Type:    req.Type,
		UserId:      authenticatedUser.Id,
		AssetId:     asset.Id,
		AssetName:   asset.Name,
		AssetSymbol: asset.Symbol,

		PriceUSD:        req.Price,
		MinInvestAmount: req.MinInvestAmount,
		MaxInvestAmount: req.MaxInvestAmount,
		MaxToRaise:      req.MaxRiseAmount,
		MinToRaise:      req.MinRaiseAmount,
		Currency:        currency,
		Quantity:        req.Quantity,

		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	createdListing, err := h.app.Repositories.Listing.Create(listing)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrActiveListing):
			return echo.NewHTTPError(http.StatusConflict, "asset is currently listed by user")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}

	resp := response.ListingResponseFromModel(createdListing)
	return ctx.JSON(http.StatusCreated, resp)

}

func (h listingHandler) Cancel(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed listing id")
	}

	// fetch listing to compare user ids
	listing, err := h.app.Repositories.Listing.GetById(parsedId.String())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset listing not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if authenticatedUser.Id != listing.UserId {
		return echo.NewHTTPError(http.StatusForbidden, "asset not listed by user")
	}

	// proceed to cancel listing if everything checks out
	err = h.app.Repositories.Listing.Cancel(parsedId.String())
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "item listing cancelled")
}

func (h listingHandler) GetChainCurrencies(ctx echo.Context) error {
	chainId := ctx.Param("chainId")
	currencies, ok := internal.ChainCurrencyListMap[models.Chain(chainId)]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "chain not found")
	}

	currenciesResp := []response.ListingCurrencyResponse{}
	for _, id := range currencies {
		name := internal.CurrencyNameMap[id]
		currenciesResp = append(currenciesResp, response.ListingCurrencyResponse{
			Id:   id,
			Name: name,
		})
	}
	return ctx.JSON(http.StatusOK, currenciesResp)
}
