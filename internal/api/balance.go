package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository"
)

func HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	c := context.Background()

	clientID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, errs.ErrAccountNotFound.Error(), http.StatusNotFound)
		return
	}

	balance, err := repository.GetBalance(c, clientID)
	if err != nil {
		if errors.Is(err, errs.ErrAccountNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str, _ := json.Marshal(balance)
	fmt.Fprint(w, string(str))
}
