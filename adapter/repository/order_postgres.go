package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Vractos/dolly/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderPostgreSQL struct {
	db *pgxpool.Pool
}

func NewOrderPostgreSQL(db *pgxpool.Pool) *OrderPostgreSQL {
	return &OrderPostgreSQL{db: db}
}

// RegisterOrder implements order.Repository
func (r *OrderPostgreSQL) RegisterOrder(o *entity.Order) error {
	_, err := r.db.Exec(context.Background(), `
  INSERT INTO orders(id, marketplace_id, date_created, status)
  VALUES($1,$2,$3,$4)
  `, o.ID, o.MarketplaceID, o.DateCreated, o.Status)
	if err != nil {
		return err
	}
	return nil
}

// GetOrder implements order.Repository
func (r *OrderPostgreSQL) GetOrder(orderMarketplaceId string) (*entity.Order, error) {
	var order entity.Order

	err := r.db.QueryRow(context.Background(), `
	SELECT
    id,
    marketplace_id,
    status
  FROM
    orders
  WHERE
  marketplace_id=$1
	`, orderMarketplaceId).Scan(&order.ID, &order.MarketplaceID, &order.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(pgErr.Message)
			log.Println(pgErr.Code)
			return nil, err
		} else if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}
