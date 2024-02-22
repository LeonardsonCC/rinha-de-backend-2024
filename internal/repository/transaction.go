package repository

import (
	"context"
	"fmt"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
	"github.com/jackc/pgx/v5"
)

func AddTransaction(c context.Context, clientID int, txType rune, v int, d string) (*contracts.TransactionSuccess, error) {
	db, _ := db.GetConnection()

	tx, _ := db.Begin(c)
	defer tx.Rollback(c)

	value := v
	if txType == 'd' {
		value = v * (-1)
	}

	b := new(pgx.Batch)

	b.Queue("INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)", clientID, string(txType), v, d)
	b.Queue("UPDATE clientes SET saldo=(saldo+$1) WHERE id=$2 RETURNING saldo, limite", value, clientID)

	bResult := tx.SendBatch(c, b)
	_, err := bResult.Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to insert and update client: %v", err)
	}

	var newBalance, limit int
	balanceResult, err := bResult.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to get returning: %v", err)
	}

	for balanceResult.Next() {
		if err := balanceResult.Scan(&newBalance, &limit); err != nil {
			return nil, fmt.Errorf("failed to updated: %v", err)
		}
	}

	balanceResult.Close()

	if newBalance < limit*-1 {
		return nil, errs.ErrInsufficientLimit
	}

	bResult.Close()
	_ = tx.Commit(c)

	return &contracts.TransactionSuccess{
		Limit:   int(limit),
		Balance: int(newBalance),
	}, nil
}
