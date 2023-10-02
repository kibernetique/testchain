package blockchain

import (
	"context"
	"testchain/internal/models"
)

type ClientInterface interface {
	GetTxInfo(ctx context.Context, hash string) (*models.TxInfo, error)
	GetBalance(ctx context.Context, hash string) (*models.Balance, error)
}
