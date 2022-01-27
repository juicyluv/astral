package store

import (
	"context"

	"github.com/juicyluv/astral/internal/model"
)

type UserRepository interface {
	Create(context.Context, *model.User) (int, error)
	FindAll(context.Context) ([]model.User, error)
	FindById(context.Context, int) (*model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
	Update(context.Context, *model.UpdateUserDto) error
	Delete(context.Context, int) error
}

type PostRepository interface {
	Create(context.Context, *model.Post) (int, error)
	FindAll(context.Context) ([]model.Post, error)
	FindById(context.Context, int) (*model.Post, error)
	FindByEmail(context.Context, string) (*model.Post, error)
	Update(context.Context, *model.UpdatePostDto) error
	Delete(context.Context, int) error
}
