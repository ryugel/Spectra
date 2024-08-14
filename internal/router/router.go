package router

import (
	"net/http"
	"spectra/internal/handlers"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", WelcomeHandler)
	mux.HandleFunc("/api/releases", handlers.GetReleasesHandler)
	mux.HandleFunc("/api/anime", handlers.GetAnimeInfoHandler)
	return mux
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Welcome to Spectra API"))
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
	}
}
