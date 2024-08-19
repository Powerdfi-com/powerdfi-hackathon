package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type OrderRepository interface {
	Create(order models.Order) (models.Order, error)
	Update(order models.Order) error
	GetById(id string) (models.Order, error)
	Cancel(id string) error
	GetUnFilledBuyOrders(filter models.Filter) ([]models.Order, error)
	FindMatchingSellOrder(order models.Order) (models.Order, error)
}
