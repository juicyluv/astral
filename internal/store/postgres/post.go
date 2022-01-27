package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/internal/model"
	"go.uber.org/zap"
)

type PostRepository struct {
	db     *pgx.Conn
	logger *zap.SugaredLogger
}

func NewPostRepository(db *pgx.Conn, logger *zap.SugaredLogger) *PostRepository {
	return &PostRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostRepository) Create(ctx context.Context, post *model.Post) (int, error) {
	query := `
	INSERT INTO posts(title, subtitle, author_id) 
	VALUES($1, $2, $3)
	RETURNING post_id`

	err := r.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Subtitle,
		post.Author.Id,
	).Scan(&post.Id)

	if err != nil {
		return 0, err
	}

	return post.Id, nil
}

func (r *PostRepository) FindAll(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post

	query := `
	SELECT 
	p.post_id, p.title, p.subtitle, 
	TO_CHAR(p.created_at, 'DD-MM-YYYY') as created_at, 
	TO_CHAR(p.updated_at, 'DD-MM-YYYY') as updated_at, 
	u.user_id, u.username, 
	FROM posts p
	INNER JOIN users u 
	ON u.user_id = p.author_id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Subtitle,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author.Id,
			&post.Author.Username,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) FindById(ctx context.Context, userId int) (*model.Post, error) {
	var post model.Post

	query := `
	SELECT user_id, username, email, 
	TO_CHAR(registered_at, 'DD-MM-YYYY') as registered_at
	FROM users
	WHERE user_id = $1`

	err := r.db.QueryRow(ctx, query, userId).Scan()

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) FindByEmail(ctx context.Context, email string) (*model.Post, error) {
	var post model.Post

	query := `
	SELECT user_id, username, email, 
	TO_CHAR(registered_at, 'DD-MM-YYYY') as registered_at
	FROM users
	WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan()

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) Update(ctx context.Context, user *model.UpdatePostDto) error {
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, postId int) error {
	query := `
	DELETE FROM users
	WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, postId)
	if err != nil {
		return err
	}

	return nil
}
