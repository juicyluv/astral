package repository

import (
	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/internal/model"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(*model.User) (int, error) {
	return 0, nil
}

func (r *UserRepository) Find(int) (*model.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByEmail(string) (*model.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(*model.UpdateUserDto) error {
	return nil
}

func (r *UserRepository) Delete(int) error {
	return nil
}
