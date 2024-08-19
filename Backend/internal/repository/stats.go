package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type StatsRepository interface {
	TopAssetPerfs(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error)
	TrendingAssetPerfs(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error)
	TopAssets(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error)
	TrendingAssets(filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error)
	GetAssetPriceData(assetId string, filter models.Filter) ([]models.AssetPriceData, error)
	GetAssetCreationSurvey() ([]models.CreatorSurveyMonthlyStat, error)
	GetAssetCategorySurvey() ([]models.AssetSurveyMonthlyStat, error)
	TrendingUserAssetPerfs(userId string, filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error)
	TopUserAssetPerfs(userId string, filter models.Filter, categoryId *int, blockchain *string) ([]models.AssetStat, int, error)
	SalesTrend(userId string, filter models.Filter) ([]models.SalesTrend, error)
	CreatorsWeekIncrement() (float64, error)
	UsersWeekIncrement() (float64, error)
	AssetsCount() (int, error)
	CreatorsCount() (int, error)
	UsersCount() (int, error)
	TopCreators(filter models.Filter, categoryId *int, blockchain *string) ([]models.CreatorStats, int, error)
}
