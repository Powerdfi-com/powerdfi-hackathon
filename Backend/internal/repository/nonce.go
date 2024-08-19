package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type NonceRepository interface {
	Create(address string, message string) error
	Get(address string) (models.Nonce, error)
	Update(address string, message string) error
	DeleteExpired() error
}
