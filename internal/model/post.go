package model

type Post struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Author    User   `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdatePostDto struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	AuthorId int    `json:"author_id"`
}
