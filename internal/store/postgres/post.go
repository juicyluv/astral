package postgres

import (
	"context"
	"fmt"
	"strings"

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
	INSERT INTO posts(title, content, author_id) 
	VALUES($1, $2, $3)
	RETURNING post_id`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Author.Id,
	).Scan(&post.Id)

	if err != nil {
		if err = tx.Rollback(ctx); err != nil {
			return 0, err
		}
		return 0, err
	}

	query = `INSERT INTO user_post VALUES ($1, $2)`
	_, err = tx.Exec(ctx, query, post.Author.Id, post.Id)
	if err != nil {
		if err = tx.Rollback(ctx); err != nil {
			return 0, err
		}
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return post.Id, nil
}

func (r *PostRepository) FindAll(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post

	query := `
	SELECT 
	p.post_id, p.title, p.content, 
	TO_CHAR(p.created_at, 'DD-MM-YYYY') as created_at, 
	TO_CHAR(p.updated_at, 'DD-MM-YYYY') as updated_at, 
	u.user_id, u.username 
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
			&post.Content,
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

func (r *PostRepository) FindById(ctx context.Context, postId int) (*model.Post, error) {
	var post model.Post

	query := `
	SELECT 
	p.post_id, p.title, p.content, 
	TO_CHAR(p.created_at, 'DD-MM-YYYY') as created_at, 
	TO_CHAR(p.updated_at, 'DD-MM-YYYY') as updated_at, 
	u.user_id, u.username 
	FROM posts p
	INNER JOIN users u 
	ON u.user_id = p.author_id
	WHERE post_id = $1`

	err := r.db.QueryRow(ctx, query, postId).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Author.Id,
		&post.Author.Username,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) Update(ctx context.Context, postId int, post *model.UpdatePostDto) error {
	values := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if post.Title != nil {
		values = append(values, fmt.Sprintf("title=$%d", argId))
		args = append(args, *post.Title)
		argId++
	}

	if post.Content != nil {
		values = append(values, fmt.Sprintf("content=$%d", argId))
		args = append(args, *post.Content)
		argId++
	}

	if post.AuthorId != nil {
		values = append(values, fmt.Sprintf("author_id=$%d", argId))
		args = append(args, *post.AuthorId)
		argId++
	}

	valuesQuery := strings.Join(values, ", ")
	query := fmt.Sprintf("UPDATE posts SET %s WHERE post_id = $%d", valuesQuery, argId)
	args = append(args, postId)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *PostRepository) Delete(ctx context.Context, postId int) error {
	query := `
	DELETE FROM posts
	WHERE post_id = $1`

	_, err := r.db.Exec(ctx, query, postId)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) FindUserPosts(ctx context.Context, userId int) ([]model.Post, error) {
	var posts []model.Post

	query := `
	SELECT 
	p.post_id, p.title, p.content, 
	TO_CHAR(p.created_at, 'DD-MM-YYYY') as created_at, 
	TO_CHAR(p.updated_at, 'DD-MM-YYYY') as updated_at, 
	u.user_id, u.username 
	FROM posts p
	INNER JOIN users u 
	ON u.user_id = p.author_id
	WHERE p.author_id = $1`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
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
