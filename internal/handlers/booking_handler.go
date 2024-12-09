package handlers

import (
	"beauty-salon/internal/db"
	"fmt"
	"net/http"
	"time"
)

// package handlers

// import (
//     "beauty-salon/internal/db"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	_ "github.com/lib/pq" // PostgreSQL driver
// )

func HandleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Извлекаем данные из формы
		masterID := r.FormValue("master_id")
		fmt.Printf("Master_id: %d", masterID)
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
            VALUES ($1, $2, $3, $4);`, 1, scheduleID, customerName, customerPhone) // Предполагаем, что service_id равно 1

		if err != nil {
			fmt.Fprintf(w, "Error inserting booking: %v", err)
			return
		}

		fmt.Fprintf(w, "Booking successful!")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
