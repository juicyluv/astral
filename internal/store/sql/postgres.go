package sql

import (
	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/internal/store"
)

type Store struct {
	user store.UserRepository
	db   *pgx.Conn
}

func NewPostgres(conn *pgx.Conn) *Store {
	return &Store{
		db: conn,
	}
}

func (s *Store) User() store.UserRepository {
	return s.user
}
