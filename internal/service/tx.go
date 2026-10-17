package service

import (
	"context"
	"shop_api/pkg/postgres"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(postgres.Querier) error) error
}
