package usecases

import (
	"fmt"
	"time"

	"booking_app/internal/entities"
)

type Logger interface {
	Info(format string, v ...any)
	Error(format string, v ...any)
}

type Transaction interface {
	Commit()
	Rollback()
}

type OrderRepository interface {
	HasTransactions() bool
	BeginTransaction() Transaction
	AddOrder(order entities.Order) error
	FindOrdersByEmail(email string) []entities.Order
	ChangeRoomAvailabilityQuota(id string, delta int) error
	FindRoomAvailabilityByHotelAndRoom(hotelID, roomID string) map[entities.Date]entities.RoomAvailability
}

type OrderHandler struct {
	Repository OrderRepository
}

func NewOrderHandler(repository OrderRepository) *OrderHandler {
	return &OrderHandler{
		Repository: repository,
	}
}

func (h *OrderHandler) Create(order entities.Order) error {
	days := order.Days()
	if len(days) == 0 {
		return entities.AppError{Message: "Order has no days"}
	}

	var tx Transaction
	if h.Repository.HasTransactions() {
		tx = h.Repository.BeginTransaction()
		defer tx.Rollback() // Works if called without tx.Commit() else noop
	}

	availabilityItems := h.Repository.FindRoomAvailabilityByHotelAndRoom(order.HotelID, order.RoomID)
	wantedAvailabilityItems := make(map[entities.Date]entities.RoomAvailability)

	for _, day := range days {
		availabilityItem, ok := availabilityItems[day]
		if !ok || availabilityItem.Quota < 1 {
			return entities.AppError{Message: "Hotel room is not available for selected dates"}
		}
		wantedAvailabilityItems[day] = availabilityItem

	}

	bookedIds := make([]string, 0, len(wantedAvailabilityItems))
	for _, availabilityItem := range wantedAvailabilityItems {
		err := h.Repository.ChangeRoomAvailabilityQuota(availabilityItem.ID, -1)
		if err != nil {
			if tx == nil {
				// Rollback by hand
				for _, id := range bookedIds {
					h.Repository.ChangeRoomAvailabilityQuota(id, 1)
				}
			}
			return entities.AppError{Message: fmt.Sprintf("Failed to book hotel room for %s. %s", availabilityItem.Date.Format(time.DateOnly), err)}
		} else {
			bookedIds = append(bookedIds, availabilityItem.ID)
		}
	}

	err := h.Repository.AddOrder(order)
	if err != nil {
		return entities.AppError{Message: "Failed to create order"}
	}

	if tx != nil {
		tx.Commit()
	}

	return nil
}
