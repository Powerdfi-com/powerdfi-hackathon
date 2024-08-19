package models

import "time"

type AssetPriceData struct {
	Timestamp time.Time
	Price     float64
}

type AssetStat struct {
	AssetId          string
	CategoryId       *int
	CategoryName     *string
	Symbol           string
	Name             string
	Blockchain       string
	BlockchainLogo   string
	CreatorId        string
	CreatorUsername  string
	URLs             []string
	TotalVolume      float64
	FloorPrice       float64
	Owners           int
	Currency         string
	Status           string
	ActivityCount    int64
	PercentageChange float64

	IsVerified   bool
	IsMinted     bool
	PriceChanges []AssetPriceData
}

type AssetSurveyMonthlyStat struct {
	Month     int
	Category1 int
	Category2 int
	Category3 int
	Category4 int
}
type CreatorSurveyMonthlyStat struct {
	Month              int
	HeritageUsersCount int
	NewUsersCount      int
}
type VerifiedSurveyMonthlyStat struct {
	Month           int
	VerifiedCount   int
	UnVerifiedCount int
}

type SalesTrend struct {
	Period     string
	Sales      int
	SalesValue float64
}

type CreatorStats struct {
	CreatorId       string
	CreatorUsername string
	AssetsCount     int
}
