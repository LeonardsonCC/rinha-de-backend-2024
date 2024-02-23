package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/worker"
)

var (
	ErrValueMustBeGreaterThan0 = errors.New("invalid value")
	ErrDescriptionInvalid      = errors.New("invalid description")
	ErrTypeInvalid             = errors.New("invalid type")
)

func HandleNewTransaction(w http.ResponseWriter, r *http.Request) {
	c := context.Background()

	clientID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, errs.ErrAccountNotFound.Error(), http.StatusNotFound)
		return
	}

	var tt contracts.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tt.ClientID = clientID
	// to keep sorted when request balance
	tt.CreatedAt = time.Now()

	if err := validateTx(tt); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tx, err := repository.AddTransaction(c, tt)
	if err != nil {
		if errors.Is(err, errs.ErrInsufficientLimit) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, errs.ErrAccountNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	worker.AddTransactionToAdd(tt)

	str, _ := json.Marshal(tx)
	fmt.Fprint(w, string(str))
}

func validateTx(tt contracts.Transaction) error {
	if l := len(tt.Description); l < 1 || l > 10 {
		return ErrDescriptionInvalid
	}
	if tt.Type != "c" && tt.Type != "d" {
		return ErrTypeInvalid
	}

	return nil
}
