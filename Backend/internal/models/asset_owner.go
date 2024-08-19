package models

type AssetOwner struct {
	UserId        string
	AssetId       string
	SerialNumbers []int64
}
