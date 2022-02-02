package store

import (
	"context"

	"github.com/juicyluv/astral/internal/handler/filter"
	"github.com/juicyluv/astral/internal/model"
)

type UserRepository interface {
	Create(context.Context, *model.User) (int, error)
	FindAll(context.Context) ([]model.User, error)
	FindById(context.Context, int) (*model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
	Update(context.Context, int, *model.UpdateUserDto) error
	Delete(context.Context, int) error
	ConfirmEmail(context.Context, int) error
}

type PostRepository interface {
	Create(context.Context, *model.Post) (int, error)
	FindAll(context.Context, *filter.PostFilter) ([]model.Post, error)
	FindById(context.Context, int) (*model.Post, error)
	FindUserPosts(context.Context, int) ([]model.Post, error)
	Update(context.Context, int, *model.UpdatePostDto) error
	Delete(context.Context, int) error
}
