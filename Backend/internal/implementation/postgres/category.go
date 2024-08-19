package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
)

type categoryImplementation struct {
	Db *sql.DB
}

func NewCategoryImplementation(db *sql.DB) repository.CategoryRepository {
	return categoryImplementation{Db: db}
}

// List returns a list of all categories
func (c categoryImplementation) List() ([]models.Category, error) {
	stmt := `
	SELECT id, name, url_slug
	FROM categories
	`

	categories := []models.Category{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := c.Db.QueryContext(ctx, stmt)

	if err != nil {

		return categories, err
	}

	defer rows.Close()

	for rows.Next() {
		cat := models.Category{}
		rows.Scan(
			&cat.Id,
			&cat.Name,
			&cat.UrlSlug,
		)

		categories = append(categories, cat)

	}

	return categories, nil
}

// ValidateSlug checks if the given URL slug has a matching category and returns that.
func (c categoryImplementation) GetById(id int) (models.Category, error) {
	stmt := `
	SELECT id, name, url_slug
	FROM categories
	WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cat := models.Category{}
	err := c.Db.QueryRowContext(ctx, stmt, id).Scan(
		&cat.Id,
		&cat.Name,
		&cat.UrlSlug,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Category{}, repository.ErrRecordNotFound

		default:
			return models.Category{}, err
		}
	}

	return cat, nil
}
func (c categoryImplementation) ValidateSlug(slug string) (models.Category, error) {
	stmt := `
	SELECT id, name, url_slug
	FROM categories
	WHERE url_slug = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cat := models.Category{}
	err := c.Db.QueryRowContext(ctx, stmt, slug).Scan(
		&cat.Id,
		&cat.Name,
		&cat.UrlSlug,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Category{}, repository.ErrRecordNotFound

		default:
			return models.Category{}, err
		}
	}

	return cat, nil
}
