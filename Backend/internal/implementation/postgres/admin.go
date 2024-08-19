package postgres

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
)

const (
	duplicateAdminEmail = "admins_email_key"
)

type adminImplementation struct {
	Db *sql.DB
}

func NewAdminImplementation(db *sql.DB) repository.AdminRepository {
	return adminImplementation{Db: db}
}

func (u adminImplementation) Create(admin models.Admin) (models.Admin, error) {
	// generate an ID for the user if it doesn't have one
	if admin.Id == "" {
		admin.Id = uuid.NewString()
	}

	stmt := `
	INSERT INTO admins(id, email, password_hash, name)
	VALUES ($1, $2, $3,$4)
	RETURNING 
	    id,
	    role_mask, 
	    created_at, 
	    updated_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newAdmin := admin
	err := u.Db.QueryRowContext(ctx, stmt, admin.Id, admin.Email, admin.PasswordHash, admin.Name).Scan(
		&newAdmin.Id,
		&newAdmin.RoleMask,
		&newAdmin.CreatedAt,
		&newAdmin.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateAdminEmail):
			return models.Admin{}, repository.ErrDuplicateDetails

		default:
			return models.Admin{}, err
		}
	}

	return newAdmin, err
}

func (u adminImplementation) FindByEmail(email string) (models.Admin, error) {
	stmt := `
	SELECT 
		id,
		email,
		password_hash,
		name,
		role_mask,
		created_at,
		updated_at
	FROM 
		admins
	WHERE 
		email = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	admin := models.Admin{}
	err := u.Db.QueryRowContext(ctx, stmt, email).Scan(
		&admin.Id,
		&admin.Email,
		&admin.PasswordHash,
		&admin.Name,
		&admin.RoleMask,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, repository.ErrRecordNotFound
		}
		return models.Admin{}, err
	}

	return admin, nil
}
func (u adminImplementation) FindByID(id string) (models.Admin, error) {
	stmt := `
	SELECT 
		id,
		email,
		password_hash,
		name,
		role_mask,
		created_at,
		updated_at
	FROM 
		admins
	WHERE 
		id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	admin := models.Admin{}
	err := u.Db.QueryRowContext(ctx, stmt, id).Scan(
		&admin.Id,
		&admin.Email,
		&admin.PasswordHash,
		&admin.Name,
		&admin.RoleMask,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, repository.ErrRecordNotFound
		}
		return models.Admin{}, err
	}

	return admin, nil
}
func (u adminImplementation) UpdatePasswordHash(id string, passwordHash []byte) error {
	stmt := `
	UPDATE admins
	SET password_hash = $1
	WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := u.Db.ExecContext(ctx, stmt, passwordHash, id)
	if err != nil {
		return err
	}

	return nil
}
func (u adminImplementation) Update(admin models.Admin) error {
	stmt := `
	UPDATE admins
	SET 
	email=$2,
	role_mask= $3,
	password_hash = $4
	WHERE id = $1
	RETURNING 
	    updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, stmt, admin.Id, admin.Email, admin.RoleMask, admin.PasswordHash).Scan(
		&admin.UpdatedAt,
	)
	if err != nil {
		switch {
		case
			strings.Contains(err.Error(), duplicateAdminEmail):
			return repository.ErrDuplicateDetails

		default:
			return err
		}
	}

	return nil
}

func (u adminImplementation) List(filter models.Filter) ([]models.Admin, int, error) {
	var totalCount int

	stmt := `
	SELECT 
		id,
		email,
		password_hash,
		name,
		role_mask,
		created_at,
		updated_at
	FROM 
		admins
		LIMIT $1 OFFSET $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	admins := []models.Admin{}

	rows, err := u.Db.QueryContext(
		ctx,
		stmt,
		filter.Limit,
		filter.Offset(),
	)

	if err != nil {
		return admins, totalCount, err
	}
	defer rows.Close()

	totalCountQuery := `
	SELECT COUNT(*)
	FROM admins;
	`
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = u.Db.QueryRowContext(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return nil, totalCount, err
	}
	for rows.Next() {
		admin := models.Admin{}
		err := rows.Scan(
			&admin.Id,
			&admin.Email,
			&admin.PasswordHash,
			&admin.Name,
			&admin.RoleMask,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)

		if err != nil {
			log.Printf("err, retrieving admin")
			continue

		}

		admins = append(admins, admin)
	}
	if err = rows.Err(); err != nil {
		return admins, totalCount, err
	}

	return admins, totalCount, nil
}
