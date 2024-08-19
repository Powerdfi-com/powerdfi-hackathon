package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	app internal.Application
}

func NewOrderHandler(app internal.Application) orderHandler {
	return orderHandler{app: app}
}

func (h orderHandler) Create(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	req := struct {
		AssetId  string  `json:"assetId" validate:"required"`
		Price    float64 `json:"price" `
		Quantity int64   `json:"quantity" validate:"required,min=1,max=100000"`
		Type     string  `json:"type" validate:"required"`
		Kind     string  `json:"kind"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	asset, err := h.app.Repositories.Asset.GetById(req.AssetId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {

			return echo.NewHTTPError(http.StatusNotFound, "asset not found")
		}

		return helpers.ErrInternalServer(ctx, err)
	}

	// TODO: prevent creating sell order if not creator of asset
	var orderPrice *float64

	switch req.Type {
	// TODO: if sell order ensure user has the quantity to sell
	case string(models.ORDER_BUY_TYPE):
	case string(models.ORDER_SELL_TYPE):
		// if authenticatedUser.Id != asset.CreatorUserID {
		// 	return echo.NewHTTPError(http.StatusForbidden, "user is not the creator of the asset")
		// }
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "order type not supported")

	}
	switch req.Kind {
	case string(models.ORDER_LIMIT_KIND):
		if req.Price <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "limit order requires a specified price")
		}
		orderPrice = &req.Price
	case string(models.ORDER_MARKET_KIND):
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "order kind not supported")

	}
	fmt.Println("orderPrice", orderPrice)
	// TODO: validate the request quantitiy against the asset total supply
	order := models.Order{
		// Type:    req.Type,
		UserId:  authenticatedUser.Id,
		AssetId: asset.Id,

		Price: req.Price,

		Quantity: req.Quantity,
		Type:     models.OrderType(req.Type),
		Status:   models.ORDER_OPEN_STATUS,
		Kind:     models.OrderKind(req.Kind),
	}

	createdOrder, err := h.app.Repositories.Order.Create(order)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrActiveListing):
			return echo.NewHTTPError(http.StatusConflict, "asset is currently listed by user")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}

	resp := response.OrderResponseFromModel(createdOrder)
	return ctx.JSON(http.StatusCreated, resp)

}

func (h orderHandler) Cancel(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed order id")
	}

	// fetch order to compare user ids
	order, err := h.app.Repositories.Order.GetById(parsedId.String())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "order not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	if authenticatedUser.Id != order.UserId {
		return echo.NewHTTPError(http.StatusForbidden, "order not created by user")
	}

	// proceed to cancel listing if everything checks out
	err = h.app.Repositories.Order.Cancel(parsedId.String())
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "order cancelled cancelled")
}
