package handlers

import (
    "beauty-salon/internal/db"
    "html/template"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
)

func BookingHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)  // Извлекаем параметры из URL
    serviceIDStr := vars["service_id"]  // Параметр из URL
    serviceID, err := strconv.Atoi(serviceIDStr)
    if err != nil {
        http.Error(w, "Invalid service ID", http.StatusBadRequest)
        return
    }

    if r.Method == "POST" {
        customerName := r.FormValue("customer_name")
        customerPhone := r.FormValue("customer_phone")

        err := db.BookService(serviceID, customerName, customerPhone)
        if err != nil {
            http.Error(w, "Failed to book service", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/", http.StatusFound)
    }

    tmpl, err := template.ParseFiles("templates/booking.html")
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    data := struct {
        ServiceID int
    }{
        ServiceID: serviceID,
    }

    tmpl.Execute(w, data)
}
