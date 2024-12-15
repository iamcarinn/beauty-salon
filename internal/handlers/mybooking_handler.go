package handlers

import (
	"beauty-salon/internal/db"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ViewMyBookingsHandler(w http.ResponseWriter, r *http.Request) {
	var bookings []db.BookingInfo

	// Обрабатываем POST запрос
	if r.Method == "POST" {
		// Извлекаем номер телефона из формы
		phone := r.FormValue("phone")

		// Получаем записи по номеру телефона
		var err error
		bookings, err = db.GetBookingsByPhone(phone)
		if err != nil {
			http.Error(w, "Ошибка при получении записей: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Загружаем шаблон
	t, err := template.ParseFiles("templates/mybooking.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке страницы", http.StatusInternalServerError)
		return
	}

	// Отправляем данные на страницу (если записи найдены, передаем их в шаблон)
	t.Execute(w, bookings)
}

func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Извлекаем ID бронирования из формы
		bookingID := r.FormValue("booking_id")
		scheduleID := r.FormValue("schedule_id")

		// Преобразуем в нужные типы
		intBookingID, err := strconv.Atoi(bookingID)
		if err != nil {
			http.Error(w, "Неверный формат ID бронирования", http.StatusBadRequest)
			return
		}

		intScheduleID, err := strconv.Atoi(scheduleID)
		if err != nil {
			http.Error(w, "Неверный формат ID расписания", http.StatusBadRequest)
			return
		}

		// Обновляем состояние is_available на true для записи в schedules
		_, err = db.Exec(`
            UPDATE schedules
            SET is_available = true
            WHERE id = $1
        `, intScheduleID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при обновлении расписания: %v", err), http.StatusInternalServerError)
			return
		}

		// Удаляем запись из таблицы bookings
		_, err = db.Exec(`
            DELETE FROM bookings
            WHERE id = $1
        `, intBookingID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при удалении записи: %v", err), http.StatusInternalServerError)
			return
		}

		// Перенаправляем пользователя обратно на страницу с записями
		http.Redirect(w, r, "/mybooking", http.StatusSeeOther)
	}
}

func UpdatePhoneHandler(w http.ResponseWriter, r *http.Request) {
    // Устанавливаем заголовок Content-Type как JSON
    w.Header().Set("Content-Type", "application/json")

    // Извлекаем параметр bookingID из URL
    vars := mux.Vars(r)
    bookingID := vars["bookingID"]

    // Чтение тела запроса (новый номер телефона)
    var requestData struct {
        Phone string `json:"phone"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        log.Println("Ошибка при чтении тела запроса:", err)
        http.Error(w, `{"error": "Ошибка при чтении запроса"}`, http.StatusBadRequest)
        return
    }

    // Проверка, что номер телефона не пустой
    if requestData.Phone == "" {
        http.Error(w, `{"error": "Номер телефона не может быть пустым"}`, http.StatusBadRequest)
        return
    }

    // Обновление номера телефона в базе данных
    err = UpdatePhoneNumber(bookingID, requestData.Phone)
    if err != nil {
        log.Println("Ошибка при обновлении номера телефона:", err)
        http.Error(w, `{"error": "Произошла ошибка при обновлении номера телефона"}`, http.StatusInternalServerError)
        return
    }

    // Отправка успешного ответа в формате JSON
    response := map[string]bool{"success": true}
    json.NewEncoder(w).Encode(response)
}


// Функция для обновления номера телефона в базе данных
func UpdatePhoneNumber(bookingID, newPhone string) error {
    query := `UPDATE bookings SET customer_phone = $1 WHERE id = $2`

    // Выполнение запроса
    _, err := db.Exec(query, newPhone, bookingID)
    if err != nil {
        return fmt.Errorf("ошибка при обновлении номера: %w", err)
    }

    return nil
}