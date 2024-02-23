package repository

import (
	"context"
	"fmt"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
)

func AddTransaction(c context.Context, clientID int, txType rune, v int, d string) (*contracts.TransactionSuccess, error) {
	db, _ := db.GetConnection()

	tx, _ := db.Begin(c)
	defer tx.Rollback(c)

	value := v
	if txType == 'd' {
		value = v * (-1)
	}

	_, err := tx.Exec(c, "INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)", clientID, string(txType), v, d)
	if err != nil {
		return nil, fmt.Errorf("failed to insert and update client: %v", err)
	}

	var newBalance, limit int
	balanceResult := tx.QueryRow(c, "UPDATE clientes SET saldo=(saldo+$1) WHERE id=$2 RETURNING saldo, limite", value, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get returning: %v", err)
	}

	if err := balanceResult.Scan(&newBalance, &limit); err != nil {
		return nil, fmt.Errorf("failed to updated: %v", err)
	}

	if newBalance < limit*-1 {
		return nil, errs.ErrInsufficientLimit
	}

	_ = tx.Commit(c)

	return &contracts.TransactionSuccess{
		Limit:   int(limit),
		Balance: int(newBalance),
	}, nil
}
