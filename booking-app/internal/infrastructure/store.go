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

type StoreHandler struct {
	Mux              sync.Mutex
	Orders           []entities.Order
	RoomAvailability []entities.RoomAvailability
}

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

func (h *StoreHandler) GetOrders() []entities.Order {
	return h.Orders
}

func (h *StoreHandler) GetRoomAvailability() []entities.RoomAvailability {
	return h.RoomAvailability
}

func (h *StoreHandler) AddOrder(order entities.Order) {
	h.Mux.Lock()
	defer h.Mux.Unlock()

	h.Orders = append(h.Orders, order)
}

func (h *StoreHandler) GetRoomAvailabilityById(id string) entities.RoomAvailability {
	for _, item := range h.RoomAvailability {
		if item.ID == id {
			return item
		}
	}
	var zero entities.RoomAvailability
	return zero
}

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
