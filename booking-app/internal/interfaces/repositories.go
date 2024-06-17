package interfaces

import (
	"fmt"

	"booking_app/internal/entities"
	"booking_app/internal/usecases"
)

// A StoreHandler represents interface for store.
type StoreHandler interface {
	GetOrders() []entities.Order
	GetRoomAvailability() []entities.RoomAvailability
	AddOrder(order entities.Order)
	GetRoomAvailabilityById(id string) entities.RoomAvailability
	UpdateRoomAvailability(id string, roomAvailability entities.RoomAvailability) error
}

// A StoreRepository represents repository with StoreHandler as its store handler.
type StoreRepository struct {
	// store handler
	storeHandler StoreHandler
}

type StoreOrderRepository StoreRepository

// NewStoreOrderRepository returns a new StoreOrderRepository provided by storeHandler passed as argument.
func NewStoreOrderRepository(storeHandler StoreHandler) *StoreOrderRepository {
	return &StoreOrderRepository{
		storeHandler: storeHandler,
	}
}

// HasTransactions returns false as StoreOrderRepository doesn't support transactions.
func (repo *StoreOrderRepository) HasTransactions() bool {
	return false
}

// BeginTransaction returns nil as StoreOrderRepository doesn't support transactions.
func (repo *StoreOrderRepository) BeginTransaction() usecases.Transaction {
	return nil
}

// AddOrder adds an order to store.
func (repo *StoreOrderRepository) AddOrder(order entities.Order) error {
	repo.storeHandler.AddOrder(order)
	return nil
}

// FindOrdersByEmail searches orders by user email and returns found or empty slice if could not find.
func (repo *StoreOrderRepository) FindOrdersByEmail(email string) []entities.Order {
	allOrders := repo.storeHandler.GetOrders()
	var orders []entities.Order
	for _, order := range allOrders {
		if order.UserEmail == email {
			orders = append(orders, order)
		}
	}
	return orders
}

// ChangeRoomAvailabilityQuota increments or decrements room availability, specified by its id,
// quota by delta. If no items with given id or quota is not enough, returns AppError.
func (repo *StoreOrderRepository) ChangeRoomAvailabilityQuota(id string, delta int) error {
	var zeroRoomAvailability entities.RoomAvailability
	roomAvailability := repo.storeHandler.GetRoomAvailabilityById(id)
	if roomAvailability == zeroRoomAvailability {
		return entities.AppError{Message: fmt.Sprintf("Could not find roomAvailability with id='%s'.", id)}
	}
	if roomAvailability.Quota+delta < 0 {
		return entities.AppError{Message: fmt.Sprintf("RoomAvailability with id='%s' has no enough quota.", id)}
	}
	roomAvailability.Quota += delta
	repo.storeHandler.UpdateRoomAvailability(id, roomAvailability)
	return nil
}

// FindRoomAvailabilityByHotelAndRoom searches room availability items by given hotelID and roomID.
// For convenience returns map of date to room availability item.
func (repo *StoreOrderRepository) FindRoomAvailabilityByHotelAndRoom(hotelID, roomID string) map[entities.Date]entities.RoomAvailability {
	roomAvailability := repo.storeHandler.GetRoomAvailability()
	m := make(map[entities.Date]entities.RoomAvailability)
	for _, item := range roomAvailability {
		if item.HotelID == hotelID && item.RoomID == roomID {
			m[item.Date] = item
		}
	}
	return m
}
