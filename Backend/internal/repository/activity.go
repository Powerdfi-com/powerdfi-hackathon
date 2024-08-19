package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type ActivityRepository interface {
	Add(activity models.Activity) (models.Activity, error)
	ListByAssetID(assetID string, filter models.Filter) ([]models.Activity, error)
	ListForUser(userId string, filter models.Filter) ([]models.Activity, int, error)
	// RecordListingOrDeListing(activity models.Activity) error
	// RecordBidOrBidCancel(activity models.Activity) error
	// GetItemActivity(tokenId, contractAddress string, filter models.Filter) ([]models.Activity, error)
	// GetCollectionActivity(contractAddress string, filter models.Filter) ([]models.Activity, error)
	// GetUserActivity(userAddress string, filter models.Filter) ([]models.Activity, error)
	// UpdateToAddress(toAddress string, transactionHash string) error
}
