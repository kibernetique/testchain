package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"testchain/config"
	"testchain/internal/blockchain"

	"github.com/gorilla/mux"
)

type Balance struct {
	logger           *log.Logger
	blockchainClient *blockchain.Client
	config           *config.Config
}

func NewBalance(logger *log.Logger, blockchainClient *blockchain.Client, config *config.Config) *Balance {
	return &Balance{
		logger:           logger,
		blockchainClient: blockchainClient,
		config:           config,
	}
}

func (h *Balance) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr, ok := vars["addr"]
	if !ok {
		h.handleBadRequestError(w, errors.New("Can't gat an address"))
		return
	}

	// validating address
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(addr) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}

	// getting balance
	ctx, cancelCtx := context.WithTimeout(context.Background(), h.config.BcRequestTimeout)
	defer cancelCtx()
	responce, err := h.blockchainClient.GetBalance(ctx, addr)
	if err != nil {
		h.handleBadRequestError(w, err)
		return
	}

	// building responce
	responceString, err := json.Marshal(responce)
	if err != nil {
		h.handleBadRequestError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(responceString))
	return
}

func (h *Balance) handleBadRequestError(w http.ResponseWriter, err error) {
	h.logger.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
