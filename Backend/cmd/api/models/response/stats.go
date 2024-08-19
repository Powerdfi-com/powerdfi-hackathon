package response

import (
	"fmt"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type AssetStatsResponse struct {
	AssetId      string  `json:"assetId"`
	CategoryId   *int    `json:"category,omitempty"`
	CategoryName *string `json:"categoryName,omitempty"`
	Name         string  `json:"name"`

	Blockchain      string `json:"blockchain"`
	BlockchainLogo  string `json:"blockchainLogo"`
	CreatorId       string `json:"creatorId"`
	CreatorUsername string `json:"creatorUsername"`
	Logo            string `json:"logo"`
	Volume          int    `json:"volume"`
	Status          string `json:"status,omitempty"`

	Owners           int                          `json:"owners"`
	FloorPrice       float64                      `json:"floorPrice"`
	IsVerified       bool                         `json:"isVerified"`
	PercentageChange string                       `json:"percentageChange"`
	PriceChanges     []AssetStatPriceDataResponse `json:"priceChanges,omitempty"`
}

type AssetStatPriceDataResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
}

type CreatorSurveyMonthlyStatResponse struct {
	Month              int `json:"month"`
	HeritageUsersCount int `json:"heritageUsersCount"`
	NewUsersCount      int `json:"newUsersCount"`
}
type VerifiedSurveyMonthlyStatResponse struct {
	Month           int `json:"month"`
	VerifiedCount   int `json:"verifiedAssetsCount"`
	UnverifiedCount int `json:"unverifiedAssetsCount"`
}

type SalesTrendResponse struct {
	Period     string  `json:"period"`
	Sales      int     `json:"sales"`
	SalesValue float64 `json:"salesValue"`
}
type AppStatsResponse struct {
	UsersCount               int    `json:"usersCount"`
	CreatorsCount            int    `json:"creatorsCount"`
	AssetsCount              int    `json:"assetsCount"`
	PercentageChangeCreators string `json:"percentageChangeCreators"`
	PercentageChangeUsers    string `json:"percentageChangeUsers"`
}

type CreatorStatsResponse struct {
	CreatorId       string `json:"creatorId"`
	CreatorUsername string `json:"creatorUsername"`
	AssetsCount     int    `json:"assetsCount"`
}

func AssetStatPriceDataResponseFromModel(priceData models.AssetPriceData) AssetStatPriceDataResponse {
	return AssetStatPriceDataResponse{
		Timestamp: priceData.Timestamp,
		Price:     priceData.Price,
	}
}

func AssetStatsResponseFromModel(assetStat models.AssetStat) AssetStatsResponse {

	priceChanges := []AssetStatPriceDataResponse{}

	for _, PriceChange := range assetStat.PriceChanges {
		priceChanges = append(priceChanges, AssetStatPriceDataResponseFromModel(PriceChange))
	}

	return AssetStatsResponse{
		AssetId:      assetStat.AssetId,
		Name:         assetStat.Name,
		CategoryName: assetStat.CategoryName,
		CategoryId:   assetStat.CategoryId,

		Blockchain: assetStat.Blockchain,
		Logo:       assetStat.URLs[0],

		Volume: int(assetStat.TotalVolume),

		BlockchainLogo:  assetStat.BlockchainLogo,
		CreatorId:       assetStat.CreatorId,
		CreatorUsername: assetStat.CreatorUsername,

		Owners:           assetStat.Owners,
		Status:           assetStat.Status,
		FloorPrice:       assetStat.FloorPrice,
		IsVerified:       assetStat.IsVerified,
		PercentageChange: fmt.Sprintf("%.2f%%", assetStat.PercentageChange),
		PriceChanges:     priceChanges,
	}
}
