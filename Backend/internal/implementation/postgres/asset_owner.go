package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
)

type assetOwnerImplementation struct {
	Db *sql.DB
}

func NewAssetOwnerImplementation(db *sql.DB) repository.AssetOwnerRepository {
	return assetOwnerImplementation{
		Db: db,
	}
}

func (ao assetOwnerImplementation) Update(assetOwner models.AssetOwner) error {

	stmt := `
	INSERT INTO asset_owners (
    user_id, asset_id, serial_numbers
) VALUES (
    $1, $2, $3
) ON CONFLICT (user_id, asset_id) DO UPDATE SET
    serial_numbers = EXCLUDED.serial_numbers;


    `

	serialNumbers, err := json.Marshal(&assetOwner.SerialNumbers)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = ao.Db.ExecContext(
		ctx,
		stmt,
		assetOwner.UserId,
		assetOwner.AssetId,
		string(serialNumbers),
	)
	if err != nil {
		return err
	}

	return nil
}
func (ao assetOwnerImplementation) GetOwnerAsset(assetId, ownerId string) (models.AssetOwner, error) {

	query := `
	SELECT 
	user_id, 
	asset_id,
	serial_numbers
	FROM  asset_owners

	WHERE user_id = $1
	AND asset_id=$2;

    `

	assetOwner := models.AssetOwner{
		UserId:        ownerId,
		AssetId:       assetId,
		SerialNumbers: make([]int64, 0),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	serialNumbers := ""
	err := ao.Db.QueryRowContext(ctx, query, ownerId, assetId).Scan(
		&assetOwner.UserId,
		&assetOwner.AssetId,
		&serialNumbers,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return assetOwner, repository.ErrRecordNotFound

		default:
			return assetOwner, err
		}

	}

	err = json.NewDecoder(strings.NewReader(serialNumbers)).Decode(&assetOwner.SerialNumbers)
	if err != nil {
		return assetOwner, err
	}
	return assetOwner, nil
}
