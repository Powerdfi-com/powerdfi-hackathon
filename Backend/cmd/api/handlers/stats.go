package handlers

import (
	"fmt"
	"net/http"

	"github.com/Powerdfi-com/Backend/cmd/api/helpers"
	"github.com/Powerdfi-com/Backend/cmd/api/models/response"
	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/labstack/echo/v4"
)

type statHandler struct {
	app internal.Application
}

func NewStatHandler(app internal.Application) statHandler {
	return statHandler{app: app}
}

func (h statHandler) TopAssets(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
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

	topAssets, totalCount, err := h.app.Repositories.Stats.TopAssets(filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("totalCount", totalCount)

	topAssetsResponse := []response.AssetStatsResponse{}
	for _, topAsset := range topAssets {
		logo := internal.ChainLogoMap[models.Chain(topAsset.Blockchain)]
		topAsset.BlockchainLogo = logo
		topAssetsResponse = append(topAssetsResponse, response.AssetStatsResponseFromModel(topAsset))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": topAssetsResponse,
	})
}
func (h statHandler) TrendingAssets(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
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

	trendingAssets, totalCount, err := h.app.Repositories.Stats.TrendingAssets(filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	trendingAssetsResponse := []response.AssetStatsResponse{}

	for _, trendingAsset := range trendingAssets {
		logo := internal.ChainLogoMap[models.Chain(trendingAsset.Blockchain)]
		trendingAsset.BlockchainLogo = logo
		trendingAssetsResponse = append(trendingAssetsResponse, response.AssetStatsResponseFromModel(trendingAsset))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": trendingAssetsResponse,
	})
}

func (h statHandler) TopAssetsPerfs(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
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
	// authenticatedUser := helpers.ContextGetUser(ctx)

	topAssets, totalCount, err := h.app.Repositories.Stats.TopAssetPerfs(filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	topAssetsResponse := []response.AssetStatsResponse{}
	for _, topAsset := range topAssets {
		logo := internal.ChainLogoMap[models.Chain(topAsset.Blockchain)]
		topAsset.BlockchainLogo = logo
		priceChanges, err := h.app.Repositories.Stats.GetAssetPriceData(topAsset.AssetId, filter)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		topAsset.PriceChanges = priceChanges
		topAssetsResponse = append(topAssetsResponse, response.AssetStatsResponseFromModel(topAsset))
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": topAssetsResponse,
	})
}

func (h statHandler) TrendingAssetsPerfs(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
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

	// authenticatedUser := helpers.ContextGetUser(ctx)

	trendingAssets, totalCount, err := h.app.Repositories.Stats.TrendingAssetPerfs(filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	trendingAssetsResponse := []response.AssetStatsResponse{}

	for _, trendingAsset := range trendingAssets {
		logo := internal.ChainLogoMap[models.Chain(trendingAsset.Blockchain)]
		trendingAsset.BlockchainLogo = logo
		priceChanges, err := h.app.Repositories.Stats.GetAssetPriceData(trendingAsset.AssetId, filter)
		if err != nil {
			return helpers.ErrInternalServer(ctx, err)
		}
		trendingAsset.PriceChanges = priceChanges
		trendingAssetsResponse = append(trendingAssetsResponse, response.AssetStatsResponseFromModel(trendingAsset))
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":  totalCount,
		"assets": trendingAssetsResponse,
	})
}

func (h statHandler) TopCreators(ctx echo.Context) error {
	req := struct {
		Page       int     `query:"page"`
		PageSize   int     `query:"size"`
		Range      uint    `query:"range"`
		CategoryId *int    `query:"categoryId"`
		Blockchain *string `query:"blockchain"`
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

	// authenticatedUser := helpers.ContextGetUser(ctx)

	topCreators, totalCount, err := h.app.Repositories.Stats.TopCreators(filter, req.CategoryId, req.Blockchain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	topCreatorsResponse := []response.CreatorStatsResponse{}

	for _, topCreator := range topCreators {
		topCreatorsResponse = append(topCreatorsResponse, response.CreatorStatsResponse{
			CreatorUsername: topCreator.CreatorUsername,
			CreatorId:       topCreator.CreatorId,
			AssetsCount:     topCreator.AssetsCount,
		})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"total":    totalCount,
		"creators": topCreatorsResponse,
	})
}
