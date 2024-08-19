package models

import "time"

type Trade struct {
	ID          string
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Quantity    float64
	CreatedAt   time.Time
}
