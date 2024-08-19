package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type AdminRepository interface {
	Create(admin models.Admin) (models.Admin, error)
	FindByEmail(email string) (models.Admin, error)
	UpdatePasswordHash(id string, passwordHash []byte) error
	List(filter models.Filter) ([]models.Admin, int, error)
	FindByID(id string) (models.Admin, error)
	Update(admin models.Admin) error
}
