package postgres

import (
	"context"
	"fmt"
	"log"
	"shop_api/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var ErrNoRows = pgx.ErrNoRows

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(lc fx.Lifecycle, cfg *config.DatabaseConfig) (*Postgres, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database does not respond: %w", err)
	}

	db := &Postgres{
		Pool: dbpool,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection")

			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			db.close()

			log.Println("Database connection closed")
			return nil
		},
	})

	return db, nil
}

func (p *Postgres) close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
