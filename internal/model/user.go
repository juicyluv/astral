package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email,omitempty"`
	RegisteredAt string `json:"registered_at,omitempty"`
	Password     string `json:"password,omitempty"`
	IsVerified   bool   `json:"verified"`
}

type UpdateUserDto struct {
	Username   *string `json:"username"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
	IsVerified *bool   `json:"verified"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, is.Alphanumeric, validation.Length(3, 20), validation.Required),
		validation.Field(&u.Email, is.Email, validation.Required),
		validation.Field(&u.Password, is.Alphanumeric, validation.Required),
	)
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return err
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) ClearPassword() {
	u.Password = ""
}

func (u *UpdateUserDto) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Username, is.Alphanumeric, validation.Length(3, 20)),
		validation.Field(&u.Email, is.Email),
		validation.Field(&u.Password, is.Alphanumeric),
	)
}
