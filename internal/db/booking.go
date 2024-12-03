package db

import (
    "fmt"
)

func BookService(serviceID, scheduleID int, customerName, customerPhone string) error {
    query := `
        INSERT INTO bookings (service_id, schedule_id, customer_name, customer_phone)
        VALUES ($1, $2, $3, $4)
    `
    _, err := DB.Exec(query, serviceID, scheduleID, customerName, customerPhone)
    if err != nil {
        return fmt.Errorf("failed to book service: %v", err)
    }

    // Обновляем статус слота на занятый
    _, err = DB.Exec("UPDATE schedules SET is_available = FALSE WHERE id = $1", scheduleID)
    if err != nil {
        return fmt.Errorf("failed to update slot availability: %v", err)
    }

    return nil
}

func CreateBooking(serviceID, masterID, slotID int, customerName, customerPhone string) error {
    query := `
        INSERT INTO bookings (service_id, schedule_id, customer_name, customer_phone, master_id)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := DB.Exec(query, serviceID, slotID, customerName, customerPhone, masterID)
    if err != nil {
        return fmt.Errorf("failed to create booking: %v", err)
    }
    return nil
}

