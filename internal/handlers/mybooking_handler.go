package handlers

import (
	"beauty-salon/internal/db"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
