package handlers

import (
    "beauty-salon/internal/db"
    "net/http"
    "html/template"
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
