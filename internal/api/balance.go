package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/errs"
	"github.com/LeonardsonCC/rinha-de-backend-2024/internal/repository"
	"github.com/gofiber/fiber/v3"
)

func HandleGetBalance(c fiber.Ctx) error {
	clientID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString(errs.ErrAccountNotFound.Error())
	}

	balance, err := repository.GetBalance(c.Context(), clientID)
	if err != nil {
		if errors.Is(err, errs.ErrAccountNotFound) {
			return c.Status(http.StatusNotFound).SendString(err.Error())

		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	str, _ := json.Marshal(balance)
	return c.Status(200).Send(str)
}
