package blockchain

import (
	"context"
	"testchain/internal/models"
)

type Client struct {
	client ClientInterface
}

func NewClient(client ClientInterface) *Client {
	return &Client{
		client: client,
	}
}

func (bc *Client) GetTxInfo(ctx context.Context, hash string) (*models.TxInfo, error) {
	return bc.client.GetTxInfo(ctx, hash)
}

func (bc *Client) GetBalance(ctx context.Context, hash string) (*models.Balance, error) {
	return bc.client.GetBalance(ctx, hash)
}
