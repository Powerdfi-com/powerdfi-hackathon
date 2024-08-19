package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/cmd/api/middleware"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setUserRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewUserHandler(app)
	user := engine.Group("user")

	user.PATCH("/activate", handler.ActivateUser, middleware.Authentication(app))
	user.GET("/wallet", handler.GetWalletDeets, middleware.Authentication(app))
	user.GET("/:id", handler.Get)
	user.GET("/me", handler.GetOwnProfile, middleware.Authentication(app))
	user.PATCH("/me", handler.UpdateProfile, middleware.Authentication(app))
	user.GET("/kyc-link", handler.GetKycLink, middleware.Authentication(app))
	user.POST("/kyb", handler.InitiateKYB, middleware.Authentication(app))
	user.GET("/kyb/status", handler.GetKYBStatus, middleware.Authentication(app))
	user.GET("/kyb/countries", handler.GetKYBCountries, middleware.Authentication(app))
	user.GET("/:userId/created", handler.GetCreatedAssets, middleware.PublicAuthentication(app))
	user.GET("/:userId/listings", handler.GetListedItems)
	user.GET("/:userId/unlisted", handler.GetUnListedItems)
	user.GET("/portfolio", handler.GetPortfolio, middleware.Authentication(app))
	user.GET("/orders", handler.GetUserOrders, middleware.Authentication(app))
	user.GET("/activities", handler.GetUserActivities, middleware.Authentication(app))
	user.GET("/top-assets/perfs", handler.TopAssetsPerfs, middleware.Authentication(app))
	user.GET("/trending-assets/perfs", handler.TrendingAssetsPerfs, middleware.Authentication(app))
	user.GET("/sales-trend", handler.SalesTrend, middleware.Authentication(app))

}
