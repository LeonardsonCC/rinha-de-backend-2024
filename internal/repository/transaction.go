package repository

import (
	"context"
	"fmt"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
)

func AddTransaction(c context.Context, t contracts.Transaction) (*contracts.TransactionSuccess, error) {
	db, _ := db.GetConnection()

	tx, _ := db.Begin(c)
	defer tx.Rollback(c)

	value := t.Value
	if []rune(t.Type)[0] == 'd' {
		value = t.Value * (-1)
	}

	var newBalance, limit int
	balanceResult := tx.QueryRow(c, "UPDATE clientes SET saldo=(saldo+$1) WHERE id=$2 RETURNING saldo, limite", value, t.ClientID)

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

func AddSuccessTransaction(c context.Context, t contracts.Transaction) error {
	db, _ := db.GetConnection()

	_, err := db.Exec(c, "INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)", t.ClientID, t.Type, t.Value, t.Description)
	if err != nil {
		return fmt.Errorf("failed to insert and update client: %v", err)
	}

	return nil
}
