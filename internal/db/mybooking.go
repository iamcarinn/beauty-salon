package db

import (
    "fmt"
    "time"
)

// Обновленная структура для отображения информации о записи с категорией
type BookingInfo struct {
    Date          string `json:"date"`
    Time          string `json:"time"`
    MasterName    string `json:"master_name"`
    ServiceName   string `json:"service_name"`
    CategoryName  string `json:"category_name"`
    CustomerName  string `json:"customer_name"`
    CustomerPhone string `json:"customer_phone"`
    ScheduleID    int    `json:"schedule_id"`
    BookingID     int    `json:"booking_id"`
}

// Функция для форматирования даты в формат "DD.MM.YYYY"
func (b *BookingInfo) FormattedDate() string {
    parsedDate, err := time.Parse("2006-01-02T15:04:05Z", b.Date)
    if err != nil {
        return b.Date // если не удалось распарсить, возвращаем исходную дату
    }
    return parsedDate.Format("02.01.2006")
}

// Функция для форматирования времени в формат "HH:MM"
func (b *BookingInfo) FormattedTime() string {
    parsedTime, err := time.Parse("2006-01-02T15:04:05Z", b.Time)
    if err != nil {
        return b.Time // если не удалось распарсить, возвращаем исходное время
    }
    return parsedTime.Format("15:04")
}

// Обновленная функция для получения записей по номеру телефона с категорией
func GetBookingsByPhone(phone string) ([]BookingInfo, error) {
    query := `
        SELECT 
            s.date,
            s.time,
            m.name AS master_name,
            srv.name AS service_name,
            cat.name AS category_name,
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
        JOIN
            categories cat ON srv.category_id = cat.id
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
            &booking.CategoryName,
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
