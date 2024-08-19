package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/cmd/api/middleware"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setAdminRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewAdminHandler(app)
	admin := engine.Group("admin")

	notification := admin.Group("/notifications")

	notification.GET("", handler.GetNotifications, middleware.AdminAuthentication(app))
	notification.PUT("/read", handler.MarkAllNotificationsAsRead, middleware.AdminAuthentication(app))
	notification.GET("/count", handler.CountNotifications, middleware.AdminAuthentication(app))
	notification.GET("/prefs", handler.GetNotificationPrefs, middleware.AdminAuthentication(app))
	notification.POST("/prefs", handler.UpdateNotificationPrefs, middleware.AdminAuthentication(app))

	admin.POST("", handler.Create, middleware.AdminAuthentication(app), middleware.RequiresSuperAdmin(app))
	admin.GET("", handler.List, middleware.AdminAuthentication(app))
	admin.GET("/:id", handler.Get, middleware.AdminAuthentication(app))
	admin.PATCH("", handler.UpdatePassword, middleware.AdminAuthentication(app))

	admin.PATCH("/:id", handler.UpdateAdmin, middleware.AdminAuthentication(app), middleware.RequiresSuperAdmin(app))
	admin.GET("/roles", handler.GetRoles, middleware.AdminAuthentication(app))
	admin.POST("/approve", handler.Approve, middleware.AdminAuthentication(app))
	admin.POST("/reject", handler.Reject, middleware.AdminAuthentication(app))
	admin.GET("/assets", handler.ListAssets, middleware.AdminAuthentication(app))
	admin.GET("/stats", handler.AppStats, middleware.AdminAuthentication(app))
	admin.GET("/survey/creators", handler.CreatorsSurvey, middleware.AdminAuthentication(app))
	admin.GET("/survey/verified-assets", handler.VerifiedSurvey, middleware.AdminAuthentication(app))
	admin.GET("/survey/categories", handler.AssetSurvey, middleware.AdminAuthentication(app))

}
