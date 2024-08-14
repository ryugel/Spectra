package models

type AnimeRelease struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Episode string `json:"episode"`
	Image   string `json:"image"`
}

type AnimeInfo struct {
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Genres      []string  `json:"genres"`
	ReleaseDate string    `json:"release_date"`
	Status      string    `json:"status"`
	Episodes    []Episode `json:"episodes"`
}

type Episode struct {
	Number string `json:"number"`
	Url    string `json:"url"`
}
