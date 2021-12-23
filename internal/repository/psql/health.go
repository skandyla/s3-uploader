package psql

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/skandyla/s3-uploader/internal/models"
)

type HealthRepository struct {
	db *sqlx.DB
}

func NewHealthRepository(db *sqlx.DB) *HealthRepository {
	return &HealthRepository{db: db}
}

// Health checks availability of storage
func (r *HealthRepository) Ping(ctx context.Context) error {
	return r.db.Ping()
}

func (r *HealthRepository) Info(ctx context.Context) (models.InfoDependencyItem, error) {

	start := time.Now()

	info := models.InfoDependencyItem{}
	info.Status = 200

	var result int
	err := r.db.QueryRow("select(1 + 1)").
		Scan(&result)

	if err != nil {
		return info, err
	}

	if result != 2 {
		err = errors.New("unknown error")
	}

	diff := time.Since(start).Microseconds()
	info.Latency = float64(diff) / 1000000.0

	return info, err
}
