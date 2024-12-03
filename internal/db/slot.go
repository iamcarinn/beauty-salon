package db

import (
    "fmt"
    "time"
)

// Структура для хранения слота
type Slot struct {
    ID        int
    Time      time.Time
    Available bool
}

// Функция для получения доступных слотов для услуги
func GetAvailableSlots(serviceID int) ([]Slot, error) {
    query := `
        SELECT id, time, is_available
        FROM schedules
        WHERE is_available = TRUE AND service_id = $1 AND date >= CURRENT_DATE
    `
    
    rows, err := DB.Query(query, serviceID)
    if err != nil {
        return nil, fmt.Errorf("failed to get slots: %v", err)
    }
    defer rows.Close()

    var slots []Slot
    for rows.Next() {
        var slot Slot
        if err := rows.Scan(&slot.ID, &slot.Time, &slot.Available); err != nil {
            return nil, fmt.Errorf("failed to scan slot: %v", err)
        }
        slots = append(slots, slot)
    }

    return slots, nil
}
