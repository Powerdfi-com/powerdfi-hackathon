package response

import (
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
)

type OrderResponse struct {
	Id      string `json:"id"`
	AssetId string `json:"assetId"`
	UserId  string `json:"userId"`

	Type      models.OrderType   `json:"type"`
	Kind      models.OrderKind   `json:"kind"`
	Price     float64            `json:"price"`
	Quantity  int64              `json:"quantity"`
	Status    models.OrderStatus `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func OrderResponseFromModel(order models.Order) OrderResponse {
	return OrderResponse{
		Id:        order.Id,
		AssetId:   order.AssetId,
		UserId:    order.UserId,
		Type:      order.Type,
		Kind:      order.Kind,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Price:     order.Price,
		Quantity:  order.Quantity,
	}
}
