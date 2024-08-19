package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type ListingRepository interface {
	Create(listing models.Listing) (models.Listing, error)
	GetById(id string) (models.Listing, error)
	Cancel(listingId string) error
}
