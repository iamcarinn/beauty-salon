package db

import "fmt"

type Schedule struct {
    ID          int
    MasterID    int
    Date        string
    Time        string
    IsAvailable bool
    ServiceID   int
}


func GetAvailableSlotsForService(serviceID int) ([]Schedule, error) {
    query := `
        SELECT s.id, s.master_id, s.date, s.time, s.is_available, s.service_id
        FROM schedules s
        WHERE s.is_available = TRUE
        AND s.date >= CURRENT_DATE -- или любое другое условие для фильтрации
        AND s.service_id = $1 -- Фильтрация по выбранной услуге
    `

    rows, err := DB.Query(query, serviceID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch slots: %v", err)
    }
    defer rows.Close()

    var slots []Schedule
    for rows.Next() {
        var slot Schedule
        if err := rows.Scan(&slot.ID, &slot.MasterID, &slot.Date, &slot.Time, &slot.IsAvailable, &slot.ServiceID); err != nil {
            return nil, fmt.Errorf("failed to scan slot: %v", err)
        }
        slots = append(slots, slot)
    }

    return slots, nil
}
