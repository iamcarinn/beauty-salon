package handlers

import (
    "beauty-salon/internal/db"
    "encoding/json"
    "net/http"
)

func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
    services, err := db.GetAllServices()
    if err != nil {
        http.Error(w, "Failed to get services", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(services)
}
