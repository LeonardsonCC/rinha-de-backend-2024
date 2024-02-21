package repository

import (
	"fmt"
	"time"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
)

func GetBalance(clientID int64) (*contracts.Balance, error) {
	db, _ := db.GetConnection()

	rows, err := db.Query("SELECT c.saldo, c.limite FROM clientes c WHERE id = $1", clientID)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	balance := &contracts.Balance{
		Balance: contracts.BalanceData{
			BalanceDate: time.Now(),
		},
		LastTransactions: make([]*contracts.Transaction, 0, 10),
	}

	rows.Next()
	err = rows.Scan(&balance.Balance.Total, &balance.Balance.Limit)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	txRows, err := db.Query("SELECT t.valor, t.tipo, t.descricao, t.realizada_em FROM transacoes t WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last transactions: %s", err)
	}

	for txRows.Next() {
		tx := &contracts.Transaction{}
		_ = txRows.Scan(&tx.Value, &tx.Type, &tx.Description, &tx.CreatedAt)
		balance.LastTransactions = append(balance.LastTransactions, tx)
	}

	return balance, nil
}