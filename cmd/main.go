package main

import (
    "beauty-salon/internal/db"
    "beauty-salon/internal/router"
    "log"
    "net/http"
)

func main() {
    db.InitDB()
    defer db.DB.Close()

    r := router.SetupRouter()
    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
