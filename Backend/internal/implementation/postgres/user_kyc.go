package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
)

type userKycImplementation struct {
	Db *sql.DB
}

func NewUserKycImplementation(db *sql.DB) repository.UserKycRepository {
	return userKycImplementation{Db: db}
}

func (u userKycImplementation) Create(userkyc models.UserKyc) error {

	stmt := `
	INSERT INTO users_kyc
(user_id, url, platform,reference_id, status,
 comment)
VALUES($1, $2, $3, $4 , $5 , $6);

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(
		ctx,
		stmt,
		userkyc.UserId,
		userkyc.URL,
		userkyc.Platform,
		userkyc.ReferenceId,
		userkyc.Status,
		userkyc.Comment,
	)
	if err != nil {
		return err
	}

	return nil

}
func (u userKycImplementation) GetByUserId(userId string) (models.UserKyc, error) {
	stmt := `

SELECT
 user_id, 
 url, 
 platform,
 reference_id,
 status,
 comment
FROM users_kyc
WHERE user_id = $1;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userKyc := models.UserKyc{}
	err := u.Db.QueryRowContext(ctx, stmt, userId).Scan(
		&userKyc.UserId,
		&userKyc.URL,
		&userKyc.Platform,
		&userKyc.ReferenceId,
		&userKyc.Status,
		&userKyc.Comment,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.UserKyc{}, repository.ErrRecordNotFound

		default:
			return models.UserKyc{}, err
		}
	}

	return userKyc, nil
}
func (u userKycImplementation) GetByReferenceId(referenceId string) (models.UserKyc, error) {
	stmt := `

SELECT
 user_id, 
 url, 
 platform,
 reference_id,
 status,
 comment
FROM users_kyc
WHERE reference_id = $1;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userKyc := models.UserKyc{}
	err := u.Db.QueryRowContext(ctx, stmt, referenceId).Scan(
		&userKyc.UserId,
		&userKyc.URL,
		&userKyc.Platform,
		&userKyc.ReferenceId,
		&userKyc.Status,
		&userKyc.Comment,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.UserKyc{}, repository.ErrRecordNotFound

		default:
			return models.UserKyc{}, err
		}
	}

	return userKyc, nil
}
func (u userKycImplementation) Update(userKyc models.UserKyc) error {
	stmt := `

UPDATE users_kyc
SET 
 status=$2, 
 "comment"=$3,
 platform=$4,
 reference_id=$5,
 url=$6
WHERE reference_id = $1;


	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(ctx, stmt,
		userKyc.ReferenceId,
		userKyc.Status,
		userKyc.Comment,
		userKyc.Platform,
		userKyc.ReferenceId,
		userKyc.URL)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrRecordNotFound

		default:
			return err
		}
	}

	return nil
}
func (u userKycImplementation) Delete(userId string) error {
	stmt := `

DELETE FROM users_kyc
WHERE user_id = $1;


	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(ctx, stmt, userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrRecordNotFound

		default:
			return err
		}
	}

	return nil
}
