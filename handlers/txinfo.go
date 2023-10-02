package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"testchain/config"
	"testchain/internal/blockchain"

	"github.com/gorilla/mux"
)

type TxInfo struct {
	logger           *log.Logger
	blockchainClient *blockchain.Client
	config           *config.Config
}

func NewTxInfo(logger *log.Logger, blockchainClient *blockchain.Client, config *config.Config) *TxInfo {
	return &TxInfo{
		logger:           logger,
		blockchainClient: blockchainClient,
		config:           config,
	}
}

func (h *TxInfo) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash, ok := vars["txhash"]
	if !ok {
		h.handleBadRequestError(w, errors.New("Can't get a transaction hash"))
		return
	}

	// getting transaction info
	ctx, cancelCtx := context.WithTimeout(context.Background(), h.config.BcRequestTimeout)
	defer cancelCtx()
	responce, err := h.blockchainClient.GetTxInfo(ctx, hash)
	if err != nil {
		h.handleBadRequestError(w, err)
		return
	}

	// responsing
	responceString, err := json.Marshal(responce)
	if err != nil {
		h.handleBadRequestError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(responceString))
	return

}

func (h *TxInfo) handleBadRequestError(w http.ResponseWriter, err error) {
	h.logger.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
