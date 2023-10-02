package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"testchain/config"
	"testchain/handlers"
	"testchain/internal/blockchain"
	"testchain/internal/blockchain/ethclient"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "ethereum-test", log.LstdFlags)

	ctx := context.Background()

	config, err := config.GetFromEnv(ctx)
	if err != nil {
		logger.Fatal("Config loading fail: " + err.Error())
	}

	// creating blockchain client with ethereum client
	ethclient, err := ethclient.NewEthClient(config.BcConnectionUrl)
	blockchainClient := blockchain.NewClient(ethclient)

	// creating hanldlers
	balanceHandler := handlers.NewBalance(logger, blockchainClient, config)
	txInfoHandler := handlers.NewTxInfo(logger, blockchainClient, config)

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/balance/{addr}", balanceHandler.Get).Methods(http.MethodGet)
	r.HandleFunc("/txinfo/{txhash}", txInfoHandler.Get).Methods(http.MethodGet)
	http.Handle("/", r)

	s := &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      r,
		IdleTimeout:  config.ServerIdleTimeout * time.Second,
		ReadTimeout:  config.ServerReadTimeout * time.Minute,
		WriteTimeout: config.ServerWriteTimeout * time.Second,
	}

	// start the server
	go func() {
		logger.Printf("Starting server at addr %s...\n", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	logger.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting for current operations to complete
	shutdownCtx, _ := context.WithTimeout(context.Background(), config.ServerShutdownTimeout*time.Second)
	s.Shutdown(shutdownCtx)
}
