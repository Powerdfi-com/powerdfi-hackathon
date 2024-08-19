package middleware

import (
	"net/http"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/labstack/echo/v4"
)

// RequireActivation ensures the user attempting to access a resource
// has been activated.
func RequireActivated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := helpers.ContextGetUser(ctx)
		if !user.IsActive {
			return &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "user is unactivated",
			}
		}

		return next(ctx)
	}
}

// func RequiresAdmin(app internal.Application) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(ctx echo.Context) error {
// 			authenticatedUser := helpers.ContextGetUser(ctx)
// 			user, err := app.Repositories.User.GetByAddress(authenticatedUser.Address)
// 			if err != nil {
// 				return helpers.ErrInternalServer(ctx, err)
// 			}
// 			if !user.HasRole(models.ROLE_ADMIN) {
// 				return &echo.HTTPError{
// 					Code:    http.StatusForbidden,
// 					Message: "user has no admin role",
// 				}
// 			}

//				return next(ctx)
//			}
//		}
//	}
func RequiresSuperAdmin(app internal.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authenticatedAdmin := helpers.ContextGetAdmin(ctx)
			admin, err := app.Repositories.Admin.FindByID(authenticatedAdmin.Id)
			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}
			if !admin.HasRole(models.ROLE_SUPER_ADMIN) {
				return &echo.HTTPError{
					Code:    http.StatusForbidden,
					Message: "privileged not sufficient",
				}
			}

			return next(ctx)
		}
	}
}
