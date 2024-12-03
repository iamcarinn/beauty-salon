// handlers/slot_handler.go
package handlers

import (
    "beauty-salon/internal/db"
    "html/template"
    "net/http"
    "strconv"
)

func SlotHandler(w http.ResponseWriter, r *http.Request) {
    serviceIDStr := r.URL.Query().Get("service_id")
    serviceID, err := strconv.Atoi(serviceIDStr)
    if err != nil {
        http.Error(w, "Invalid service ID", http.StatusBadRequest)
        return
    }

    slots, err := db.GetAvailableSlots(serviceID)
    if err != nil {
        http.Error(w, "Failed to get available slots", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/slots.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    data := struct {
        ServiceID int
        Slots     []db.Slot
    }{
        ServiceID: serviceID,
        Slots:     slots,
    }

    tmpl.Execute(w, data)
}
