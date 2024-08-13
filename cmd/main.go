package main

import (
    "log"
    "net/http"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    _, err := w.Write([]byte("Welcome to Spectra API"))
    if err != nil {
        http.Error(w, "Unable to write response", http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/", WelcomeHandler)

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
