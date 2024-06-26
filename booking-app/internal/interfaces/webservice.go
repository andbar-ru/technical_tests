package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"booking_app/internal/entities"
	"booking_app/internal/usecases"
	"booking_app/internal/utils"
)

// A Webservice represents web service
type Webservice struct {
	// order handler
	OrderHandler *usecases.OrderHandler
	// logger
	Logger usecases.Logger
}

// NewWebservice returns a new web service provided by orderHandler and logger passed as arguments.
func NewWebservice(orderHandler *usecases.OrderHandler, logger usecases.Logger) *Webservice {
	return &Webservice{
		OrderHandler: orderHandler,
		Logger:       logger,
	}
}

// CreateOrder creates new order according to reqest body.
// If successed responses with new order and status 201.
// If failed responses with errors and some erroneous status.
func (service Webservice) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder entities.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = entities.AppError{Message: "Empty Body"}
		}
		service.handleErrors(w, []error{err}, http.StatusBadRequest)
		return
	}
	newOrder.ID = utils.Uuid()
	errs := newOrder.Validate()
	if len(errs) > 0 {
		service.handleErrors(w, errs, http.StatusBadRequest)
		return
	}

	err = service.OrderHandler.Create(newOrder)
	if err != nil {
		service.handleErrors(w, []error{err}, http.StatusConflict)
		err = fmt.Errorf("Order %v. %w", newOrder, err)
		service.Logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)

	service.Logger.Info("Order successfully created: %v", newOrder)
}

// handleErrors responses with given errors and some erroneous status.
func (service Webservice) handleErrors(w http.ResponseWriter, errs []error, status int) {
	w.Header().Set("Content-Type", "application/json")
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Error()
	}
	o := struct {
		Errors []string `json:"errors"`
	}{msgs}
	b, err := json.Marshal(o)
	if err != nil {
		http.Error(w, "Invalid error", status)
		return
	}
	http.Error(w, string(b), status)
}
