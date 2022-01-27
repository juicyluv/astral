package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Post struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    User   `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdatePostDto struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	AuthorId *int    `json:"author_id"`
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Title, is.Alphanumeric, validation.Length(1, 200), validation.Required),
		validation.Field(&p.Content, is.ASCII, validation.Length(1, 0), validation.Required),
	)
}

func (p *UpdatePostDto) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Title, is.Alphanumeric, validation.Length(1, 200)),
		validation.Field(&p.Content, is.ASCII, validation.Length(1, 0)),
	)
}
