package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type AssetRepository interface {
	List(filter models.Filter) ([]models.Asset, error)
	Create(asset models.Asset) (models.Asset, error)
	GetById(id string) (models.Asset, error)
	Update(asset models.Asset) (models.Asset, error)
	AddFavourite(userId string, assetId string) error
	RemoveFavourite(userId string, assetId string) error
	IsFavourite(userId string, assetId string) (bool, error)
	CountFavourites(assetId string) (int64, error)
	AddView(userId string, assetId string) error
	IsViewed(userId string, assetId string) (bool, error)
	CountViews(assetId string) (int64, error)
	GetListings(assetId string, filter models.Filter) ([]models.Listing, error)
	// GetListings(assetId string, filter models.Filter) ([]models.Listing, int, error)
	ListRecommended(asset models.Asset, filter models.Filter) ([]models.AssetStat, int, error)
	GetOrders(assetId string, status *models.OrderStatus, orderType *models.OrderType, filter models.Filter) ([]models.Order, int, error)

	GetApprovedUnmintedAssets(filter models.Filter) ([]models.Asset, error)

	UpdateMintStatus(assetId string) error
	// Approve(assetId string) error
	// Reject(assetId string) error

	UpdateStatus(assetId string, status models.AssetStatus) error
	IsListed(userId string, assetId string) (bool, error)

	// GetCreator(contractAddress string, tokenId string) (models.User, error)
	// Delete(id string) error
}
