package interfaces

import (
	"fmt"

	"booking_app/internal/entities"
	"booking_app/internal/usecases"
)

type StoreHandler interface {
	GetOrders() []entities.Order
	GetRoomAvailability() []entities.RoomAvailability
	AddOrder(order entities.Order)
	GetRoomAvailabilityById(id string) entities.RoomAvailability
	UpdateRoomAvailability(id string, roomAvailability entities.RoomAvailability) error
}

type StoreRepository struct {
	storeHandler StoreHandler
}

func NewStoreRepository(storeHandler StoreHandler) *StoreRepository {
	return &StoreRepository{
		storeHandler: storeHandler,
	}
}

type StoreOrderRepository StoreRepository

func NewStoreOrderRepository(storeHandler StoreHandler) *StoreOrderRepository {
	return &StoreOrderRepository{
		storeHandler: storeHandler,
	}
}

func (repo *StoreOrderRepository) HasTransactions() bool {
	return false
}

func (repo *StoreOrderRepository) BeginTransaction() usecases.Transaction {
	return nil
}

func (repo *StoreOrderRepository) AddOrder(order entities.Order) error {
	repo.storeHandler.AddOrder(order)
	return nil
}

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
