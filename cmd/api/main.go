package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
)

func main() {
	_, err := db.GetConnection()
	if err != nil {
		panic(err)
	}
	defer db.CloseConnection()

	// mux.HandleFunc("POST /clientes/{id}/transacoes", api.HandleNewTransaction)
	// mux.HandleFunc("GET /clientes/{id}/extrato", api.HandleGetBalance)

	// _ = http.ListenAndServe("0.0.0.0:8888", mux)

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/clientes/:id/extrato", api.HandleGetBalance)
	app.Post("/clientes/:id/transacoes", api.HandleNewTransaction)

	fmt.Println("listening to: 0.0.0.0:8888")
	err = app.Listen(":8888")
	if err != nil {
		panic(err)
	}
}
