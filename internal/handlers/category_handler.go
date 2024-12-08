package handlers

import (
    "beauty-salon/internal/db"
    "beauty-salon/internal/models"
    "html/template"
    "net/http"
    "strconv"
)

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
    categoryIDStr := r.URL.Query().Get("category_id")
    categoryID, err := strconv.Atoi(categoryIDStr)
    if err != nil {
        http.Error(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    services, err := db.GetServicesByCategory(categoryID)
    if err != nil {
        http.Error(w, "Failed to get services", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/category.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    data := struct {
        Services []models.Service
    }{
        Services: services,
    }

    tmpl.Execute(w, data)
}
