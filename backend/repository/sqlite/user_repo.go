package sqlite

import (
	"context"
	"database/sql"

	"github.com/homemenu/backend/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (username, password_hash, nickname) VALUES (?, ?, ?)",
		user.Username, user.PasswordHash, user.Nickname,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, password_hash, nickname, created_at FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Nickname, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, password_hash, nickname, created_at FROM users WHERE username = ?", username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Nickname, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
