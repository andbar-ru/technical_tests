package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"booking_app/internal/infrastructure"
	"booking_app/internal/interfaces"
	"booking_app/internal/usecases"
)

func main() {
	logger := infrastructure.NewLogger()
	storeHandler := infrastructure.NewStoreHandler()
	repository := interfaces.NewStoreOrderRepository(storeHandler)
	orderHandler := usecases.NewOrderHandler(repository)
	webservice := interfaces.NewWebservice(orderHandler, logger)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		webservice.CreateOrder(w, r)
	})

	logger.Info("Server listening on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if errors.Is(err, http.ErrServerClosed) {
		logger.Info("Server closed")
	} else if err != nil {
		logger.Error("Server failed: %s", err)
		os.Exit(1)
	}
}
