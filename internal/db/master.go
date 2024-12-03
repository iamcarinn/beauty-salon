package db

import "fmt"

type Master struct {
	ID         int
	Name       string
	CategoryID int // ID категории
}

func GetAvailableMastersForService(serviceID int) ([]Master, error) {
	// Получаем доступных мастеров для выбранной услуги и доступных слотов, а также категорию мастера
	query := `
        SELECT m.id, m.name, m.category_id
        FROM masters m
        JOIN schedules s ON s.master_id = m.id
        JOIN services srv ON srv.id = s.service_id
        JOIN categories c ON c.id = m.category_id
        WHERE srv.id = $1 AND s.is_available = TRUE;
    `
	rows, err := DB.Query(query, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch masters: %v", err)
	}
	defer rows.Close()

	var masters []Master
	for rows.Next() {
		var master Master
		if err := rows.Scan(&master.ID, &master.Name, &master.CategoryID); err != nil {
			return nil, fmt.Errorf("failed to scan master row: %v", err)
		}
		masters = append(masters, master)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over master rows: %v", err)
	}

	return masters, nil
}
