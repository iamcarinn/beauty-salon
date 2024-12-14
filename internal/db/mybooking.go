package db

import (
    "fmt"
)

// Структура для отображения информации о записи
type BookingInfo struct {
    Date         string `json:"date"`
    Time         string `json:"time"`
    MasterName   string `json:"master_name"`
    ServiceName  string `json:"service_name"`
    CustomerName string `json:"customer_name"`
    CustomerPhone string `json:"customer_phone"`
    ScheduleID   int    `json:"schedule_id"`
    BookingID    int    `json:"booking_id"`
}

// Функция для получения записей по номеру телефона
func GetBookingsByPhone(phone string) ([]BookingInfo, error) {
    query := `
        SELECT 
            s.date,
            s.time,
            m.name AS master_name,
            srv.name AS service_name,
            b.customer_name,
            b.customer_phone,
            s.id AS schedule_id,
            b.id AS booking_id
        FROM 
            bookings b
        JOIN 
            schedules s ON b.schedule_id = s.id
        JOIN 
            masters m ON s.master_id = m.id
        JOIN 
            services srv ON b.service_id = srv.id
        WHERE 
            b.customer_phone = $1;
    `

    rows, err := DB.Query(query, phone)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }
    defer rows.Close()

    var bookings []BookingInfo
    for rows.Next() {
        var booking BookingInfo
        err := rows.Scan(
            &booking.Date,
            &booking.Time,
            &booking.MasterName,
            &booking.ServiceName,
            &booking.CustomerName,
            &booking.CustomerPhone,
            &booking.ScheduleID,
            &booking.BookingID,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }
        bookings = append(bookings, booking)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during row iteration: %v", err)
    }

    return bookings, nil
}
