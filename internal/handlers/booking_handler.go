package handlers

import (
	"beauty-salon/internal/db"
	"fmt"
	"net/http"
	"strconv"
	"html/template"
	"time"

	"github.com/gorilla/mux"
)


func HandleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
        // Получаем service_id из URL
        vars := mux.Vars(r)
        serviceID, err := strconv.Atoi(vars["service_id"])
        if err != nil {
            fmt.Printf("invalid serviceID = %v", serviceID)
            http.Error(w, "Invalid service ID", http.StatusBadRequest)
            return
        }        

		// Извлекаем данные из формы
		masterID := r.FormValue("master_id")
		dateStr := r.FormValue("date") // дата в формате "DD.MM.YYYY"
		timeStr := r.FormValue("time") // время в формате "HH:MM"
		customerName := r.FormValue("name")
		customerPhone := r.FormValue("phone")

		// Преобразуем дату в формат Date (ISO 8601: YYYY-MM-DDT00:00:00Z)
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
		if err != nil {
			fmt.Fprintf(w, "Error parsing date: %v", err)
			return
		}
		formattedDate := parsedDate.Format("2006-01-02")

		// Преобразуем время в формат Time (HH:MM -> HH:MM:00)
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
		if err != nil {
			fmt.Fprintf(w, "Error parsing time: %v", err)
			return
		}
		formattedTime := parsedTime.Format("15:04:00")

		// 2. Обновляем таблицу schedules, чтобы сделать время недоступным
		var scheduleID int
		err = db.QueryRow(
			`UPDATE schedules
            SET is_available = false
            WHERE master_id = $1
              AND date = $2
              AND time = $3
            RETURNING id;`, masterID, formattedDate, formattedTime).Scan(&scheduleID)
		if err != nil {
			fmt.Fprintf(w, "Error updating schedule: %v", err)
			return
		}

		// 3. Вставляем запись в таблицу bookings
		_, err = db.Exec(
			`INSERT INTO bookings (service_id, schedule_id, customer_name, customer_phone)
            VALUES ($1, $2, $3, $4);`, serviceID, scheduleID, customerName, customerPhone) // Предполагаем, что service_id равно 1

		if err != nil {
			fmt.Fprintf(w, "Error inserting booking: %v", err)
			return
		}

		// Рендерим шаблон
		tmpl, err := template.ParseFiles("templates/booking_successful.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		
			// Передаем данные в шаблон (если нужно, можно передать дополнительные данные)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
