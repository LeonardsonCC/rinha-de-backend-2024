package contracts

import "time"

type Balance struct {
	Balance          BalanceData    `json:"saldo"`
	LastTransactions []*Transaction `json:"ultimas_transacoes"`
}

type BalanceData struct {
	Total       int       `json:"total"`
	BalanceDate time.Time `json:"data_extrato"`
	Limit       int       `json:"limite"`
}
