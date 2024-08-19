package middleware

import (
	"strings"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Authentication(app internal.Application) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.CustomAccessJwtClaims)
		},
		SigningKey: []byte(app.Config.Jwt.Access),
	})
}

func AdminAuthentication(app internal.Application) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.CustomAdminAccessJwtClaims)
		},
		SigningKey: []byte(app.Config.Jwt.AdminAccess),
	})
}
func PublicAuthentication(app internal.Application) echo.MiddlewareFunc {

	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.CustomAccessJwtClaims)
		},
		SigningKey:             []byte(app.Config.Jwt.Access),
		ContinueOnIgnoredError: true,
		ErrorHandler: func(ctx echo.Context, err error) error {
			token := ctx.Request().Header.Get(echo.HeaderAuthorization)
			if strings.TrimSpace(token) == "" {
				// set an empty token context for requests without a provided JWTs
				ctx.Set(helpers.ContextKeyUser, new(jwt.Token))
				return nil
			}

			// return an error for an invalid provided JWT
			return &echo.HTTPError{
				Code:    middleware.ErrJWTInvalid.Code,
				Message: middleware.ErrJWTInvalid.Message,
			}
		},
	})

}
