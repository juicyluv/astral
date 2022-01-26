package model

import "time"

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registered_at"`
	Password     string
}

type UpdateUserDto struct {
	Username     *string    `json:"username"`
	Email        *string    `json:"email"`
	RegisteredAt *time.Time `json:"registered_at"`
	Password     *string    `json:"password"`
}
