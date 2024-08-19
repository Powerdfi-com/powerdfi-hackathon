package models

import "time"

type ListingCurrency string

type ListingType string

type Listing struct {
	Id          string
	AssetId     string
	AssetName   string
	AssetSymbol string
	Type        ListingType
	UserId      string

	PriceUSD float64

	MinInvestAmount *float64
	MaxInvestAmount *float64
	MinToRaise      *float64
	MaxToRaise      *float64

	Currency  []ListingCurrency
	Quantity  int64
	StartDate *time.Time
	EndDate   *time.Time

	IsActive    bool
	IsCancelled bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
