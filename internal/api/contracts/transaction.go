package contracts

import "time"

type Transaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type TransactionSuccess struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}
