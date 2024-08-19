package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
)

type activityImpl struct {
	Db *sql.DB
}

func NewActivityImplementation(db *sql.DB) repository.ActivityRepository {
	return activityImpl{
		Db: db,
	}
}

func (a activityImpl) Add(activity models.Activity) (models.Activity, error) {
	// generate an ID for the user if it doesn't have one
	if activity.Id == "" {
		activity.Id = uuid.NewString()
	}

	stmt := `

	INSERT INTO activities
(
	id, 
	"action", 
	asset_id, 
	from_user_id, 
	to_user_id, 
	price, 
	currency, 
	quantity
	)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING 
	    id,
	   created_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newActivity := activity

	err := a.Db.QueryRowContext(ctx, stmt,
		activity.Id,
		activity.Action,
		activity.AssetId,
		activity.FromUserId,
		activity.ToUserId,
		activity.Price,
		activity.Currency,
		activity.Quantity,
	).Scan(
		&newActivity.Id,
		&newActivity.CreatedAt,
	)
	if err != nil {

		return models.Activity{}, repository.ErrDuplicateDetails
	}

	return newActivity, err
}

func (a activityImpl) ListByAssetID(assetID string, filter models.Filter) ([]models.Activity, error) {

	stmt := `
	SELECT 
	ac.id, 
	ac."action", 
	ac.asset_id,
	a."name" ,
	ac.from_user_id, 
	ac.to_user_id, 
	ac.price, 
	ac.currency, 
	ac.quantity, 
	ac.created_at
	FROM activities as ac
	JOIN assets AS a ON ac.asset_id = a.id
	WHERE asset_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := a.Db.QueryContext(ctx, stmt, assetID, filter.Limit, filter.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activity

	for rows.Next() {
		var activity models.Activity
		err := rows.Scan(
			&activity.Id,
			&activity.Action,
			&activity.AssetId,
			&activity.AssetName,
			&activity.FromUserId,
			&activity.ToUserId,
			&activity.Price,
			&activity.Currency,
			&activity.Quantity,
			&activity.CreatedAt,
		)
		if err != nil {
			return activities, err
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return activities, err
	}

	return activities, nil
}
func (a activityImpl) ListForUser(userId string, filter models.Filter) ([]models.Activity, int, error) {
	var totalCount int

	stmt := `
	SELECT 
	ac.id, 
	ac."action", 
	ac.asset_id,
	a."name" ,
	ac.from_user_id, 
	ac.to_user_id, 
	ac.price, 
	ac.currency, 
	ac.quantity, 
	ac.created_at
	FROM activities as ac
	JOIN assets AS a ON ac.asset_id = a.id
	WHERE ac.from_user_id = $1
	OR  ac.to_user_id=$1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := a.Db.QueryContext(ctx, stmt, userId, filter.Limit, filter.Offset())
	if err != nil {
		return nil, totalCount, err
	}
	defer rows.Close()

	var activities []models.Activity

	totalCountQuery := `
	SELECT COUNT(*)
   FROM activities as ac
	WHERE ac.from_user_id = $1
	OR ac.to_user_id = $1;
 `

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = a.Db.QueryRowContext(ctx, totalCountQuery, userId).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		var activity models.Activity
		err := rows.Scan(
			&activity.Id,
			&activity.Action,
			&activity.AssetId,
			&activity.AssetName,
			&activity.FromUserId,
			&activity.ToUserId,
			&activity.Price,
			&activity.Currency,
			&activity.Quantity,
			&activity.CreatedAt,
		)
		if err != nil {
			return activities, totalCount, err
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return activities, totalCount, err
	}

	return activities, totalCount, nil
}
