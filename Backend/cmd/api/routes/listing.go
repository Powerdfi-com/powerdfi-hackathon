package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/cmd/api/middleware"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setListingRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewListingHandler(app)
	listing := engine.Group("listings")
	listing.GET("/:chainId/tokens", handler.GetChainCurrencies, middleware.PublicAuthentication(app))
	listing.POST("", handler.Create, middleware.Authentication(app))

	listing.POST("/:id/cancel", handler.Cancel, middleware.Authentication(app), middleware.RequireActivated)

}
