package router

import (
	"beauty-salon/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Роут для главной страницы
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Роут для категории
	r.HandleFunc("/category", handlers.CategoryHandler).Methods("GET")

	// Роут для страницы бронирования с параметром service_id
	r.HandleFunc("/booking/{service_id}", handlers.ViewAvailableSlotsHandler).Methods("GET")

	return r
}
