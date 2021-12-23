package psql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/skandyla/s3-uploader/internal/models"
)

type Users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *Users {
	return &Users{db}
}

func (r *Users) Create(ctx context.Context, user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (r *Users) GetByCredentials(ctx context.Context, email, password string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}
