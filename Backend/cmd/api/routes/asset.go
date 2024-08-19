package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/cmd/api/middleware"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setAssetsRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewAssetsHandler(app)
	assets := engine.Group("assets")

	assets.GET("/categories", handler.GetCategories)
	assets.GET("/chains", handler.GetChains)
	assets.POST("", handler.Create, middleware.Authentication(app))
	assets.GET("/:id", handler.Get, middleware.PublicAuthentication(app))
	assets.GET("/:id/recommended", handler.RecommendedList, middleware.PublicAuthentication(app))
	assets.GET("/:asset-id/listings", handler.GetListings)
	assets.POST("/:asset-id/favourite", handler.Favourite, middleware.Authentication(app), middleware.RequireActivated)
	assets.GET("/:asset-id/orderbook", handler.GetAssetOrderBook, middleware.PublicAuthentication(app))
	assets.GET("/:asset-id/activities", handler.GetAssetActivities, middleware.PublicAuthentication(app))
}
