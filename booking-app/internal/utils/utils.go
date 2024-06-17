package utils

import (
	"crypto/rand"
	"fmt"
	"time"

	"booking_app/internal/entities"
)

// Date returns a Date object by given year, month and day.
func Date(year, month, day int) entities.Date {
	return entities.Date{Time: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
}

// Uuid return new uuid v4.
func Uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Unexpected error
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
