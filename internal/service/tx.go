package service

import (
	"context"
	"shop_api/pkg/postgres"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(postgres.Querier) error) error
}

type txManager struct {
	db *postgres.Postgres
}

func (m *txManager) WithTx(ctx context.Context, fn func(postgres.Querier) error) error {
	return m.db.WithTx(ctx, fn)
}

func NewTxManager(db *postgres.Postgres) TxManager {
	return &txManager{db: db}
}
