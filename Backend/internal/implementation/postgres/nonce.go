package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
)

type nonce struct {
	Db *sql.DB
}

func NewNonceImplementation(db *sql.DB) repository.NonceRepository {
	return nonce{
		Db: db,
	}
}

func (n nonce) Create(address string, message string) error {
	stmt := `
	INSERT INTO nonces(user_address, message, action) 
	VALUES ($1, $2, $3);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := n.Db.ExecContext(ctx, stmt, address, message, models.NonceActionCreate)
	if err != nil {
		return err
	}

	return nil
}

func (n nonce) Get(address string) (models.Nonce, error) {
	stmt := `
	SELECT message, action FROM nonces
	WHERE user_address = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	nonce := models.Nonce{}
	err := n.Db.QueryRowContext(ctx, stmt, address).Scan(&nonce.Message, &nonce.Action)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Nonce{}, repository.ErrRecordNotFound

		default:
			return models.Nonce{}, err
		}
	}

	return nonce, nil
}

func (n nonce) Update(address string, message string) error {
	stmt := `
	UPDATE nonces
	SET message = $1, action = $2
	WHERE user_address = $3;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// set the nonce action to "update" so future queries don't create
	// an account with the address
	_, err := n.Db.ExecContext(ctx, stmt, message, models.NonceActionUpdate, address)
	if err != nil {
		return err
	}

	return nil
}

// DeleteExpired invalidates the nonces of addresses with uncreated accounts
// by deleting those created before a given duration.
func (n nonce) DeleteExpired() error {
	stmt := `
	DELETE FROM nonces
	WHERE action = $1 AND created_at < now() - $2::INTERVAL;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := n.Db.ExecContext(ctx, stmt, models.NonceActionCreate, models.NonceInactiveAddressDuration)
	if err != nil {
		return err
	}

	return nil
}
