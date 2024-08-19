package helpers

import (
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ContextKey string

const ContextKeyUser = "user"

// ContextGetUser returns the current user data stored in the request context.
// Only the ID and Address of a user are stored.
func ContextGetUser(ctx echo.Context) models.User {
	token, ok := ctx.Get(ContextKeyUser).(*jwt.Token)
	if !ok {
		return models.User{}
	}
	claims, ok := token.Claims.(*CustomAccessJwtClaims)
	if !ok {
		return models.User{}
	}

	return models.User{
		Id:       claims.Id,
		Address:  claims.Address,
		IsActive: claims.Activated,
	}
}
func ContextGetAdmin(ctx echo.Context) models.Admin {
	token, ok := ctx.Get(ContextKeyUser).(*jwt.Token)
	if !ok {
		return models.Admin{}
	}
	claims, ok := token.Claims.(*CustomAdminAccessJwtClaims)
	if !ok {
		return models.Admin{}
	}

	return models.Admin{
		Id:    claims.Id,
		Email: claims.Email,
	}
}
