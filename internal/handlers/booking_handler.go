package handlers

import (
	"beauty-salon/internal/db"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookingData struct {
    ServiceID    int
    Masters      []db.Master
    Slots        []db.Schedule
    CustomerName string
    CustomerPhone string
    MasterID     int
    SlotID       int
}

func BookServiceHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serviceID := vars["service_id"]

    // Преобразуем serviceID в int
    serviceIDInt, err := strconv.Atoi(serviceID)
    if err != nil {
        http.Error(w, "Invalid service ID", http.StatusBadRequest)
        return
    }

    // Получаем мастеров и доступные слоты для услуги
    masters, err := db.GetAvailableMastersForService(serviceIDInt)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching masters: %v", err), http.StatusInternalServerError)
        return
    }
    
    // Логируем полученных мастеров
    log.Printf("Masters: %+v", masters)

    slots, err := db.GetAvailableSlotsForService(serviceIDInt)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching slots: %v", err), http.StatusInternalServerError)
        return
    }

    // Логируем доступные слоты
    log.Printf("Slots: %+v", slots)

    if r.Method == http.MethodPost {
        // Извлекаем данные из формы
        masterID, err := strconv.Atoi(r.FormValue("master_id"))
        if err != nil {
            http.Error(w, "Invalid master ID", http.StatusBadRequest)
            return
        }

        slotID, err := strconv.Atoi(r.FormValue("slot_id"))
        if err != nil {
            http.Error(w, "Invalid slot ID", http.StatusBadRequest)
            return
        }

        customerName := r.FormValue("customer_name")
        customerPhone := r.FormValue("customer_phone")

        // Создаем новое бронирование в базе данных
        err = db.CreateBooking(serviceIDInt, masterID, slotID, customerName, customerPhone)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error creating booking: %v", err), http.StatusInternalServerError)
            return
        }

        // Перенаправляем на страницу подтверждения
        http.Redirect(w, r, "/booking/confirmation", http.StatusSeeOther)
        return
    }

    // Подготавливаем данные для передачи в шаблон
    data := BookingData{
        ServiceID: serviceIDInt,
        Masters:   masters,
        Slots:     slots,
    }

    tmpl, err := template.New("booking").ParseFiles("templates/booking.html")
    if err != nil {
        log.Printf("Error loading template: %v", err)
        http.Error(w, "Internal Server Error: Template not found or invalid", http.StatusInternalServerError)
        return
    }


    // Рендерим шаблон с данными
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
    }
}

// Страница подтверждения бронирования
func BookingConfirmationHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Ваше бронирование успешно выполнено!")
}
