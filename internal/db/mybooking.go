package db

import (
	"fmt"
)

// Структура для записи клиента
type Booking struct {
	ID          int    `json:"id"`
	ServiceID   int    `json:"service_id"`
	ScheduleID  int    `json:"schedule_id"`
	CustomerName string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
}

// Получение всех записей для клиента по телефону
func GetBookingsByPhone(phone string) ([]Booking, error) {
	rows, err := DB.Query(`
		SELECT id, service_id, schedule_id, customer_name, customer_phone
		FROM bookings
		WHERE customer_phone = $1`, phone)
	if err != nil {
		return nil, fmt.Errorf("Error fetching bookings for phone %s: %v", phone, err)
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		err := rows.Scan(&booking.ID, &booking.ServiceID, &booking.ScheduleID, &booking.CustomerName, &booking.CustomerPhone)
		if err != nil {
			return nil, fmt.Errorf("Error scanning booking row: %v", err)
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return bookings, nil
}
