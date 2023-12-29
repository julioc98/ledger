package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	CreateDepositHandler       http.HandlerFunc
	CreateTransferHandler      http.HandlerFunc
	GetBalanceHandler          http.HandlerFunc
	GetTransfersHistoryHandler http.HandlerFunc
}

func (a *API) Routes(router *chi.Mux) {
	router.Post("/api/v1/account/{account}/deposit", a.CreateDepositHandler)
	router.Post("/api/v1/account/{account}/transfers", a.CreateTransferHandler)
	router.Get("/api/v1/account/{account}/balance", a.GetBalanceHandler)
	router.Get("/api/v1/account/{account}/transfers", a.GetTransfersHistoryHandler)
}
