package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setAuthRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewAuthHandler(app)
	auth := engine.Group("auth")
	auth.GET("/:address/nonce", handler.GetNonce)
	auth.POST("/:address/verify", handler.VerifySignature)
	auth.POST("/refresh-token", handler.RefreshToken)
	auth.POST("/admin", handler.AuthenticateAdmin)

}
