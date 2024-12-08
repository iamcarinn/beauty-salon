package handlers

import (
	"beauty-salon/internal/db"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SlotData struct {
	MasterID int
	Date     string
	Time     string
}

type SlotViewData struct {
	Slots []SlotData
}

func ViewAvailableSlotsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем service_id из URL
	vars := mux.Vars(r)
	serviceID, err := strconv.Atoi(vars["service_id"])
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Получаем category_id для данной услуги
	categoryID, err := db.GetCategoryIDByService(serviceID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching category ID: %v", err), http.StatusInternalServerError)
		return
	}

	// Получаем доступные слоты для категории
	dbSlots, err := db.GetAvailableSlotsForCategory(categoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching slots: %v", err), http.StatusInternalServerError)
		return
	}

	// Подготавливаем данные для шаблона
	var slots []SlotData
	for _, dbSlot := range dbSlots {
		slots = append(slots, SlotData{
			MasterID: dbSlot.MasterID,
			Date:     dbSlot.Date,
			Time:     dbSlot.Time,
		})
	}

	data := SlotViewData{
		Slots: slots,
	}

	// Загружаем и обрабатываем шаблон
	tmpl, err := template.ParseFiles("templates/booking.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	// Рендерим шаблон
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}
