package db

import (
	"database/sql"
	"log"
	//"fmt"
	"beauty-salon/internal/models"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=yourpassword dbname=beauty_salon sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connected to the database!")
}

func GetServicesByCategory(categoryID int) ([]models.Service, error) {
    query := `SELECT id, category_id, name, price FROM services WHERE category_id = $1`
    rows, err := DB.Query(query, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var services []models.Service
    for rows.Next() {
        var service models.Service
        if err := rows.Scan(&service.ID, &service.CategoryID, &service.Name, &service.Price); err != nil {
            return nil, err
        }
        services = append(services, service)
    }

    return services, nil
}