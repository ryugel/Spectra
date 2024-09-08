package scraper

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"spectra/internal/models"
	"spectra/internal/utils"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func ConfigureCollector() *colly.Collector {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		logrus.Infof("Visiting: %s", r.URL.String())
	})
	return c
}

func FetchReleases() ([]models.AnimeRelease, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	homePageURL := utils.GetEnv("HOME_PAGE_URL", "https://anitaku.pe/home.html")

	var releases []models.AnimeRelease
	c := ConfigureCollector()

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		releases = append(releases, models.AnimeRelease{
			Title:   e.ChildAttr("p.name a", "title"),
			Link:    baseURL + e.ChildAttr("a", "href"),
			Episode: e.ChildText("p.episode"),
			Image:   e.ChildAttr("img", "src"),
		})
	})

	if err := c.Visit(homePageURL); err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}

	return releases, nil
}


func FetchAnimeInfo(title string) (*models.AnimeInfo, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	slug := strings.ReplaceAll(strings.ToLower(title), " ", "-")
	animeURL := fmt.Sprintf("%s/category/%s", baseURL, slug)

	var animeInfo models.AnimeInfo
	c := ConfigureCollector()

	c.OnHTML("div.anime_info_body_bg", func(e *colly.HTMLElement) {
		animeInfo.Title = e.ChildText("h1")
		animeInfo.Image = e.ChildAttr("img", "src")
	})

	c.OnHTML("div.anime_info_body", func(e *colly.HTMLElement) {
		if desc := e.ChildText("div.description"); desc != "" {
			animeInfo.Description = desc
		}

		e.ForEach("p.type:contains('Genre:') a", func(_ int, el *colly.HTMLElement) {
			genre := strings.TrimSpace(el.Text)
			if genre != "" && genre != "," {
				animeInfo.Genres = append(animeInfo.Genres, genre)
			}
		})

		var cleanedGenres []string
		for _, genre := range animeInfo.Genres {
			trimmedGenre := strings.Trim(genre, ", ")
			if trimmedGenre != "" {
				cleanedGenres = append(cleanedGenres, trimmedGenre)
			}
		}
		animeInfo.Genres = cleanedGenres
	})

c.OnHTML("div.anime_video_body", func(e *colly.HTMLElement) {
		e.ForEach("ul#episode_page li a", func(_ int, el *colly.HTMLElement) {
			epStart, err1 := strconv.Atoi(el.Attr("ep_start"))
			epEnd, err2 := strconv.Atoi(el.Attr("ep_end"))

			if err1 != nil || err2 != nil {
				return
			}

			for i := epStart; i <= epEnd; i++ {
				if i == 0 {
					continue
				}
				episodeURL := fmt.Sprintf("%s-episode-%d", slug, i)
				episode := models.Episode{
					Number: fmt.Sprintf("%d", i),
					Url:    baseURL + "/" + episodeURL,
				}
				animeInfo.Episodes = append(animeInfo.Episodes, episode)
			}
		})
  })

	if err := c.Visit(animeURL); err != nil {
		return nil, fmt.Errorf("failed to fetch anime info: %w", err)
	}

	if animeInfo.Title == "" {
		return nil, errors.New("anime not found")
	}

	return &animeInfo, nil
}

func FetchMovies() ([]models.Movie, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	moviesPageURL := fmt.Sprintf("%s/anime-movies.html", baseURL)

	var movies []models.Movie
	c := ConfigureCollector()

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("p.name a", "title")
		isSubbed := "Sub"
		if strings.Contains(title, "Dub") {
			isSubbed = "Dub"
		}

		movies = append(movies, models.Movie{
			Title:       title,
			Link:        baseURL + e.ChildAttr("a", "href"),
			Image:       e.ChildAttr("img", "src"),
			ReleaseDate: e.ChildText("p.released"),
			IsSubbed:    isSubbed,
		})
	})

	if err := c.Visit(moviesPageURL); err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	return movies, nil
}

func SearchQuery(keyword string, genres []string) ([]models.Movie, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	searchURL := fmt.Sprintf("%s/filter.html?keyword=%s", baseURL, url.QueryEscape(keyword))

	if len(genres) > 0 {
		searchURL += "&genre[]=" + strings.Join(genres, "&genre[]=")
	}

	var movies []models.Movie
	c := ConfigureCollector()

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("p.name a", "title")
		isSubbed := "Sub"
		if strings.Contains(title, "Dub") {
			isSubbed = "Dub"
		}

		movies = append(movies, models.Movie{
			Title:       title,
			Link:        baseURL + e.ChildAttr("a", "href"),
			Image:       e.ChildAttr("img", "src"),
			ReleaseDate: e.ChildText("p.released"),
			IsSubbed:    isSubbed,
		})
	})

	if err := c.Visit(searchURL); err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	return movies, nil
}

func FetchPopularAnime() ([]models.Anime, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	popularPageURL := fmt.Sprintf("%s/popular.html", baseURL)

	var popularAnime []models.Anime
	c := ConfigureCollector()

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		popularAnime = append(popularAnime, models.Anime{
			Title:       e.ChildAttr("p.name a", "title"),
			Link:        baseURL + e.ChildAttr("a", "href"),
			Image:       e.ChildAttr("img", "src"),
			ReleaseDate: e.ChildText("p.released"),
		})
	})

	if err := c.Visit(popularPageURL); err != nil {
		return nil, fmt.Errorf("failed to fetch popular anime: %w", err)
	}

	return popularAnime, nil
}

func FetchEpisodeVideoURL(episodeURL string) (string, error) {
	var videoURL string
	c := ConfigureCollector()

	c.OnHTML("div.play-video iframe", func(e *colly.HTMLElement) {
		videoURL = e.Attr("src")
	})

	if err := c.Visit(episodeURL); err != nil {
		return "", fmt.Errorf("failed to fetch video URL: %w", err)
	}

	if videoURL == "" {
		return "", errors.New("video URL not found")
	}

	return videoURL, nil
}


func extractTextAfter(e *colly.HTMLElement, selector, prefix string) string {
	text := e.ChildText(selector)
	return strings.TrimSpace(strings.Replace(text, prefix, "", 1))
}

