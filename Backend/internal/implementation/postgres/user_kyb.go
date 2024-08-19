package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
)

type userKybImplementation struct {
	Db *sql.DB
}

func NewUserKybImplementation(db *sql.DB) repository.UserKYBRepository {
	return userKybImplementation{Db: db}
}

func (u userKybImplementation) Create(userkyb models.UserKyB) error {

	stmt := `
	INSERT INTO users_kyb
	(
	user_id, 
	certificate_of_inc,
	platform,
	reference_id, 
	status,
	comment,
	company_name,
	company_location,
	company_address,
	company_reg_no
	)
	VALUES($1, $2, $3, $4 , $5 , $6,$7,$8,$9,$10);

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(
		ctx,
		stmt,
		userkyb.UserId,
		userkyb.CertificateOfInc,
		userkyb.Platform,
		userkyb.ReferenceId,
		userkyb.Status,
		userkyb.Comment,
		userkyb.CompanyName,
		userkyb.CompanyLocation,
		userkyb.CompanyAddress,
		userkyb.CompanyRegNo,
	)
	if err != nil {
		return err
	}

	return nil

}
func (u userKybImplementation) GetByUserId(userId string) (models.UserKyB, error) {
	stmt := `

SELECT
	user_id, 
	certificate_of_inc,
	platform,
	reference_id,
	status,
	comment,
 	company_name,
	company_location,
	company_address,
	company_reg_no
FROM users_kyb
WHERE user_id = $1;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userKyb := models.UserKyB{}
	err := u.Db.QueryRowContext(ctx, stmt, userId).Scan(
		&userKyb.UserId,
		&userKyb.CertificateOfInc,
		&userKyb.Platform,
		&userKyb.ReferenceId,
		&userKyb.Status,
		&userKyb.Comment,
		&userKyb.CompanyName,
		&userKyb.CompanyLocation,
		&userKyb.CompanyAddress,
		&userKyb.CompanyRegNo,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.UserKyB{}, repository.ErrRecordNotFound

		default:
			return models.UserKyB{}, err
		}
	}

	return userKyb, nil
}
func (u userKybImplementation) GetByReferenceId(referenceId string) (models.UserKyB, error) {
	stmt := `

SELECT
	user_id, 
	certificate_of_inc,
	platform,
	reference_id,
	status,
	comment,
	company_name,
	company_location,
	company_address,
	company_reg_no
FROM users_kyb
WHERE reference_id = $1;

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userKyb := models.UserKyB{}
	err := u.Db.QueryRowContext(ctx, stmt, referenceId).Scan(
		&userKyb.UserId,
		&userKyb.CertificateOfInc,
		&userKyb.Platform,
		&userKyb.ReferenceId,
		&userKyb.Status,
		&userKyb.Comment,
		&userKyb.CompanyName,
		&userKyb.CompanyLocation,
		&userKyb.CompanyAddress,
		&userKyb.CompanyRegNo,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.UserKyB{}, repository.ErrRecordNotFound

		default:
			return models.UserKyB{}, err
		}
	}

	return userKyb, nil
}
func (u userKybImplementation) Update(userKyb models.UserKyB) error {
	stmt := `

UPDATE users_kyb
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
		userKyb.ReferenceId,
		userKyb.Status,
		userKyb.Comment,
		userKyb.Platform,
		userKyb.ReferenceId,
		userKyb.CertificateOfInc)
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
func (u userKybImplementation) Delete(userId string) error {
	stmt := `

DELETE FROM users_kyb
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
