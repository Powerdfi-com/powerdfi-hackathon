package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/external/shufti"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/labstack/echo/v4"
)

type webhookHandler struct {
	app internal.Application
}

func NewWebHooksHandler(app internal.Application) webhookHandler {
	return webhookHandler{app: app}
}

func (h webhookHandler) Shufti(ctx echo.Context) error {
	req := shufti.VerificationParams{}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadRequest
	}
	defer ctx.Request().Body.Close()

	signature := ctx.Request().Header.Get("Signature")

	if !h.app.ShuftiClient.IsValidSig(signature, body) {
		ctx.Logger().Errorf("invalid sig")
		return echo.ErrBadRequest
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadRequest
	}

	user_kyc, err := h.app.Repositories.UserKyc.GetByReferenceId(req.ReferenceId)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadRequest
	}

	switch req.Event {
	case shufti.EVENT_TYPE_ACCEPTED:
		// TODO: update first_name and last_name
		user_kyc.Status = models.KYC_STATUS_SUCCESS
		user_kyc.Comment = ""
		err := h.app.Repositories.UserKyc.Update(user_kyc)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		err = h.app.Repositories.User.Verify(user_kyc.UserId)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		return ctx.JSON(http.StatusOK, "processed")

	case shufti.EVENT_TYPE_DECLINED:
		user_kyc.Status = models.KYC_STATUS_FAILED
		user_kyc.Comment = req.DeclinedReason
		err := h.app.Repositories.UserKyc.Update(user_kyc)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}

		return ctx.JSON(http.StatusOK, "processed")
	default:
		return ctx.JSON(http.StatusOK, "skipped")

	}

}
