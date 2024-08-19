package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/cmd/api/middleware"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setNotificationRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewNotificationHandler(app)
	notification := engine.Group("notifications")

	notification.GET("", handler.Get, middleware.Authentication(app), middleware.RequireActivated)
	notification.PUT("/read", handler.MarkAllAsRead, middleware.Authentication(app), middleware.RequireActivated)
	notification.PUT("/:id/read", handler.MarkAsRead, middleware.Authentication(app), middleware.RequireActivated)
	notification.GET("/count", handler.Count, middleware.Authentication(app), middleware.RequireActivated)
	notification.GET("/prefs", handler.GetPrefs, middleware.Authentication(app), middleware.RequireActivated)
	notification.POST("/prefs", handler.UpdatePrefs, middleware.Authentication(app), middleware.RequireActivated)

}
