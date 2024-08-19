package routes

import (
	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Engine initializes, configures and returns the server engine instance.
func Engine(app internal.Application) *echo.Echo {
	var whiteList = []string{"*"}

	engine := echo.New()

	engine.HideBanner = true
	engine.Validator = helpers.NewCustomValidator()

	// set Echo to debug mode on development
	if app.Config.Env == "production" {
		engine.Debug = true

		// TODO: load approved origins from env
		whiteList = []string{}
	}

	engine.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Before(func() {
				origin := c.Request().Header.Get("Origin")
				for _, allowedOrigins := range whiteList {
					// we are in development
					if allowedOrigins == "*" {
						c.Response().Header().Set("Access-Control-Allow-Origin", "*")
						break
					}
					if allowedOrigins == origin {
						c.Response().Header().Set("Access-Control-Allow-Origin", origin)
					}
				}

				c.Response().Header().Set("Access-Control-Allow-Headers", "*")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				c.Response().Header().Add("Vary", "Origin")

			})
			return next(c)
		}
	}))

	// set middleware
	engine.Pre(middleware.RemoveTrailingSlash())
	engine.Use(middleware.Recover())

	engine.Use(middleware.Logger())

	// set routes
	setAuthRoutes(app, engine)
	setUserRoutes(app, engine)
	setAssetsRoutes(app, engine)
	setAdminRoutes(app, engine)
	setListingRoutes(app, engine)
	setupStatRoutes(app, engine)
	setOrderRoutes(app, engine)
	setNotificationRoutes(app, engine)
	// setPortfolioRoutes(app, engine)
	setupWebHooksRoutes(app, engine)

	return engine
}
