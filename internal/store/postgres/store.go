package postgres

import (
	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/internal/store"
	"go.uber.org/zap"
)

type Store struct {
	user store.UserRepository
	post store.PostRepository
	db   *pgx.Conn
}

func NewPostgres(conn *pgx.Conn, logger *zap.SugaredLogger) *Store {
	return &Store{
		db:   conn,
		user: NewUserRepository(conn, logger),
		post: NewPostRepository(conn, logger),
	}
}

func (s *Store) User() store.UserRepository {
	return s.user
}

func (s *Store) Post() store.PostRepository {
	return s.post
}
