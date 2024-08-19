package postgres

import (
	"database/sql"

	"github.com/Powerdfi-com/Backend/internal/repository"
)

func NewRepositories(db *sql.DB) repository.Repositories {
	return repository.Repositories{
		User:         NewUserImplementation(db),
		Nonce:        NewNonceImplementation(db),
		Asset:        NewAssetImplementation(db),
		Listing:      NewListingImplementation(db),
		UserKyc:      NewUserKycImplementation(db),
		Stats:        NewStatImplementation(db),
		Category:     NewCategoryImplementation(db),
		AssetOwner:   NewAssetOwnerImplementation(db),
		Order:        NewOrderImplementation(db),
		Trade:        NewTradeImplementation(db),
		Activity:     NewActivityImplementation(db),
		Admin:        NewAdminImplementation(db),
		Notification: NewNotificationImplementation(db),
		UserKYB:      NewUserKybImplementation(db),
	}
}
