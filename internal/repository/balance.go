package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
	"github.com/jackc/pgx/v5"
)

func GetBalance(c context.Context, clientID int) (*contracts.Balance, error) {
	db, _ := db.GetConnection()

	conn, _ := db.Acquire(c)
	defer conn.Release()

	b := new(pgx.Batch)

	b.Queue("SELECT c.saldo, c.limite FROM clientes c WHERE id = $1", clientID)
	b.Queue("SELECT t.valor, t.tipo, t.descricao, t.realizada_em FROM transacoes t WHERE cliente_id=$1 ORDER BY realizada_em DESC LIMIT 10", clientID)

	results := conn.SendBatch(c, b)

	balance := &contracts.Balance{
		Balance: contracts.BalanceData{
			BalanceDate: time.Now(),
		},
		LastTransactions: make([]*contracts.Transaction, 0, 10),
	}

	rows := results.QueryRow()
	err := rows.Scan(&balance.Balance.Total, &balance.Balance.Limit)
	if err != nil {
		return nil, errs.ErrAccountNotFound
	}

	txRows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to get last transactions: %s", err)
	}

	for txRows.Next() {
		tx := &contracts.Transaction{}
		_ = txRows.Scan(&tx.Value, &tx.Type, &tx.Description, &tx.CreatedAt)
		balance.LastTransactions = append(balance.LastTransactions, tx)
	}
	txRows.Close()

	results.Close()

	return balance, nil
}
