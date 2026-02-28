package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type SettingsRepo struct {
	db *sql.DB
}

func NewSettingsRepo(db *sql.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

func (r *SettingsRepo) Get(ctx context.Context, userID int64, key string) (string, error) {
	var value string
	err := r.db.QueryRowContext(ctx,
		"SELECT value FROM settings WHERE user_id = ? AND key = ?", userID, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (r *SettingsRepo) Set(ctx context.Context, userID int64, key, value string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO settings (user_id, key, value, updated_at)
		 VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(user_id, key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`,
		userID, key, value)
	return err
}

func (r *SettingsRepo) GetMulti(ctx context.Context, userID int64, keys []string) (map[string]string, error) {
	result := make(map[string]string)
	if len(keys) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(keys))
	args := make([]interface{}, 0, len(keys)+1)
	args = append(args, userID)
	for i, k := range keys {
		placeholders[i] = "?"
		args = append(args, k)
	}

	query := fmt.Sprintf(
		"SELECT key, value FROM settings WHERE user_id = ? AND key IN (%s)",
		strings.Join(placeholders, ","),
	)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, nil
}

func (r *SettingsRepo) SetMulti(ctx context.Context, userID int64, values map[string]string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO settings (user_id, key, value, updated_at)
		 VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(user_id, key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for k, v := range values {
		if _, err := stmt.ExecContext(ctx, userID, k, v); err != nil {
			return err
		}
	}
	return tx.Commit()
}
