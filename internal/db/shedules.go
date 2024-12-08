package db

import (
	"log"
)

type SlotData struct {
	MasterID int
	Date     string
	Time     string
}

func GetCategoryIDByService(serviceID int) (int, error) {
	query := `SELECT category_id FROM services WHERE id = $1`
	var categoryID int
	err := DB.QueryRow(query, serviceID).Scan(&categoryID)
	if err != nil {
		return 0, err
	}
	return categoryID, nil
}


func GetAvailableSlotsForCategory(categoryID int) ([]SlotData, error) {
	query := `
		SELECT schedules.master_id, schedules.date, schedules.time
		FROM schedules
		JOIN masters ON schedules.master_id = masters.id
		WHERE masters.category_id = $1 AND schedules.is_available = true;
	`

	rows, err := DB.Query(query, categoryID)
	if err != nil {
		log.Printf("Error querying available slots: %v", err)
		return nil, err
	}
	defer rows.Close()

	var slots []SlotData
	for rows.Next() {
		var slot SlotData
		if err := rows.Scan(&slot.MasterID, &slot.Date, &slot.Time); err != nil {
			log.Printf("Error scanning slot data: %v", err)
			return nil, err
		}
		slots = append(slots, slot)
	}

	return slots, nil
}
