package handlers

import (
    "beauty-salon/internal/db"
    "beauty-salon/internal/models"
    //"encoding/json"
    "html/template"
    "net/http"
    "log"
    "strconv"
    "os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    cwd, err := os.Getwd()  // Получаем текущую рабочую директорию
    if err != nil {
        http.Error(w, "Failed to get current working directory", http.StatusInternalServerError)
        return
    }
    log.Println("Current working directory:", cwd)  // Выводим в лог текущую рабочую директорию

    tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}


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
