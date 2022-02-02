package store

import "context"

type Store interface {
	User() UserRepository
	Post() PostRepository
	Close(context.Context) error
}
