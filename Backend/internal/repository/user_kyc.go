package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type UserKycRepository interface {
	Create(userkycLink models.UserKyc) error
	GetByUserId(userId string) (models.UserKyc, error)
	GetByReferenceId(referenceId string) (models.UserKyc, error)
	Update(userKyc models.UserKyc) error
	Delete(userId string) error
}
