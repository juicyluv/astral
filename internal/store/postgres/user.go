package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/astral/internal/model"
	"go.uber.org/zap"
)

type UserRepository struct {
	db     *pgx.Conn
	logger *zap.SugaredLogger
}

func NewUserRepository(db *pgx.Conn, logger *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int, error) {
	query := `
	INSERT INTO users(username, email, password) 
	VALUES($1, $2, $3)
	RETURNING user_id`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.Id)

	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User

	query := `
	SELECT user_id, username, email, 
	TO_CHAR(registered_at, 'DD-MM-YYYY') as registered_at
	FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.RegisteredAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindById(ctx context.Context, userId int) (*model.User, error) {
	var user model.User

	query := `
	SELECT user_id, username, email, 
	TO_CHAR(registered_at, 'DD-MM-YYYY') as registered_at
	FROM users
	WHERE user_id = $1`

	err := r.db.QueryRow(ctx, query, userId).Scan(
		&user.Id, &user.Username, &user.Email, &user.RegisteredAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := `
	SELECT user_id, username, email, 
	TO_CHAR(registered_at, 'DD-MM-YYYY') as registered_at
	FROM users
	WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id, &user.Username, &user.Email, &user.RegisteredAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.UpdateUserDto) error {
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userId int) error {
	query := `
	DELETE FROM users
	WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
