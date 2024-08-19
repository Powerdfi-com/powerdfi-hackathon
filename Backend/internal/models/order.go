package models

import "time"

type OrderType string
type OrderStatus string
type OrderKind string

var (
	ORDER_BUY_TYPE  OrderType = "buy"
	ORDER_SELL_TYPE OrderType = "sell"
)
var (
	ORDER_OPEN_STATUS             OrderStatus = "open"
	ORDER_CANCELLED_STATUS        OrderStatus = "cancelled"
	ORDER_FILLED_STATUS           OrderStatus = "filled"
	ORDER_PARTIALLY_FILLED_STATUS OrderStatus = "partial"
	ORDER_FAILED_STATUS           OrderStatus = "failed"
)
var (
	ORDER_LIMIT_KIND  OrderKind = "limit"
	ORDER_MARKET_KIND OrderKind = "market"
)

type Order struct {
	Id          string
	UserId      string
	AssetId     string
	Type        OrderType
	Kind        OrderKind
	Price       float64
	Quantity    int64
	InitialQty  int64
	FilledPrice int64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
