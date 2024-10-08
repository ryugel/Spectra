package handlers

import (
	"encoding/json"
	"net/http"
	"spectra/internal/scraper"
	"strings"
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

func GetMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := scraper.FetchMovies()
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
	}
}

func SearchQueriesHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	genres := r.URL.Query()["genre[]"]

	movies, err := scraper.SearchQuery(keyword, genres)
	if err != nil {
		http.Error(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		http.Error(w, "Failed to encode movies", http.StatusInternalServerError)
	}
}

func GetPopularAnimeHandler(w http.ResponseWriter, r *http.Request) {
	popularAnime, err := scraper.FetchPopularAnime()
	if err != nil {
		http.Error(w, "Failed to fetch popular anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(popularAnime); err != nil {
		http.Error(w, "Failed to encode popular anime", http.StatusInternalServerError)
	}
}
func GetEpisodeVideoURLHandler(w http.ResponseWriter, r *http.Request) {
	episodeURL := r.URL.Query().Get("episode_url")
	if episodeURL == "" {
		http.Error(w, "Missing episode_url parameter", http.StatusBadRequest)
		return
	}

	videoURL, err := scraper.FetchEpisodeVideoURL(episodeURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"video_url": videoURL}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode video URL", http.StatusInternalServerError)
		return
	}

	jsonString := strings.ReplaceAll(string(jsonData), `\u0026`, "&")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonString))
}
