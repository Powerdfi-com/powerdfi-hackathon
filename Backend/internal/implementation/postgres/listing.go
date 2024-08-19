package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
)

type listingImpl struct {
	Db *sql.DB
}

func NewListingImplementation(db *sql.DB) repository.ListingRepository {
	return listingImpl{
		Db: db,
	}
}

func (l listingImpl) Create(listing models.Listing) (models.Listing, error) {
	err := l.checkActivenessForItem(listing.UserId, listing.AssetId)
	if err != nil {
		return models.Listing{}, err
	}

	if listing.Id == "" {
		listing.Id = uuid.NewString()
	}

	stmt := `
		INSERT INTO listings
	(
		id, 
		type, 
		user_id, 
		asset_id,
		price, 
		min_invest_amount,
		max_invest_amount,
		min__raise_amount,
		max_raise_amount,
		currency, 
		quantity, 
		start_date, 
		end_date,
		is_active
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 
        CASE WHEN CURRENT_DATE >= Date($12) THEN true ELSE false END)
	RETURNING 
	    id,
	   created_at, 
	   updated_at,
	   is_active;
	`

	currency, err := json.Marshal(listing.Currency)
	if err != nil {
		return models.Listing{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newListing := listing
	err = l.Db.QueryRowContext(ctx, stmt,
		listing.Id,
		&listing.Type,
		&listing.UserId,
		&listing.AssetId,
		&listing.PriceUSD,
		&listing.MinInvestAmount,
		&listing.MaxInvestAmount,
		&listing.MinToRaise,
		&listing.MaxToRaise,
		currency,
		&listing.Quantity,
		&listing.StartDate,
		&listing.EndDate,
	).Scan(
		&newListing.Id,
		&newListing.CreatedAt,
		&newListing.UpdatedAt,
		&newListing.IsActive,
	)
	if err != nil {
		// TODO: handle duplicate error if any field is duplicate
		return models.Listing{}, err
	}

	return newListing, err
}

func (l listingImpl) GetById(id string) (models.Listing, error) {
	stmt := `
	SELECT 
	   id,
      "type",
       user_id, 
       asset_id,
       price,
       min_invest_amount, 
       max_invest_amount, 
       min__raise_amount, 
       max_raise_amount,
       currency, 
       quantity, 
       start_date, 
       end_date, 
       is_active,
       is_cancelled, 
       created_at, 
       updated_at  
	FROM listings AS l 	
	WHERE l.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	currency := ""
	listing := models.Listing{}
	err := l.Db.QueryRowContext(ctx, stmt, id).Scan(
		&listing.Id,
		&listing.Type,
		&listing.UserId,
		&listing.AssetId,
		&listing.PriceUSD,
		&listing.MinInvestAmount,
		&listing.MaxInvestAmount,
		&listing.MinToRaise,
		&listing.MaxToRaise,
		&currency,
		&listing.Quantity,
		&listing.StartDate,
		&listing.EndDate,
		&listing.IsActive,
		&listing.IsCancelled,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Listing{}, repository.ErrRecordNotFound

		default:
			return models.Listing{}, err
		}
	}

	err = json.Unmarshal([]byte(currency), &listing.Currency)
	if err != nil {
		return models.Listing{}, err
	}

	return listing, nil

}
func (l listingImpl) Cancel(id string) error {
	stmt := `
	UPDATE listings
	SET is_active = false,
	is_cancelled = true
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := l.Db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (l listingImpl) checkActivenessForItem(userId string, assetId string) error {
	stmt := `
	SELECT is_active
	FROM listings
	WHERE 
		user_id = $1
		AND asset_id = $2
		AND (
			(is_active = true) 
			OR (NOW() <= start_date AND is_cancelled = false)  -- Added condition for non-canceled items
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var isActive bool
	err := l.Db.QueryRowContext(ctx, stmt, userId, assetId).Scan(&isActive)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// return a success if no active listing is found
		return nil

	case isActive:
		// return a specific error if an active listing is found
		return repository.ErrActiveListing

	default:
		return err
	}
}
