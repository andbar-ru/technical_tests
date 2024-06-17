package infrastructure

import (
	"fmt"
	"sync"

	"booking_app/internal/entities"
	"booking_app/internal/utils"
)

var orders = []entities.Order{}

var availability = []entities.RoomAvailability{
	{ID: utils.Uuid(), HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 1), Quota: 1},
	{ID: utils.Uuid(), HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 2), Quota: 1},
	{ID: utils.Uuid(), HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 3), Quota: 1},
	{ID: utils.Uuid(), HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 4), Quota: 1},
	{ID: utils.Uuid(), HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 5), Quota: 0},
}

var storeHandler *StoreHandler

// StoreHandler represents realization of StoreHandler interface.
type StoreHandler struct {
	// Mutex for concurrent-safe write operations.
	Mux sync.Mutex
	// orders
	Orders []entities.Order
	// room availability items
	RoomAvailability []entities.RoomAvailability
}

// NewStoreHandler return a new instance of StoreHandler if it doesn't exist else returns existing instance.
func NewStoreHandler() *StoreHandler {
	if storeHandler != nil {
		return storeHandler
	}
	storeHandler = &StoreHandler{
		Mux:              sync.Mutex{},
		Orders:           orders,
		RoomAvailability: availability,
	}
	return storeHandler
}

// GetOrders return all orders.
func (h *StoreHandler) GetOrders() []entities.Order {
	return h.Orders
}

// GetRoomAvailability returns all availability items.
func (h *StoreHandler) GetRoomAvailability() []entities.RoomAvailability {
	return h.RoomAvailability
}

// AddOrder appends new order to existing orders.
func (h *StoreHandler) AddOrder(order entities.Order) {
	h.Mux.Lock()
	defer h.Mux.Unlock()

	h.Orders = append(h.Orders, order)
}

// GetRoomAvailabilityById returns availability item with id specified by `id` argument.
func (h *StoreHandler) GetRoomAvailabilityById(id string) entities.RoomAvailability {
	for _, item := range h.RoomAvailability {
		if item.ID == id {
			return item
		}
	}
	var zero entities.RoomAvailability
	return zero
}

// UpdateRoomAvailability replaces availability item, specified by `id` argument with new item,
// specified by roomAvailability argument. Return AppError, if store doesn't contain availability
// item with given id.
func (h *StoreHandler) UpdateRoomAvailability(id string, roomAvailability entities.RoomAvailability) error {
	h.Mux.Lock()
	defer h.Mux.Unlock()

	for i, item := range h.RoomAvailability {
		if item.ID == id {
			h.RoomAvailability[i] = roomAvailability
			return nil
		}
	}
	return entities.AppError{Message: fmt.Sprintf("Could not find roomAvailability with id '%s'", id)}
}
