package postgres

import (
	"database/sql"

	"github.com/Powerdfi-com/Backend/internal/repository"
)

type tradeImpl struct {
	Db *sql.DB
}

func NewTradeImplementation(db *sql.DB) repository.TradeRepository {
	return tradeImpl{
		Db: db,
	}
}
