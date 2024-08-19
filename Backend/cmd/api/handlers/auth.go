package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	app internal.Application
}

func NewAuthHandler(app internal.Application) authHandler {
	return authHandler{app: app}
}

func (a authHandler) GetNonce(ctx echo.Context) error {
	address := ctx.Param("address")

	// delete expired nonces of unregistered users before attempting any further queries
	if err := a.app.Repositories.Nonce.DeleteExpired(); err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := struct {
		Nonce string `json:"nonce"`
	}{}

	// find the nonce associated with the account
	nonce, err := a.app.Repositories.Nonce.Get(address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			// create nonce for non-existing users
			nonce.Message, err = helpers.GenerateNonceMessage(models.NonceActionCreate)
			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}

			err = a.app.Repositories.Nonce.Create(address, nonce.Message)
			if err != nil {
				return helpers.ErrInternalServer(ctx, err)
			}

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	resp.Nonce = nonce.Message
	return ctx.JSON(http.StatusOK, resp)

}

// VerifySignature generates and returns access and refresh JWTs alongside the user profile of an address
// if the signer of the signature in the request body matches the given address.
// A new user account is created if one with the address is not found.
func (a authHandler) VerifySignature(ctx echo.Context) error {
	address := ctx.Param("address")

	reqBody := struct {
		Signature string `json:"signature" validate:"required"`
	}{}
	if err := ctx.Bind(&reqBody); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}
	// get the nonce with the unsigned message
	nonce, err := a.app.Repositories.Nonce.Get(address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	// retrieve the signer's address
	signer, err := helpers.RecoverAddressFromSignature(nonce.Message, reqBody.Signature)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed signature")
	}

	// verify the signer is the same as the given address

	if !strings.EqualFold(address, signer) {

		return echo.NewHTTPError(http.StatusBadRequest, "invalid signature")
	}

	// retrieve an existing user, and create an account for a new one
	statusCode := http.StatusOK
	var user models.User

	if nonce.Action == models.NonceActionCreate {
		statusCode = http.StatusCreated
		newUser := models.User{Address: address}

		user, err = a.app.Repositories.User.Create(newUser)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
	} else {
		user, err = a.app.Repositories.User.GetByAddress(address)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				return echo.ErrNotFound

			default:
				return helpers.ErrInternalServer(ctx, err)
			}
		}
	}

	// generate and update new nonce for user for the next time it's needed
	nonceMessage, err := helpers.GenerateNonceMessage(models.NonceActionUpdate)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	err = a.app.Repositories.Nonce.Update(address, nonceMessage)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	// generate access and refresh tokens
	accessToken, refreshToken, expiration, err := helpers.GenerateTokens(a.app, user)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := map[string]interface{}{
		"user": response.UserResponseFromModel(user),
		"tokens": response.TokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiration,
		},
	}
	return ctx.JSON(statusCode, resp)
}

func (a authHandler) RefreshToken(ctx echo.Context) error {
	reqBody := struct {
		RefreshToken string `json:"refreshToken"`
	}{}

	if err := ctx.Bind(&reqBody); err != nil {
		return echo.ErrBadRequest
	}

	// validate refresh token
	claims, err := helpers.ValidateRefreshToken(a.app.Config.Jwt.Refresh, reqBody.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired refresh token")
	}

	// retrieve user details to generate new access token with
	user, err := a.app.Repositories.User.GetByAddress(claims.Address)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.ErrNotFound

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	accessToken, _, expiration, err := helpers.GenerateTokens(a.app, user)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := response.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: reqBody.RefreshToken,
		ExpiresAt:    expiration,
	}
	return ctx.JSON(http.StatusCreated, resp)
}

func (h authHandler) AuthenticateAdmin(ctx echo.Context) error {

	req := struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	admin, err := h.app.Repositories.Admin.FindByEmail(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	match, err := helpers.MatchPassword(admin.PasswordHash, []byte(req.Password))

	if err != nil || !match {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	accessToken, refreshToken, expiration, err := helpers.GenerateAdminTokens(h.app, admin)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := map[string]interface{}{
		"data": response.AdminResponseFromModel(admin),
		"tokens": response.TokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiration,
		},
	}

	return ctx.JSON(http.StatusOK, resp)
}
