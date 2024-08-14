package models

type AnimeRelease struct {
	Title   string
	Link    string
	Episode string
	Image   string
}

type AnimeInfo struct {
	Title       string
	Image       string
	Description string
	Genres      []string
	ReleaseDate string
	Status      string
	Episodes    []Episode
}

type Episode struct {
	Number string
	Url    string
}

type Movie struct {
	Title       string
	Link        string
	Image       string
	ReleaseDate string
}

type Anime struct {
	Title       string
	Link        string
	Image       string
	ReleaseDate string
}
