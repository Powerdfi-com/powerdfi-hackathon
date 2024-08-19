package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type adminHandler struct {
	app internal.Application
}

func NewAdminHandler(app internal.Application) adminHandler {
	return adminHandler{app: app}
}

func (h adminHandler) Create(ctx echo.Context) error {

	req := struct {
		Email    string           `json:"email" validate:"required"`
		Password string           `json:"password" validate:"required"`
		Role     models.AdminRole `json:"role" validate:"required"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}
	_, ok := models.AdminRoleMappings[req.Role]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"role": "invalid role",
		})

	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	admin, err := h.app.Repositories.Admin.Create(models.Admin{
		Email:        req.Email,
		PasswordHash: passwordHash,
		RoleMask:     int(req.Role),
	})

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateDetails):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "email already exists",
			}

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, response.AdminResponseFromModel(admin))
}
func (h adminHandler) UpdatePassword(ctx echo.Context) error {

	req := struct {
		NewPassword string `json:"newPassword" validate:"required"`
		Password    string `json:"password" validate:"required"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	authenticatedAdmin := helpers.ContextGetAdmin(ctx)
	admin, err := h.app.Repositories.Admin.FindByEmail(authenticatedAdmin.Email)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	match, err := helpers.MatchPassword(admin.PasswordHash, []byte(req.Password))

	if err != nil || !match {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	passwordHash, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	err = h.app.Repositories.Admin.UpdatePasswordHash(admin.Id, passwordHash)

	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response.AdminResponseFromModel(admin))
}
func (h adminHandler) UpdateAdmin(ctx echo.Context) error {

	req := struct {
		Email    *string           `json:"email"`
		Password *string           `json:"password"`
		Role     *models.AdminRole `json:"role"`
	}{}

	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "asset not found")
	}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	admin, err := h.app.Repositories.Admin.FindByID(parsedId.String())
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	if req.Role != nil {
		_, ok := models.AdminRoleMappings[*req.Role]
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
				"role": "invalid role",
			})

		}
		admin.RoleMask = int(*req.Role)
	}

	if req.Password != nil {
		passwordHash, err := helpers.HashPassword(*req.Password)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}

		admin.PasswordHash = passwordHash

	}
	if req.Email != nil {

		admin.Email = *req.Email

	}
	err = h.app.Repositories.Admin.Update(admin)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateDetails):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "email already exists",
			}

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}
	return ctx.JSON(http.StatusOK, response.AdminResponseFromModel(admin))
}

func (h adminHandler) Approve(ctx echo.Context) error {

	req := struct {
		AssetId string `json:"assetId" validate:"required"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	asset, err := h.app.Repositories.Asset.GetById(req.AssetId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	asset.Status = models.VERIFIED_ASSET_STATUS
	err = h.app.Repositories.Asset.UpdateStatus(asset.Id, asset.Status)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response.AssetResponseFromModel(asset))
}
func (h adminHandler) Reject(ctx echo.Context) error {

	req := struct {
		AssetId string `json:"assetId" validate:"required"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helpers.FormatValidationErr(err))
	}

	asset, err := h.app.Repositories.Asset.GetById(req.AssetId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "asset not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}
	}

	asset.Status = models.REJECTED_ASSET_STATUS

	err = h.app.Repositories.Asset.UpdateStatus(asset.Id, asset.Status)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response.AssetResponseFromModel(asset))
}

