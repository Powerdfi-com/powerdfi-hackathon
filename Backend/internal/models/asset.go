package models

import (
	"bytes"
	"encoding/json"
	"io"
	"path"
	"time"
)

type Chain string

type AssetStatus string

var (
	VERIFIED_ASSET_STATUS   AssetStatus = "verified"
	UNVERIFIED_ASSET_STATUS AssetStatus = "unverified"
	REJECTED_ASSET_STATUS   AssetStatus = "rejected"
)

type Category struct {
	Id      int
	Name    string
	UrlSlug string
}

type Asset struct {
	Id            string
	TokenId       string
	Name          string
	Symbol        string
	CategoryId    *int
	CategoryName  *string
	Blockchain    Chain
	CreatorUserID string
	MetadataUrl   string

	URLs                 []string
	LegalDocumentURLs    []string
	IssuanceDocumentURLs []string

	Signatories []string
	Description string

	TotalSupply  int64
	SerialNumber *int64

	Properties     AssetProperties
	Favourites     int64
	Views          int64
	Status         AssetStatus
	IsFavourite    *bool
	IsVerified     bool
	IsListedByUser bool

	IsMinted       bool
	IsRejected     bool
	ExpirationDate time.Time
	ExecutionDate  *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FloorPrice     float64
}

// AssetProperties a slice of pairs in the form [[key, value], [key, value], [key, value]...]
type AssetProperties [][2]string

type assetMetadata struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Properties  AssetProperties `json:"properties"`
}

// GenerateMetadata returns the metadata of an Asset with the format defined in:
// https://eips.ethereum.org/EIPS/eip-1155#metadata
func (i *Asset) GenerateMetadata() (io.Reader, error) {
	metadata := assetMetadata{
		Name:        i.Name,
		Description: i.Description,
		// Image:       i.ImageUrl,
		Properties: i.Properties,
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// ImagePath parses an Asset image storage path in the form:
// "metadata/contract_address/Asset_id"
func (i *Asset) ImagePath() string {
	return path.Join("assets", i.Id)
}

// MetadataPath parses an Asset metadata storage path in the form:
// "metadata/contract_address/Asset_id"
func (i *Asset) MetadataPath() string {
	return path.Join("metadata", i.Id) + ".json"
}
