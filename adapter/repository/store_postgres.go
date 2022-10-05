package repository

import (
	"context"

	"github.com/Vractos/dolly/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type StorePostgreSQL struct {
	db *pgx.Conn
}

func NewStorePostgreSQL(db *pgx.Conn) *StorePostgreSQL {
	return &StorePostgreSQL{db: db}
}

// Get implements store.Repository
func (r *StorePostgreSQL) Get(id string) (*entity.Store, error) {
	panic("unimplemented")
}

// Create implements store.Repository
func (r *StorePostgreSQL) Create(e *entity.Store) (uuid.UUID, error) {
	_, err := r.db.Exec(context.Background(), `
    INSERT INTO store (id, name, email)
    VALUES($1,$2,$3)
  `, e.ID, e.Name, e.Email)
	if err != nil {
		return e.ID, err
	}

	return e.ID, nil
}

// Delete implements store.Repository
func (r *StorePostgreSQL) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

// Update implements store.Repository
func (r *StorePostgreSQL) Update(e *entity.Store) error {
	_, err := r.db.Exec(context.Background(), `
  INSERT INTO mercadolivre_credentials(owner_id, access_token, expires_in, user_id, refresh_token)
  VALUES($1,$2,$3,$4,$5)
  `)
	if err != nil {
		return err
	}
	return nil
}