func (h adminHandler) ListAssets(ctx echo.Context) error {
	req := struct {
		Page     int  `query:"page"`
		PageSize int  `query:"size"`
		Range    uint `query:"range"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.PageLimit9
	}
	if req.Range == 0 {
		req.Range = models.Range30
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
		Range: req.Range,
	}

	assets, err := h.app.Repositories.Asset.List(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	assetsResponse := []response.AssetResponse{}
	for _, asset := range assets {
		assetsResponse = append(assetsResponse, response.AssetResponseFromModel(asset))
	}

	return ctx.JSON(http.StatusOK, assetsResponse)
}

func (h adminHandler) CreatorsSurvey(ctx echo.Context) error {

	monthlyStats, err := h.app.Repositories.Stats.GetAssetCreationSurvey()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	monthlyStatResponse := []response.CreatorSurveyMonthlyStatResponse{}
	for _, monthlyStat := range monthlyStats {
		monthlyStatResponse = append(monthlyStatResponse, response.CreatorSurveyMonthlyStatResponse{
			Month:              monthlyStat.Month,
			HeritageUsersCount: monthlyStat.HeritageUsersCount,
			NewUsersCount:      monthlyStat.NewUsersCount,
		})
	}

	return ctx.JSON(http.StatusOK, monthlyStatResponse)
}
func (h adminHandler) VerifiedSurvey(ctx echo.Context) error {

	monthlyStats, err := h.app.Repositories.Stats.GetAssetCreationSurvey()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	monthlyStatResponse := []response.VerifiedSurveyMonthlyStatResponse{}
	for _, monthlyStat := range monthlyStats {
		monthlyStatResponse = append(monthlyStatResponse, response.VerifiedSurveyMonthlyStatResponse{
			Month:           monthlyStat.Month,
			VerifiedCount:   monthlyStat.HeritageUsersCount,
			UnverifiedCount: monthlyStat.NewUsersCount,
		})
	}

	return ctx.JSON(http.StatusOK, monthlyStatResponse)
}
func (h adminHandler) AssetSurvey(ctx echo.Context) error {

	monthlyStats, err := h.app.Repositories.Stats.GetAssetCategorySurvey()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	categories, err := h.app.Repositories.Category.List()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	categoryMap := map[string]string{}
	for _, category := range categories {
		categoryMap[strconv.Itoa(category.Id)] = category.Name
	}

	monthlyStatResponse := []map[string]interface{}{}
	for _, monthlyStat := range monthlyStats {
		monthlyStatResponse = append(monthlyStatResponse, map[string]interface{}{
			"month":          monthlyStat.Month,
			categoryMap["1"]: monthlyStat.Category1,
			categoryMap["2"]: monthlyStat.Category2,
			categoryMap["3"]: monthlyStat.Category3,
			categoryMap["4"]: monthlyStat.Category4,
		})
	}

	return ctx.JSON(http.StatusOK, monthlyStatResponse)
}

func (h adminHandler) GetRoles(ctx echo.Context) error {
	roles := []map[string]interface{}{}
	for id, name := range models.AdminRoleMappings {
		roles = append(roles, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}
	return ctx.JSON(http.StatusOK, roles)
}
func (h adminHandler) List(ctx echo.Context) error {
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
	admins, totalCount, err := h.app.Repositories.Admin.List(filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	adminsResponse := []response.AdminResponse{}
	for _, admin := range admins {
		adminsResponse = append(adminsResponse, response.AdminResponseFromModel(admin))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"admins": adminsResponse,
	})
}

func (h adminHandler) Get(ctx echo.Context) error {
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "admin not found")
	}

	// authenticatedAdmin := helpers.ContextGetAdmin(ctx)

	admin, err := h.app.Repositories.Admin.FindByID(parsedId.String())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "admin not found")

		default:
			return helpers.ErrInternalServer(ctx, err)
		}

	}

	return ctx.JSON(http.StatusOK, response.AdminResponseFromModel(admin))
}

func (h adminHandler) AppStats(ctx echo.Context) error {
	usersCount, err := h.app.Repositories.Stats.UsersCount()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)

	}
	creatorsCount, err := h.app.Repositories.Stats.CreatorsCount()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)

	}
	assetsCount, err := h.app.Repositories.Stats.AssetsCount()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)

	}
	percentageChangeCreators, err := h.app.Repositories.Stats.CreatorsWeekIncrement()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)

	}
	percentageChangeUsers, err := h.app.Repositories.Stats.UsersWeekIncrement()
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)

	}

	resp := response.AppStatsResponse{
		UsersCount:               usersCount,
		CreatorsCount:            creatorsCount,
		AssetsCount:              assetsCount,
		PercentageChangeCreators: fmt.Sprintf("%.2f%%", percentageChangeCreators),
		PercentageChangeUsers:    fmt.Sprintf("%.2f%%", percentageChangeUsers),
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (h adminHandler) GetNotifications(ctx echo.Context) error {
	authenticatedAdmin := helpers.ContextGetAdmin(ctx)

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
	notifications, totalCount, err := h.app.Repositories.Notification.GetForAdmin(authenticatedAdmin.Id, filter)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := []response.AdminNotificationResponse{}
	for _, notification := range notifications {
		resp = append(resp, response.AdminNotificationResponseFromModel(notification))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":         totalCount,
		"notifications": resp,
	})

}

func (h adminHandler) MarkAllNotificationsAsRead(ctx echo.Context) error {
	authenticatedAdmin := helpers.ContextGetAdmin(ctx)
	err := h.app.Repositories.Notification.MarkAllAsReadAdmin(authenticatedAdmin.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, map[string]string{})
}

func (h adminHandler) CountNotifications(ctx echo.Context) error {
	authenticatedAdmin := helpers.ContextGetAdmin(ctx)

	count, err := h.app.Repositories.Notification.CountForAdmin(authenticatedAdmin.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

func (h adminHandler) GetNotificationPrefs(ctx echo.Context) error {
	authenticatedAdmin := helpers.ContextGetAdmin(ctx)

	prefs, err := h.app.Repositories.Notification.GetAdminUserPrefs(authenticatedAdmin.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := response.AdminNotificationResponsePrefsFromModel(prefs)
	return ctx.JSON(http.StatusOK, resp)
}

func (h adminHandler) UpdateNotificationPrefs(ctx echo.Context) error {
	authenticatedAdmin := helpers.ContextGetAdmin(ctx)

	req := struct {
		Created *bool `json:"created"`
		Login   *bool `json:"login"`
	}{}

	err := ctx.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest
	}

	prefs, err := h.app.Repositories.Notification.GetAdminUserPrefs(authenticatedAdmin.Id)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	if req.Created != nil {
		prefs.Created = *req.Created
	}

	if req.Login != nil {
		prefs.Login = *req.Login
	}

	newPrefs, err := h.app.Repositories.Notification.UpdateAdminPrefs(authenticatedAdmin.Id, prefs)
	if err != nil {
		return helpers.ErrInternalServer(ctx, err)
	}

	resp := response.AdminNotificationResponsePrefsFromModel(newPrefs)
	return ctx.JSON(http.StatusOK, resp)
}
