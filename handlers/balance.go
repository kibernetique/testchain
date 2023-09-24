// 0x5eD8Cee6b63b1c6AFce3AD7c92f4fD7E1B8fAd9F
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/big"
	"net/http"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

type Balance struct {
	logger *log.Logger
}

type BalanceResponce struct {
	Addr string `json:"addr"`
	Wei  string `json:"wei"`
	Eth  string `json:"eth"`
}

func NewBalance(logger *log.Logger) *Balance {
	return &Balance{
		logger: logger,
	}
}

func (b *Balance) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr, ok := vars["addr"]
	if !ok {
		b.handleBadRequestError(w, errors.New("Can't gat an address"))
		return
	}

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(addr) {
		b.handleBadRequestError(w, errors.New("Invalid address"))
		return
	}

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/4bf9cf69d33f4c7395938604e3e69ce0")
	if err != nil {
		b.handleBadRequestError(w, err)
		return
	}

	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		b.handleBadRequestError(w, err)
		return
	}

	ethBal := new(big.Float)
	ethBal.SetString(balance.String())
	ethValue := new(big.Float).Quo(ethBal, big.NewFloat(math.Pow10(18)))

	responce := BalanceResponce{
		Addr: addr,
		Wei:  balance.String(),
		Eth:  ethValue.Text('g', 1024),
	}
	responceString, err := json.Marshal(responce)
	if err != nil {
		b.handleBadRequestError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(responceString))
	return
}

func (b *Balance) handleBadRequestError(w http.ResponseWriter, err error) {
	b.logger.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
