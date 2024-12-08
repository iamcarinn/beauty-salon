package router

import (
    "beauty-salon/internal/handlers"
    "github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
    r.HandleFunc("/category", handlers.CategoryHandler).Methods("GET")
    r.HandleFunc("/booking/{service_id}", handlers.ViewAvailableSlotsHandler).Methods("GET")

    return r
}
