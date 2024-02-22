package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
)

func main() {
	mux := http.NewServeMux()

	_, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("GET /debug/pprof/cpu", pprof.Profile)
	mux.HandleFunc("POST /clientes/{id}/transacoes", api.HandleNewTransaction)
	mux.HandleFunc("GET /clientes/{id}/extrato", api.HandleGetBalance)

	fmt.Println("listening to: 0.0.0.0:8888")
	_ = http.ListenAndServe("0.0.0.0:8888", mux)

	db.CloseConnection()
}
