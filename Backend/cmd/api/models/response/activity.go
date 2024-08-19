package response

import (
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type ActivityResponse struct {
	Id         string    `json:"id"`
	Action     string    `json:"action"`
	AssetID    string    `json:"assetId"`
	AssetName  *string   `json:"assetName"`
	Price      float64   `json:"price,omitempty"`
	Currency   string    `json:"currency,omitempty"`
	Quantity   int64     `json:"quantity"`
	FromUserId *string   `json:"fromUserId,omitempty"`
	ToUserId   string    `json:"toUserId"`
	Blockchain string    `json:"blockchain,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

func ActivityResponseFromModel(activity models.Activity) ActivityResponse {

	return ActivityResponse{
		Id:         activity.Id,
		AssetID:    activity.AssetId,
		Action:     string(activity.Action),
		AssetName:  activity.AssetName,
		Price:      activity.Price,
		Currency:   activity.Currency,
		Quantity:   activity.Quantity,
		FromUserId: activity.FromUserId,
		ToUserId:   activity.ToUserId,
		Blockchain: activity.Blockchain,
		CreatedAt:  activity.CreatedAt,
	}
}
