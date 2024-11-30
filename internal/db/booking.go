package db

import (
	//"beauty-salon/internal/models"
	//"database/sql"
	"fmt"
)

// BookService сохраняет информацию о записи на услугу
func BookService(serviceID int, customerName, customerPhone string) error {
	// Создание запроса для вставки новой записи о записи на услугу
	query := `INSERT INTO bookings (service_id, customer_name, customer_phone) 
			  VALUES ($1, $2, $3)`

	// Выполнение запроса
	_, err := DB.Exec(query, serviceID, customerName, customerPhone)
	if err != nil {
		return fmt.Errorf("failed to book service: %v", err)
	}

	return nil
}
