package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TransactionProvider interface {
	LoadTransactions() []Transaction
}

type HttpServer struct {
	transactionProvider TransactionProvider
}

func SetupHttpServer(transactionProvider TransactionProvider) *HttpServer {
	httpServer := HttpServer{transactionProvider}
	http.HandleFunc("/health", httpServer.getHealth)
	http.HandleFunc("/transactions", httpServer.getTransactions)
	return &httpServer
}

func (httpServer *HttpServer) listenAndServe() {
	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		fmt.Println("Error serving http server")
	}
}

func (httpServer *HttpServer) getHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "UP\n")
}

func (httpServer *HttpServer) getTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := httpServer.transactionProvider.LoadTransactions()
	w.Header().Set("Content-Type", "application/json")
	responseJSON, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(responseJSON))
}
