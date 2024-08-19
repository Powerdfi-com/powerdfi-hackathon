package response

import (
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type ListingCurrencyResponse struct {
	Id   models.ListingCurrency `json:"id"`
	Name string                 `json:"name"`
}

type ListingResponse struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`

	Type        models.ListingType `json:"type,omitempty"`
	AssetId     string             `json:"assetId"`
	AssetName   string             `json:"assetName"`
	AssetSymbol string             `json:"assetSymbol"`

	Price float64 `json:"price"`

	Currency        []models.ListingCurrency `json:"currency"`
	Quantity        int64                    `json:"quantity"`
	MinInvestAmount *float64                 `json:"min_investment_amount"`
	MaxInvestAmount *float64                 `json:"max_investment_amount"`
	MaxRiseAmount   *float64                 `json:"max_raise_amount"`
	MinRaiseAmount  *float64                 `json:"min_raise_amount"`

	StartDate   *time.Time `json:"startAt"`
	EndDate     *time.Time `json:"endAt"`
	IsScheduled bool       `json:"scheduled"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func ListingResponseFromModel(listing models.Listing) ListingResponse {
	response := ListingResponse{
		Id:     listing.Id,
		UserId: listing.UserId,

		MinInvestAmount: listing.MinInvestAmount,
		MaxInvestAmount: listing.MaxInvestAmount,
		MinRaiseAmount:  listing.MinToRaise,
		MaxRiseAmount:   listing.MaxToRaise,

		AssetName:   listing.AssetName,
		AssetSymbol: listing.AssetSymbol,
		Type:        listing.Type,
		AssetId:     listing.AssetId,

		Price: listing.PriceUSD,

		Currency: listing.Currency,
		Quantity: listing.Quantity,

		StartDate: listing.StartDate,
		IsActive:  listing.IsActive,
		EndDate:   listing.EndDate,
	}
	// it is scheduled if the listing is not active but has a start date later than now
	response.IsScheduled = !listing.IsActive && listing.StartDate.After(time.Now())
	response.CreatedAt = listing.CreatedAt.UTC()
	response.UpdatedAt = listing.UpdatedAt.UTC()
	return response
}
