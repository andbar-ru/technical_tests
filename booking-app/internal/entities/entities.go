package entities

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var t time.Time
	err := t.UnmarshalJSON(data)
	if err == nil {
		d.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	}
	if _, ok := err.(*time.ParseError); !ok {
		return err
	}
	s := strings.Trim(string(data), `"`)
	t, err = time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type Order struct {
	ID        string `json:"id"`
	HotelID   string `json:"hotel_id"`
	RoomID    string `json:"room_id"`
	UserEmail string `json:"email"`
	From      Date   `json:"from"`
	To        Date   `json:"to"`
}

type RoomAvailability struct {
	ID      string `json:"id"`
	HotelID string `json:"hotel_id"`
	RoomID  string `json:"room_id"`
	Date    Date   `json:"date"`
	Quota   int    `json:"quota"`
}

func (o Order) Validate() []error {
	var errs []error
	var zeroOrder Order
	if o.ID == zeroOrder.ID {
		errs = append(errs, AppError{"Field 'id' is not specified"})
	}
	if o.HotelID == zeroOrder.HotelID {
		errs = append(errs, AppError{"Field 'hotel_id' is not specified"})
	}
	if o.RoomID == zeroOrder.RoomID {
		errs = append(errs, AppError{"Field 'room_id' is not specified"})
	}
	if o.UserEmail == zeroOrder.UserEmail {
		errs = append(errs, AppError{"Field 'email' is not specified"})
	}
	if o.From == zeroOrder.From {
		errs = append(errs, AppError{"Field 'from' is not specified"})
	}
	if o.To == zeroOrder.To {
		errs = append(errs, AppError{"Field 'to' is not specified"})
	}
	if o.From.After(o.To.Time) {
		errs = append(errs, AppError{"Date 'from' is later than date 'to'"})
	}
	return errs
}

func (o Order) Days() []Date {
	if o.From.After(o.To.Time) {
		return nil
	}
	days := make([]Date, 0)
	curDay := o.From
	for !curDay.After(o.To.Time) {
		days = append(days, curDay)
		curDay = Date{Time: curDay.AddDate(0, 0, 1)}
	}
	return days
}

type AppError struct {
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
