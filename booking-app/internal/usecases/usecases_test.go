package usecases_test

import (
	"sync"
	"testing"

	"booking_app/internal/entities"
	"booking_app/internal/infrastructure"
	"booking_app/internal/interfaces"
	"booking_app/internal/usecases"
	"booking_app/internal/utils"
)

var availabilityIds = map[entities.Date]string{
	utils.Date(2024, 1, 1): utils.Uuid(),
	utils.Date(2024, 1, 2): utils.Uuid(),
	utils.Date(2024, 1, 3): utils.Uuid(),
	utils.Date(2024, 1, 4): utils.Uuid(),
	utils.Date(2024, 1, 5): utils.Uuid(),
}

func NewStoreHandler() *infrastructure.StoreHandler {
	var orders = []entities.Order{}
	var availability = []entities.RoomAvailability{
		{ID: availabilityIds[utils.Date(2024, 1, 1)], HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 1), Quota: 1},
		{ID: availabilityIds[utils.Date(2024, 1, 2)], HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 2), Quota: 1},
		{ID: availabilityIds[utils.Date(2024, 1, 3)], HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 3), Quota: 1},
		{ID: availabilityIds[utils.Date(2024, 1, 4)], HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 4), Quota: 1},
		{ID: availabilityIds[utils.Date(2024, 1, 5)], HotelID: "reddison", RoomID: "lux", Date: utils.Date(2024, 1, 5), Quota: 0},
	}
	return &infrastructure.StoreHandler{
		Mux:              sync.Mutex{},
		Orders:           orders,
		RoomAvailability: availability,
	}
}

func TestOrderHandlerCreate(t *testing.T) {
	order1 := entities.Order{
		ID:      utils.Uuid(),
		HotelID: "reddison",
		RoomID:  "lux",
		From:    utils.Date(2023, 9, 11),
		To:      utils.Date(2023, 9, 15),
	}
	order2 := entities.Order{
		ID:      utils.Uuid(),
		HotelID: "reddison",
		RoomID:  "lux",
		From:    utils.Date(2024, 1, 1),
		To:      utils.Date(2024, 1, 6),
	}
	order3 := entities.Order{
		ID:      utils.Uuid(),
		HotelID: "reddison",
		RoomID:  "lux",
		From:    utils.Date(2024, 1, 4),
		To:      utils.Date(2024, 1, 4),
	}
	order4 := entities.Order{
		ID:      utils.Uuid(),
		HotelID: "reddison",
		RoomID:  "lux",
		From:    utils.Date(2024, 1, 2),
		To:      utils.Date(2024, 1, 4),
	}

	storeHandler := NewStoreHandler()
	repository := interfaces.NewStoreOrderRepository(storeHandler)
	orderHandler := usecases.NewOrderHandler(repository)

	want := "Hotel room is not available for selected dates"
	err := orderHandler.Create(order1)
	if err == nil || err.Error() != want {
		t.Errorf("order1: want error %q, got %s", want, err)
	}

	want = "Hotel room is not available for selected dates"
	err = orderHandler.Create(order2)
	if err == nil || err.Error() != want {
		t.Errorf("order2: want error %q, got %s", want, err)
	}

	err = orderHandler.Create(order3)
	if err != nil {
		t.Errorf("order3: want no error, got %q", err)
	}
	a := storeHandler.GetRoomAvailabilityById(availabilityIds[utils.Date(2024, 1, 4)])
	if a.Quota != 0 {
		t.Errorf("order3: want quota 0, got %d", a.Quota)
	}

	want = "Hotel room is not available for selected dates"
	err = orderHandler.Create(order4)
	if err == nil || err.Error() != want {
		t.Errorf("order4: want error %q, got %s", want, err)
	}
	// Test rollback
	a2 := storeHandler.GetRoomAvailabilityById(availabilityIds[utils.Date(2024, 1, 2)])
	if a2.Quota != 1 {
		t.Errorf("order4, a2: want quota 1, got %d", a2.Quota)
	}
	a3 := storeHandler.GetRoomAvailabilityById(availabilityIds[utils.Date(2024, 1, 3)])
	if a3.Quota != 1 {
		t.Errorf("order4, a3: want quota 1, got %d", a3.Quota)
	}

	storeHandler = NewStoreHandler()
	repository = interfaces.NewStoreOrderRepository(storeHandler)
	orderHandler = usecases.NewOrderHandler(repository)

	err = orderHandler.Create(order4)
	if err != nil {
		t.Errorf("order4: want no error, got %q", err)
	}
	a2 = storeHandler.GetRoomAvailabilityById(availabilityIds[utils.Date(2024, 1, 2)])
	if a2.Quota != 0 {
		t.Errorf("order4, a2: want quota 0, got %d", a.Quota)
	}
	a3 = storeHandler.GetRoomAvailabilityById(availabilityIds[utils.Date(2024, 1, 3)])
	if a3.Quota != 0 {
		t.Errorf("order4, a3: want quota 0, got %d", a.Quota)
	}
	a4 := storeHandler.GetRoomAvailabilityById(availabilityIds[utils.Date(2024, 1, 4)])
	if a4.Quota != 0 {
		t.Errorf("order4, a4: want quota 0, got %d", a.Quota)
	}

	want = "Hotel room is not available for selected dates"
	err = orderHandler.Create(order4)
	if err == nil || err.Error() != want {
		t.Errorf("order4: want error %q, got %s", want, err)
	}
}
