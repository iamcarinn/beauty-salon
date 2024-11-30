package main

import (
	"beauty-salon/internal/db"
    "log"
)

func main() {
    db.InitDB()
    defer db.DB.Close()
    log.Println("Server is running...")
}
