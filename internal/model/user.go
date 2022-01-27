package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email,omitempty"`
	RegisteredAt string `json:"registered_at,omitempty"`
	Password     string `json:"password,omitempty"`
}

type UpdateUserDto struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, is.Alphanumeric, validation.Length(3, 20), validation.Required),
		validation.Field(&u.Email, is.Email, validation.Required),
		validation.Field(&u.Password, is.Alphanumeric, validation.Required),
	)
}

func (u *UpdateUserDto) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, is.Alphanumeric, validation.Length(3, 20)),
		validation.Field(&u.Email, is.Email),
		validation.Field(&u.Password, is.Alphanumeric),
	)
}
