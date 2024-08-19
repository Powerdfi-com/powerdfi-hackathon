package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/handlers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
)

func setupStatRoutes(app internal.Application, engine *echo.Echo) {
	handler := handlers.NewStatHandler(app)
	stat := engine.Group("stats")
	stat.GET("/top-assets", handler.TopAssets)
	stat.GET("/trending-assets", handler.TrendingAssets)
	stat.GET("/top-assets/perfs", handler.TopAssetsPerfs)
	stat.GET("/trending-assets/perfs", handler.TrendingAssetsPerfs)
	stat.GET("/top-creators", handler.TopCreators)
}
