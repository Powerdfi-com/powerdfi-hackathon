package models

import (
	"time"
)

type ActivityAction string

const (
	ACTIVITY_ACTION_SALE = "sale"
	ACTIVITY_ACTION_MINT = "mint"
	// ACTIVITY_ACTION_LIST       = "list"
	// ACTIVITY_ACTION_DELIST     = "delist"
	// ACTIVITY_ACTION_BID        = "bid"
	// ACTIVITY_ACTION_BID_CANCEL = "bid_cancel"
)

type Activity struct {
	Id        string
	Action    ActivityAction
	AssetId   string
	AssetName *string
	Price     float64
	// Price      *big.Float
	Currency string
	// Quantity   *big.Int
	Quantity   int64
	FromUserId *string
	ToUserId   string

	Blockchain string

	// ListingId       *string
	// BlockNumber     *big.Int
	// TransactionHash string

	CreatedAt  time.Time
	OccurredAt time.Time
}
