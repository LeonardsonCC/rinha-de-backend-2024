package repository

import (
	"fmt"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository/db"
)

func AddTransaction(clientID uint, txType rune, v int, d string) (*contracts.TransactionSuccess, error) {
	db, _ := db.GetConnection()

	tx, _ := db.Begin()
	defer tx.Rollback()

	var limit, balance int
	clientBefore, err := tx.Query("SELECT limite, saldo FROM clientes WHERE id = $1", clientID)
	if err != nil {
		return nil, err
	}

	if !clientBefore.Next() {
		return nil, errs.ErrAccountNotFound
	}

	err = clientBefore.Scan(&limit, &balance)
	if err != nil {
		return nil, err
	}
	clientBefore.Close()

	var newBalance int
	if txType == 'd' {
		newBalance = balance - v
	} else {
		newBalance = balance + v
	}

	if limit+newBalance < 0 {
		return nil, errs.ErrInsufficientLimit
	}

	// TODO: make it goroutine
	_, err = tx.Exec("INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)", clientID, string(txType), v, d)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tx: %v", err)
	}

	_, err = tx.Exec("UPDATE clientes SET saldo=$1 WHERE id=$2", newBalance, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to update client: %v", err)
	}

	_ = tx.Commit()

	return &contracts.TransactionSuccess{
		Limit:   int(limit),
		Balance: int(newBalance),
	}, err
}
