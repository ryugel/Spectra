package handlers

import (
	"encoding/json"
	"net/http"
	"spectra/internal/scraper"
)

func GetReleasesHandler(w http.ResponseWriter, r *http.Request) {
	releases, err := scraper.FetchReleases()
	if err != nil {
		http.Error(w, "Failed to fetch releases", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(releases); err != nil {
		http.Error(w, "Failed to encode releases", http.StatusInternalServerError)
	}
}

func GetAnimeInfoHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing title parameter", http.StatusBadRequest)
		return
	}

	animeInfo, err := scraper.FetchAnimeInfo(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animeInfo); err != nil {
		http.Error(w, "Failed to encode anime info", http.StatusInternalServerError)
	}
}
