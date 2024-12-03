package handlers

import (
    "beauty-salon/internal/db"
    "html/template"
    "net/http"
    "strconv"
)

func MasterHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем service_id из URL-параметра
    serviceIDStr := r.URL.Query().Get("service_id")
    serviceID, err := strconv.Atoi(serviceIDStr)
    if err != nil {
        http.Error(w, "Invalid service ID", http.StatusBadRequest)
        return
    }

    // Получаем доступных мастеров для выбранной услуги
    masters, err := db.GetAvailableMastersForService(serviceID)
    if err != nil {
        http.Error(w, "Failed to get available masters", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/master.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    data := struct {
        ServiceID int
        Masters   []db.Master
    }{
        ServiceID: serviceID,
        Masters:   masters,
    }

    tmpl.Execute(w, data)
}
