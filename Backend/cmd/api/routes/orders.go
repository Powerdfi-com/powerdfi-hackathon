package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/cmd/api/middleware"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setOrderRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewOrderHandler(app)
	orders := engine.Group("orders")

	orders.POST("", handler.Create, middleware.Authentication(app), middleware.RequireActivated)
	orders.DELETE("/:id", handler.Cancel, middleware.Authentication(app), middleware.RequireActivated)

}
