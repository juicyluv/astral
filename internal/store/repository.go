package store

import "github.com/juicyluv/astral/internal/model"

type UserRepository interface {
	Create(*model.User) (int, error)
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Update(*model.UpdateUserDto) error
	Delete(int) error
}
