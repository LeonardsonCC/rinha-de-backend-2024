package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/api/contracts"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository"
	"github.com/gofiber/fiber/v3"
)

var (
	ErrValueMustBeGreaterThan0 = errors.New("invalid value")
	ErrDescriptionInvalid      = errors.New("invalid description")
	ErrTypeInvalid             = errors.New("invalid type")
)

func HandleNewTransaction(c fiber.Ctx) error {
	clientID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString(errs.ErrAccountNotFound.Error())
	}

	var tt contracts.Transaction
	if err := c.Bind().Body(&tt); err != nil {
		return c.Status(http.StatusBadGateway).SendString(err.Error())
	}

	if err := validateTx(tt); err != nil {
		return c.Status(http.StatusUnprocessableEntity).SendString(err.Error())
	}

	tx, err := repository.AddTransaction(c.Context(), clientID, []rune(tt.Type)[0], tt.Value, tt.Description)
	if err != nil {
		if errors.Is(err, errs.ErrInsufficientLimit) {
			return c.Status(http.StatusUnprocessableEntity).SendString(err.Error())
		}
		if errors.Is(err, errs.ErrAccountNotFound) {
			return c.Status(http.StatusNotFound).SendString(err.Error())
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	str, _ := json.Marshal(tx)
	return c.Status(200).Send(str)
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
