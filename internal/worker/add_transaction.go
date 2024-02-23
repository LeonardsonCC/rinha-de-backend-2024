package worker

import (
	"context"
	"fmt"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository"
)

var transactionsToAdd = make(chan contracts.Transaction)

func AddTransactionToAdd(t contracts.Transaction) {
	transactionsToAdd <- t
}

func ListenTransactionToAdd() {
	for t := range transactionsToAdd {
		err := repository.AddSuccessTransaction(context.Background(), t)
		if err != nil {
			fmt.Printf("failed to add transaction [%+v], err: %v\n", t, err)
		}
	}
}
