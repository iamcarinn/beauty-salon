package db

import (
    "beauty-salon/internal/models"
)

// Получение всех услуг
func GetAllServices() ([]models.Service, error) {
    rows, err := DB.Query("SELECT id, category_id, name, price FROM services")
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

// Услуги в зависимости от категории
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

