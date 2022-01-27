package model

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RegisteredAt string `json:"registered_at"`
	Password     string `json:"-"`
}

type UpdateUserDto struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
