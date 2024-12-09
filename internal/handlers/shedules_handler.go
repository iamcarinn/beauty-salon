package handlers

import (
	"fmt"
	"beauty-salon/internal/db"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"strconv"
	"encoding/json"
)

// Структура для отображения мастеров и слотов
type MasterSlots struct {
	Master db.Master
	Slots  []db.Slot
}

// Обработчик для отображения мастеров в выбранной категории и услуг
func ViewMastersHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем service_id из URL
	serviceID, err := strconv.Atoi(mux.Vars(r)["service_id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid service ID: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем мастеров для выбранной услуги
	masters, err := db.GetMastersForService(serviceID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching masters for service %d: %v", serviceID, err), http.StatusInternalServerError)
		return
	}

	// Загружаем шаблон
	tmpl, err := template.ParseFiles("templates/booking.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	// Передаем мастеров и услуги в шаблон
	err = tmpl.Execute(w, struct {
		Masters  []db.Master
		ServiceID int
	}{
		Masters:  masters,
		ServiceID: serviceID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

// Handler для получения доступных дат
func ViewAvailableDatesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	masterID, err := strconv.Atoi(vars["masterID"])
	if err != nil {
		http.Error(w, "Invalid master ID", http.StatusBadRequest)
		return
	}

	dates, _ := db.GetAvailableDateForMaster(masterID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"dates": dates})
}

// Handler для получения доступного времени
func ViewAvailableTimesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	masterID, err := strconv.Atoi(vars["masterID"])
	if err != nil {
		http.Error(w, "Invalid master ID", http.StatusBadRequest)
		return
	}

	date := vars["date"]
	times, _ := db.GetAvailableTimeForMasterAndDate(masterID, date)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"times": times})
}