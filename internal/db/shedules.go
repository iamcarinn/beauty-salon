package db

import (
    "fmt"
)

// Структура для мастера
type Master struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CategoryID int  `json:"category_id"`
}

// Структура для слота
type Slot struct {
	ID         int    `json:"id"`
	MasterID   int    `json:"master_id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	IsAvailable bool  `json:"is_available"`
}

// Структура для услуги
type Service struct {
	ID          int    `json:"id"`
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	Price       float64 `json:"price"`
}

// Структура для категории
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
}

// Получение мастеров для услуги
func GetMastersForService(serviceID int) ([]Master, error) {
    category, _ := GetCategoryForService(serviceID)
    categoryID := category.ID
	rows, err := DB.Query(`
		SELECT id, name, category_id
		FROM masters
		WHERE category_id = $1`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching masters for category %d: %v", categoryID, err)
	}
	defer rows.Close()

	var masters []Master
	for rows.Next() {
		var master Master
		if err := rows.Scan(&master.ID, &master.Name, &master.CategoryID); err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return masters, nil
}

// Получение категории для сервиса
func GetCategoryForService(serviceID int) (Category, error) {
	var category Category

	// Выполняем запрос для получения категории
	err := DB.QueryRow(`
        SELECT c.id, c.name
        FROM categories c
        JOIN services s ON c.id = s.category_id
        WHERE s.id = $1`, serviceID).Scan(&category.ID, &category.Name)
	
	if err != nil {
		return Category{}, fmt.Errorf("Error fetching category for service %d: %v", serviceID, err)
	}

	return category, nil
}

// Получение доступных слотов для мастера
func GetAvailableSlotsForMaster(masterID int) ([]Slot, error) {
	rows, err := DB.Query(`
		SELECT id, master_id, date, time, is_available
		FROM schedules
		WHERE master_id = $1 AND is_available = true`, masterID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching slots for master %d: %v", masterID, err)
	}
	defer rows.Close()

	var slots []Slot
	for rows.Next() {
		var slot Slot
		if err := rows.Scan(&slot.ID, &slot.MasterID, &slot.Date, &slot.Time, &slot.IsAvailable); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}



// Получение доступных дат для мастера
func GetAvailableDateForMaster(masterID int) ([]string, error) {
	rows, err := DB.Query(`
		SELECT date
		FROM schedules
		WHERE master_id = $1 AND is_available = true`, masterID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching slots for master %d: %v", masterID, err)
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
        var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dates, nil
}

// Получение доступного времени для мастера и даты
func GetAvailableTimeForMasterAndDate(masterID int, date string) ([]string, error) {
	rows, err := DB.Query(`
		SELECT time
		FROM schedules
		WHERE master_id = $1 AND is_available = true AND date = $2`, masterID, date)
	if err != nil {
		return nil, fmt.Errorf("Error fetching slots for master %d and date %s: %v", masterID, date, err)
	}
	defer rows.Close()

	var times []string
	for rows.Next() {
        var time string
		if err := rows.Scan(&time); err != nil {
			return nil, err
		}
		times = append(times, time)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return times, nil
}
