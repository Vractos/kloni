package repository

import (
	"context"

	"github.com/Vractos/dolly/entity"
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
