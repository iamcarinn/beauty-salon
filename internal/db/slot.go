// internal/db/slot.go
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
        SELECT s.id, s.time, s.is_available
        FROM schedules s
        JOIN services sv ON sv.id = s.service_id
        WHERE sv.id = $1 AND s.is_available = TRUE
        ORDER BY s.time
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
