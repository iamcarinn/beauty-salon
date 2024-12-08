// handlers/shedules_handler.go
package handlers

import (
	"fmt"
	"html/template"
	"net/http"
    "beauty-salon/internal/db"
	"github.com/gorilla/mux"
	"strconv"
)

// Структура данных для передачи в шаблон
type PageData struct {
	MasterSlots []struct {
		Master db.Master
		Slots  []db.SlotData
	}
	ServiceID int
}

func ViewAvailableSlotsHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр service_id из URL
	vars := mux.Vars(r)
	serviceID, err := strconv.Atoi(vars["service_id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid service_id: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем всех мастеров из базы данных
	masters, err := db.GetMasters()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching masters: %v", err), http.StatusInternalServerError)
		return
	}

	var masterSlots []struct {
		Master db.Master
		Slots  []db.SlotData
	}

	// Получаем слоты для каждого мастера, фильтруя по service_id
	for _, master := range masters {
		slots, err := db.GetAvailableSlotsForMaster(master.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching slots for master %d: %v", master.ID, err), http.StatusInternalServerError)
			return
		}

		// Применяем фильтрацию по service_id, если оно нужно
		var filteredSlots []db.SlotData
		for _, slot := range slots {
			// Проверяем, если ServiceID не является NULL, то фильтруем по service_id
			if slot.ServiceID != nil && *slot.ServiceID == serviceID {
				filteredSlots = append(filteredSlots, slot)
			} else if slot.ServiceID == nil {
				// Если ServiceID равно NULL, пропускаем этот слот
				filteredSlots = append(filteredSlots, slot)
			}
		}

		masterSlots = append(masterSlots, struct {
			Master db.Master
			Slots  []db.SlotData
		}{
			Master: master,
			Slots:  filteredSlots,
		})
	}

	// Отправляем данные в шаблон
	data := PageData{
		MasterSlots: masterSlots,
		ServiceID:   serviceID,
	}

	tmpl, err := template.ParseFiles("templates/booking.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}
