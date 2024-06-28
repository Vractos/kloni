package repository

import (
	"context"
	"errors"

	"github.com/Vractos/kloni/entity"
	"github.com/Vractos/kloni/pkg/metrics"
	"github.com/Vractos/kloni/usecases/common"
	"github.com/Vractos/kloni/usecases/store"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type StorePostgreSQL struct {
	db     *pgxpool.Pool
	logger metrics.Logger
}

func NewStorePostgreSQL(db *pgxpool.Pool, logger metrics.Logger) *StorePostgreSQL {
	return &StorePostgreSQL{db: db, logger: logger}
}

// Get implements store.Repository
func (r *StorePostgreSQL) Get(id string) (*entity.Store, error) {
	panic("unimplemented")
}

// Create implements store.Repository
func (r *StorePostgreSQL) Create(e *entity.Store) (entity.ID, error) {
	_, err := r.db.Exec(context.Background(), `
    INSERT INTO store (id, name, email)
    VALUES($1,$2,$3)
    `, e.ID, e.Name, e.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
		}
		return e.ID, err
	}

	return e.ID, nil
}

// Delete implements store.Repository
func (r *StorePostgreSQL) Delete(id entity.ID) error {
	panic("unimplemented")
}

// Update implements store.Repository
func (r *StorePostgreSQL) Update(e *entity.Store) error {
	panic("unimplemented")
}

// RegisterMeliCredential implements store.Repository
func (r *StorePostgreSQL) RegisterMeliCredential(id entity.ID, c *common.MeliCredential) error {
	_, err := r.db.Exec(context.Background(), `
  INSERT INTO mercadolivre_credentials(owner_id, access_token, expires_in, user_id, refresh_token, updated_at)
  VALUES($1,$2,$3,$4,$5,$6)
  `, id, c.AccessToken, c.ExpiresIn, c.UserID, c.RefreshToken, c.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
		}
		return err
	}
	return nil
}

// RetrieveMeliCredentials implements store.Repository
func (r *StorePostgreSQL) RetrieveMeliCredentialsFromStoreID(id entity.ID) (*[]store.Credentials, error) {
	var credentials []store.Credentials

	rows, err := r.db.Query(context.Background(), `
		SELECT
			mc.id as account_id,
			mc.account_name,
		  mc.user_id AS mercadolivre_user_id,
			mc.access_token,
			mc.refresh_token,
			mc.updated_at
		FROM
			mercadolivre_credentials mc
		WHERE
			mc.owner_id = $1
    `, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
		}
		return nil, err
	}

	for rows.Next() {
		credential := store.Credentials{MeliCredential: &common.MeliCredential{}}

		err = rows.Scan(
			&credential.ID,
			&credential.AccountName,
			&credential.UserID,
			&credential.AccessToken,
			&credential.RefreshToken,
			&credential.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		credentials = append(credentials, credential)
	}

	return &credentials, nil
}

// RetrieveMeliCredentials implements store.Repository
func (r *StorePostgreSQL) RetrieveMeliCredentialsFromMeliUserID(accountId string) (*store.Credentials, error) {
	var credential store.Credentials

	err := r.db.QueryRow(context.Background(), `
	SELECT
		mc.id as account_id,
		mc.owner_id,
		mc.account_name,
		mc.access_token,
		mc.refresh_token,
		mc.updated_at
  	FROM
   		mercadolivre_credentials mc
     WHERE
     	id=$1
	`, accountId).Scan(
		&credential.ID,
		&credential.OwnerID,
		&credential.AccountName,
		&credential.AccessToken,
		&credential.RefreshToken,
		&credential.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
		}
		return nil, err
	}

	return &credential, nil
}

// UpdateMeliCredentials implements store.Repository
func (r *StorePostgreSQL) UpdateMeliCredentials(accountId entity.ID, c *common.MeliCredential) error {
	_, err := r.db.Exec(context.Background(), `
    UPDATE
    	mercadolivre_credentials
    SET
    	access_token=$1,
     	refresh_token=$2,
      	updated_at=$3
    WHERE
    	id=$4
    `, c.AccessToken, c.RefreshToken, c.UpdatedAt, accountId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(pgErr.Message, pgErr, zap.String("db_error_code", pgErr.Code))
		}
		return err
	}
	return nil
}
