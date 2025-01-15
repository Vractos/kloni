package repository

import (
	"context"
	"errors"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OrderPostgreSQL struct {
	db     *pgxpool.Pool
	logger metrics.Logger
}

func NewOrderPostgreSQL(db *pgxpool.Pool, logger metrics.Logger) *OrderPostgreSQL {
	return &OrderPostgreSQL{db: db, logger: logger}
}

// RegisterOrder implements order.Repository
func (r *OrderPostgreSQL) RegisterOrder(o *entity.Order) error {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
  INSERT INTO orders(id, store_id, marketplace_id, date_created, status)
  VALUES($1,$2,$3,$4,$5)
  `, o.ID, o.AccountID, o.MarketplaceID, o.DateCreated, o.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
		}
		return err
	}

	for _, i := range o.Items {
		_, err := tx.Exec(ctx, `
    INSERT INTO order_items(id, title, sku, quantity, order_id)
    VALUES($1,$2,$3,$4,$5)
    `, i.ID, i.Title, i.Sku, i.Quantity, o.ID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
			}
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("Error to commit order", err)
		return errors.New("error to commit order")
	}

	return nil
}

// GetOrder implements order.Repository
func (r *OrderPostgreSQL) GetOrder(orderMarketplaceId string) (*entity.Order, error) {
	var order entity.Order

	err := r.db.QueryRow(context.Background(), `
	SELECT
  id,
  store_id,
  marketplace_id,
  status
  FROM
  orders
  WHERE
  marketplace_id=$1
	`, orderMarketplaceId).Scan(&order.ID, &order.AccountID, &order.MarketplaceID, &order.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
			return nil, err
		} else if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}
