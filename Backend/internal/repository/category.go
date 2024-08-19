package repository

import "github.com/Powerdfi-com/Backend/internal/models"

type CategoryRepository interface {
	List() ([]models.Category, error)
	ValidateSlug(slug string) (models.Category, error)
	GetById(id int) (models.Category, error)
}
