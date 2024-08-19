package models

import "time"

type AssetTransfer struct {
	ID             string
	AssetId        string
	SenderUserID   string
	ReceiverUserID string
	TransactionID  string
	CreatedAt      time.Time
}
