package psql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/skandyla/s3-uploader/internal/models"
)

type Tokens struct {
	db *sqlx.DB
}

func NewTokens(db *sqlx.DB) *Tokens {
	return &Tokens{db}
}

func (r *Tokens) Create(ctx context.Context, token models.RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (r *Tokens) Get(ctx context.Context, token string) (models.RefreshSession, error) {
	var t models.RefreshSession
	err := r.db.QueryRow("SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).
		Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = r.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", t.UserID)

	return t, err
}
