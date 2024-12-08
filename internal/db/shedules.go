package db

import (
    "time"
)

// Структура, представляющая мастера
type Master struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Структура для данных о слоте
type SlotData struct {
	ID              int       `json:"id"`
	MasterID        int       `json:"master_id"`
	Date            time.Time `json:"date"`
	Time            string    `json:"time"`
    IsAvailable     bool      `json:"is_available"`
    ServiceID       *int      `json:"service_id"` // Указатель на int, чтобы поддерживать NULL
}

// Функция для получения всех мастеров из базы данных
func GetMasters() ([]Master, error) {
	var masters []Master
	query := "SELECT id, name FROM masters" // Пример запроса, адаптируйте под вашу схему БД

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var master Master
		if err := rows.Scan(&master.ID, &master.Name); err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}

	return masters, nil
}

// Функция для получения доступных слотов для мастера
func GetAvailableSlotsForMaster(masterID int) ([]SlotData, error) {
	var slots []SlotData
	query := "SELECT id, master_id, service_id, date, time FROM schedules WHERE master_id = $1 AND is_available = true"

	rows, err := DB.Query(query, masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot SlotData
		if err := rows.Scan(&slot.ID, &slot.MasterID, &slot.ServiceID, &slot.Date, &slot.Time); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	return slots, nil
}

