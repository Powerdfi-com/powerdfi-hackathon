package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetByAddress(address string) (models.User, error)
	GetById(id string) (models.User, error)
	Update(user models.User) (models.User, error)
	Activate(id string) error
	Verify(id string) error
	GetListedAssets(id string, filter models.Filter) ([]models.Asset, error)
	ListCreatedAssets(id string, filter models.Filter) ([]models.Asset, error)
	GetUnListedAssets(id string, filter models.Filter) ([]models.Asset, error)
	ListOwnedAssets(id string, filter models.Filter) ([]models.Asset, error)
	HasOpenSellOrder(userId, assetId string) (bool, error)
	GetOrders(userId string, status *models.OrderStatus, orderType *models.OrderType, filter models.Filter) ([]models.Order, error)
	SetPrivateKey(userId string, encryptedPrivateKey []byte) error
}
