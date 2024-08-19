package handlers

import (
	"net/http"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type notificationHandler struct {
	app internal.Application
}

func NewNotificationHandler(app internal.Application) notificationHandler {
	return notificationHandler{app: app}
}

func (h notificationHandler) Get(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	req := struct {
		Page     int `query:"page"`
		PageSize int `query:"size"`
	}{}
	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 || req.PageSize > models.MaxPageLimit {
		req.PageSize = models.PageLimit8
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
	}
	notifications, totalCount, err := h.app.Repositories.Notification.GetForUser(authenticatedUser.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.NotificationResponse{}
	for _, notification := range notifications {
		resp = append(resp, response.NotificationResponseFromModel(notification))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":         totalCount,
		"notifications": resp,
	})

}
func (h notificationHandler) MarkAsRead(ctx echo.Context) error {
	notificationId := ctx.Param("id")
	parsedId, err := uuid.Parse(notificationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "notification not found")
	}
	authenticatedUser := helpers.ContextGetUser(ctx)
	err = h.app.Repositories.Notification.MarkAsRead(authenticatedUser.Id, parsedId.String())
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, map[string]string{})
}

func (h notificationHandler) MarkAllAsRead(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)
	err := h.app.Repositories.Notification.MarkAllAsRead(authenticatedUser.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, map[string]string{})
}

func (h notificationHandler) Count(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	count, err := h.app.Repositories.Notification.CountForUser(authenticatedUser.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

func (h notificationHandler) GetPrefs(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	prefs, err := h.app.Repositories.Notification.GetUserPrefs(authenticatedUser.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := response.NotificationResponsePrefsFromModel(prefs)
	return ctx.JSON(http.StatusOK, resp)
}

func (h notificationHandler) UpdatePrefs(ctx echo.Context) error {
	authenticatedUser := helpers.ContextGetUser(ctx)

	req := struct {
		Sale     *bool `json:"sale"`
		Verified *bool `json:"verified"`
		Rejected *bool `json:"rejected"`
		Login    *bool `json:"login"`
	}{}

	err := ctx.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest
	}

	prefs, err := h.app.Repositories.Notification.GetUserPrefs(authenticatedUser.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	if req.Sale != nil {
		prefs.Sale = *req.Sale
	}
	if req.Verified != nil {
		prefs.Verified = *req.Verified
	}
	if req.Rejected != nil {
		prefs.Rejected = *req.Rejected
	}
	if req.Login != nil {
		prefs.Login = *req.Login
	}

	newPrefs, err := h.app.Repositories.Notification.UpdateUserPrefs(authenticatedUser.Id, prefs)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := response.NotificationResponsePrefsFromModel(newPrefs)
	return ctx.JSON(http.StatusOK, resp)
}
