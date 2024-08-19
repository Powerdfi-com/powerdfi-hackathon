package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type AssetOwnerRepository interface {
	Update(assetOwner models.AssetOwner) error
	GetOwnerAsset(assetId, ownerId string) (models.AssetOwner, error)
}
