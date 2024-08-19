package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Powerdfi-com/Backend/internal/models"
	"github.com/Powerdfi-com/Backend/internal/repository"
	"github.com/google/uuid"
)

type orderImpl struct {
	Db *sql.DB
}

func NewOrderImplementation(db *sql.DB) repository.OrderRepository {
	return orderImpl{
		Db: db,
	}
}

func (o orderImpl) Create(order models.Order) (models.Order, error) {
	// generate an ID for the user if it doesn't have one
	if order.Id == "" {
		order.Id = uuid.NewString()
	}

	stmt := `
INSERT INTO public.orders
(
    id, 
 user_id, 
 asset_id, 
 "type",
 kind,
 status, 
 price, 
 quantity,
 inital_quantity
 )
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING 
	    id,
	   created_at, 
	   updated_at;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newOrder := order

	err := o.Db.QueryRowContext(ctx, stmt,
		order.Id,
		order.UserId,
		order.AssetId,
		order.Type,
		order.Kind,
		order.Status,
		order.Price,
		order.Quantity,
		order.Quantity,
	).Scan(
		&newOrder.Id,
		&newOrder.CreatedAt,
		&newOrder.UpdatedAt,
	)
	if err != nil {
		// TODO: handle duplicate error if any field is duplicate
		switch {
		case strings.Contains(err.Error(), duplicateAddress):
			return models.Order{}, repository.ErrDuplicateDetails

		default:
			return models.Order{}, err
		}
	}

	return newOrder, err
}

func (o orderImpl) GetById(id string) (models.Order, error) {
	stmt := `
	SELECT 
	id,
	user_id, 
	asset_id, 
	"type", 
	kind, 
	status, 
	price, 
	quantity, 
	inital_quantity,
	created_at, 
	updated_at
	FROM orders
WHERE id=$1

	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	order := models.Order{}

	err := o.Db.QueryRowContext(ctx, stmt, id).Scan(
		&order.Id,
		&order.UserId,
		&order.AssetId,
		&order.Type,
		&order.Kind,
		&order.Status,
		&order.Price,
		&order.Quantity,
		&order.InitialQty,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return order, repository.ErrRecordNotFound

		default:
			return order, err
		}
	}

	return order, nil
}

func (o orderImpl) Cancel(id string) error {
	stmt := `
	UPDATE orders
	SET status = $2
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := o.Db.ExecContext(ctx, stmt, id, models.ORDER_CANCELLED_STATUS)
	if err != nil {
		return err
	}

	return nil
}
func (o orderImpl) Update(order models.Order) error {
	stmt := `
	UPDATE orders
	SET quantity = $2, status = $3, updated_at = $4
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := o.Db.ExecContext(ctx, stmt, order.Id, order.Quantity, order.Status, time.Now())
	if err != nil {
		return err
	}

	return nil
}
func (o orderImpl) GetUnFilledBuyOrders(filter models.Filter) ([]models.Order, error) {
	stmt := `
	SELECT 
	id, 
	user_id, 
	asset_id,
	"type", 
	kind,
	status,
	price, 
	quantity, 
	inital_quantity, 
	created_at, 
	updated_at
	FROM orders
	WHERE status IN ($1, $2) 
	 AND "type" = $3
  	  AND ( "kind" = $4 OR ( "kind" = $5 AND price IS NULL) )
	ORDER BY created_at DESC
	LIMIT $6 OFFSET $7
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := o.Db.QueryContext(ctx, stmt,
		models.ORDER_OPEN_STATUS,
		models.ORDER_PARTIALLY_FILLED_STATUS,
		models.ORDER_BUY_TYPE,
		models.ORDER_LIMIT_KIND,
		models.ORDER_MARKET_KIND,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.AssetId,
			&order.Type,
			&order.Kind,
			&order.Status,
			&order.Price,
			&order.Quantity,
			&order.InitialQty,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o orderImpl) FindMatchingSellOrder(order models.Order) (models.Order, error) {
	stmt := `
	SELECT 
		id,
		user_id, 
		asset_id, 
		"type", 
		kind, 
		status, 
		price, 
		quantity, 
		inital_quantity,
		created_at, 
		updated_at
	FROM orders where "type" =$1
	and status in ($2,$3) and asset_id = $4
	and (($5::NUMERIC is null and kind=$6 ) or price <=$5::NUMERIC)
	order by price asc
	limit 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	matchOrder := models.Order{}

	// var price *float64
	// price = &order.Price
	err := o.Db.QueryRowContext(ctx, stmt,
		models.ORDER_SELL_TYPE,
		models.ORDER_OPEN_STATUS,
		models.ORDER_PARTIALLY_FILLED_STATUS,
		order.AssetId,
		order.Price,
		models.ORDER_MARKET_KIND,
	).Scan(
		&matchOrder.Id,
		&matchOrder.UserId,
		&matchOrder.AssetId,
		&matchOrder.Type,
		&matchOrder.Kind,
		&matchOrder.Status,
		&matchOrder.Price,
		&matchOrder.Quantity,
		&matchOrder.InitialQty,
		&matchOrder.CreatedAt,
		&matchOrder.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return matchOrder, repository.ErrRecordNotFound

		default:
			return matchOrder, err
		}
	}

	return matchOrder, nil
}
