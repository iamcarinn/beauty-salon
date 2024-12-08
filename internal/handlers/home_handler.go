package handlers

import (
    "html/template"
    "net/http"
    "log"
    "os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    cwd, err := os.Getwd()  // Получаем текущую рабочую директорию
    if err != nil {
        http.Error(w, "Failed to get current working directory", http.StatusInternalServerError)
        return
    }
    log.Println("Current working directory:", cwd)  // Выводим в лог текущую рабочую директорию

    tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}
