package repository

import (
	"context"

	"github.com/Vractos/dolly/entity"
	"github.com/Vractos/dolly/usecases/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorePostgreSQL struct {
	db *pgxpool.Pool
}

func NewStorePostgreSQL(db *pgxpool.Pool) *StorePostgreSQL {
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
	panic("unimplemented")
}

// RegisterMeliCredential implements store.Repository
func (r *StorePostgreSQL) RegisterMeliCredential(id uuid.UUID, c *store.MeliCredential) error {
	_, err := r.db.Exec(context.Background(), `
  INSERT INTO mercadolivre_credentials(owner_id, access_token, expires_in, user_id, refresh_token, updated_at)
  VALUES($1,$2,$3,$4,$5,$6)
  `, id, c.AccessToken, c.ExpiresIn, c.UserID, c.RefreshToken, c.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
