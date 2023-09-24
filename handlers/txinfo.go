// 0xd9510912ac124f7dd525ae714b4fbab2652051beff9de279315419ec3cc4ffcf
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

type TxInfo struct {
	logger *log.Logger
}

type TxInfoResponce struct {
	Hash        string `json:"hash"`
	IsPending   bool   `json:"is_pending"`
	ChainId     string `json:"chain_id"`
	Cost        string `json:"cost"`
	To          string `json:"to"`
	Sender      string `json:"sender"`
	Value       string `json:"value"`
	Data        []byte `json:"data"`
	BlockNumber string `json:"block_nuber"`
	GasUsed     uint64 `json:"gas_used"`
}

func NewTxInfo(logger *log.Logger) *TxInfo {
	return &TxInfo{
		logger: logger,
	}
}

func (b *TxInfo) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txHex, ok := vars["txhash"]
	if !ok {
		b.handleBadRequestError(w, errors.New("Can't gat a transaction hash"))
		return
	}

	ctx := context.Background()

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/4bf9cf69d33f4c7395938604e3e69ce0")
	if err != nil {
		b.handleBadRequestError(w, err)
		return
	}

	txHash := common.HexToHash(txHex)
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		b.handleBadRequestError(w, errors.New("TransactionByHash error: "+err.Error()))
		return
	}

	responce := TxInfoResponce{
		Hash:      txHex,
		IsPending: isPending,
		ChainId:   tx.ChainId().String(),
		Cost:      tx.Cost().String(),
		To:        tx.To().String(),
		Value:     tx.Value().String(),
		Data:      tx.Data(),
	}

	// get sender address
	from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		b.handleBadRequestError(w, errors.New("Sender getting error: "+err.Error()))
		return
	}
	responce.Sender = from.Hex()

	// get block number
	receipt, err := client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		b.handleBadRequestError(w, errors.New("TransactionReceipt error: "+err.Error()))
		return
	}
	responce.BlockNumber = receipt.BlockNumber.String()
	responce.GasUsed = receipt.GasUsed

	// responsing
	responceString, err := json.Marshal(responce)
	if err != nil {
		b.handleBadRequestError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(responceString))
	return

}

func (b *TxInfo) handleBadRequestError(w http.ResponseWriter, err error) {
	b.logger.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
