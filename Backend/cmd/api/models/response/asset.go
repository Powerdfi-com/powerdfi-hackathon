package response

import (
	"time"

	"github.com/Powerdfi-com/Backend/internal"
	"github.com/Powerdfi-com/Backend/internal/models"
)

type AssetResponse struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	Symbol        string                 `json:"symbol"`
	Description   string                 `json:"description"`
	CategoryName  *string                `json:"category,omitempty"`
	CategoryId    *int                   `json:"categoryId,omitempty"`
	CreatorId     string                 `json:"creatorId"`
	Blockchain    string                 `json:"blockchain"`
	TokenId       string                 `json:"tokenId"`
	MetadataUrl   string                 `json:"metadataUrl"`
	TotalSupply   int64                  `json:"totalSupply"`
	Properties    models.AssetProperties `json:"properties"`
	URLs          []string               `json:"urls"`
	TokenStandard string                 `json:"tokenStandard"`
	FloorPrice    float64                `json:"floorPrice"`
	Status        models.AssetStatus     `json:"status"`

	LegalDocumentURLs    []string `json:"legalDocumentUrls"`
	IssuanceDocumentURLs []string `json:"issuanceDocumentsUrls"`
	Favourites           int64    `json:"favourites"`
	Views                int64    `json:"views"`
	IsFavourite          *bool    `json:"isFavourite,omitempty"`
	// IsVerified           bool     `json:"isVerified"`
	// IsApproved           bool      `json:"isApproved"`
	// IsRejected           bool      `json:"isRejected"`
	IsListedByUser bool      `json:"isListedByUser"`
	IsMinted       bool      `json:"isMinted"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CategoryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
type ChainResponse struct {
	Id        models.Chain `json:"id"`
	Name      string       `json:"name"`
	Logo      string       `json:"logo"`
	IsEnabled bool         `json:"isEnabled"`
}

func AssetResponseFromModel(asset models.Asset) AssetResponse {
	response := AssetResponse{
		Id:                   asset.Id,
		Name:                 asset.Name,
		Symbol:               asset.Symbol,
		Blockchain:           internal.ChainNameMap[asset.Blockchain],
		TokenId:              asset.TokenId,
		CreatorId:            asset.CreatorUserID,
		MetadataUrl:          asset.MetadataUrl,
		Description:          asset.Description,
		TotalSupply:          asset.TotalSupply,
		Properties:           asset.Properties,
		LegalDocumentURLs:    asset.LegalDocumentURLs,
		URLs:                 asset.URLs,
		IssuanceDocumentURLs: asset.IssuanceDocumentURLs,
		Favourites:           asset.Favourites,
		CategoryName:         asset.CategoryName,
		CategoryId:           asset.CategoryId,
		Views:                asset.Views,
		IsFavourite:          asset.IsFavourite,
		IsListedByUser:       asset.IsListedByUser,
		TokenStandard:        "ERC115",
		FloorPrice:           asset.FloorPrice,
		Status:               asset.Status,
		// show badge if asset has both been minted and verified
		// IsVerified: asset.IsVerified && asset.IsMinted,
		IsMinted: asset.IsMinted,
		// IsApproved: asset.IsVerified,
		// IsRejected: asset.IsRejected,
	}

	response.CreatedAt = asset.CreatedAt.UTC()
	response.UpdatedAt = asset.UpdatedAt.UTC()
	return response
}
