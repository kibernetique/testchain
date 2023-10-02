package ethclient

import (
	"context"
	"encoding/hex"
	"math"
	"math/big"
	"testchain/internal/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	client *ethclient.Client
}

func NewEthClient(connUrl string) (*EthClient, error) {
	client, err := ethclient.Dial(connUrl)
	if err != nil {
		return nil, err
	}
	eth := EthClient{
		client: client,
	}
	return &eth, nil
}

func (ec *EthClient) GetTxInfo(ctx context.Context, hash string) (*models.TxInfo, error) {
	txHash := common.HexToHash(hash)
	tx, isPending, err := ec.client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	responce := models.TxInfo{
		Hash:      hash,
		IsPending: isPending,
		ChainId:   tx.ChainId().String(),
		Cost:      tx.Cost().String(),
		To:        tx.To().String(),
		Value:     tx.Value().String(),
		Data:      hex.EncodeToString(tx.Data()),
	}

	// get sender address
	from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		return nil, err
	}
	responce.Sender = from.Hex()

	// get block number
	receipt, err := ec.client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, err
	}
	responce.BlockNumber = receipt.BlockNumber.String()
	responce.GasUsed = receipt.GasUsed
	return &responce, nil
}

func (ec *EthClient) GetBalance(ctx context.Context, hash string) (*models.Balance, error) {
	account := common.HexToAddress(hash)
	balance, err := ec.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}

	ethBal := new(big.Float)
	ethBal.SetString(balance.String())
	ethValue := new(big.Float).Quo(ethBal, big.NewFloat(math.Pow10(18)))

	responce := models.Balance{
		Addr: hash,
		Wei:  balance.String(),
		Eth:  ethValue.Text('g', 1024),
	}
	return &responce, nil
}
