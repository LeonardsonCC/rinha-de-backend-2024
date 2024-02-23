package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/worker"
)

func main() {
	mux := http.DefaultServeMux

	_, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("POST /clientes/{id}/transacoes", api.HandleNewTransaction)
	mux.HandleFunc("GET /clientes/{id}/extrato", api.HandleGetBalance)

	// add transactions listening to channel
	for range 10 {
		go worker.ListenTransactionToAdd()
	}

	fmt.Println("listening to: 0.0.0.0:8888")
	_ = http.ListenAndServe("0.0.0.0:8888", mux)

	db.CloseConnection()
}
