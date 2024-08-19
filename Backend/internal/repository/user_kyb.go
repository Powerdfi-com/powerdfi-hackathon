package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type UserKYBRepository interface {
	Create(userkyb models.UserKyB) error
	GetByUserId(userId string) (models.UserKyB, error)
	GetByReferenceId(referenceId string) (models.UserKyB, error)
	Update(userKyB models.UserKyB) error
	Delete(userId string) error
}
