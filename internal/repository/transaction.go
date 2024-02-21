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

	var limit, balance int
	clientBefore := tx.QueryRow(c, "SELECT limite, saldo FROM clientes WHERE id = $1 FOR NO KEY UPDATE", clientID)

	err := clientBefore.Scan(&limit, &balance)
	if err != nil {
		return nil, err
	}

	var newBalance int
	if txType == 'd' {
		newBalance = balance - v
	} else {
		newBalance = balance + v
	}

	if newBalance < (limit * -1) {
		return nil, errs.ErrInsufficientLimit
	}

	b := new(pgx.Batch)

	b.Queue("INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)", clientID, string(txType), v, d)
	b.Queue("UPDATE clientes SET saldo=$1 WHERE id=$2", newBalance, clientID)

	bResult := tx.SendBatch(c, b)
	err = bResult.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to insert and update client: %v", err)
	}

	_ = tx.Commit(c)

	return &contracts.TransactionSuccess{
		Limit:   int(limit),
		Balance: int(newBalance),
	}, err
}
