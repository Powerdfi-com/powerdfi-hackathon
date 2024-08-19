package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setupWebHooksRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewWebHooksHandler(app)
	webhook := engine.Group("webhooks")
	webhook.POST("/shufti", handler.Shufti)

}
