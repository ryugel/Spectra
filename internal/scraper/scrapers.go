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
)

func FetchReleases() ([]models.AnimeRelease, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	homePageURL := utils.GetEnv("HOME_PAGE_URL", "https://anitaku.pe/home.html")

	var releases []models.AnimeRelease

	c := colly.NewCollector()

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("p.name a", "title")
		link := baseURL + e.ChildAttr("a", "href")
		episode := e.ChildText("p.episode")
		image := e.ChildAttr("img", "src")

		releases = append(releases, models.AnimeRelease{
			Title:   title,
			Link:    link,
			Episode: episode,
			Image:   image,
		})
	})

	err := c.Visit(homePageURL)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func FetchAnimeInfo(title string) (*models.AnimeInfo, error) {
	slug := strings.ReplaceAll(title, " ", "-")
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	url := baseURL + "/category/" + slug

	c := colly.NewCollector()

	var animeInfo models.AnimeInfo

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	c.OnHTML("div.anime_info_body_bg", func(e *colly.HTMLElement) {
		animeInfo.Title = e.ChildText("h1")
		animeInfo.Image = e.ChildAttr("img", "src")
	})

	c.OnHTML("div.anime_info_body", func(e *colly.HTMLElement) {
		if desc := e.ChildText("div.description"); desc != "" {
			animeInfo.Description = desc
		}

		if release := e.ChildText("p.type:contains('Released:')"); release != "" {
			animeInfo.ReleaseDate = strings.TrimSpace(strings.Replace(release, "Released:", "", 1))
		}

		e.ForEach("p.type:contains('Genre:') a", func(_ int, el *colly.HTMLElement) {
			genre := strings.TrimSpace(el.Text)
			animeInfo.Genres = append(animeInfo.Genres, genre)
		})

		if status := e.ChildText("p.type:contains('Status:') a"); status != "" {
			animeInfo.Status = strings.TrimSpace(status)
		}
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

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	if animeInfo.Title == "" {
		return nil, errors.New("anime not found")
	}

	return &animeInfo, nil
}

func FetchMovies() ([]models.Movie, error) {
	var movies []models.Movie

	c := colly.NewCollector()

	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("p.name a", "title")
		link := baseURL + e.ChildAttr("a", "href")
		image := e.ChildAttr("img", "src")
		releaseDate := e.ChildText("p.released")

		movies = append(movies, models.Movie{
			Title:       title,
			Link:        link,
			Image:       image,
			ReleaseDate: releaseDate,
		})
	})

	err := c.Visit(baseURL + "/anime-movies.html")
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func SearchQuery(keyword string, genres []string) ([]models.Movie, error) {
	var movies []models.Movie

	c := colly.NewCollector()

	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")

	searchURL := baseURL + "/filter.html?keyword=" + url.QueryEscape(keyword)

	if len(genres) > 0 {
		searchURL += "&genre[]=" + strings.Join(genres, "&genre[]=")
	}

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("p.name a", "title")
		link := baseURL + e.ChildAttr("a", "href")
		image := e.ChildAttr("img", "src")
		releaseDate := e.ChildText("p.released")

		movies = append(movies, models.Movie{
			Title:       title,
			Link:        link,
			Image:       image,
			ReleaseDate: releaseDate,
		})
	})

	err := c.Visit(searchURL)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func FetchPopularAnime() ([]models.Anime, error) {
	baseURL := utils.GetEnv("BASE_URL", "https://anitaku.pe")
	popularPageURL := baseURL + "/popular.html"

	var popularAnime []models.Anime

	c := colly.NewCollector()

	c.OnHTML(".last_episodes .items li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("p.name a", "title")
		link := baseURL + e.ChildAttr("a", "href")
		image := e.ChildAttr("img", "src")
		releaseDate := e.ChildText("p.released")

		popularAnime = append(popularAnime, models.Anime{
			Title:       title,
			Link:        link,
			Image:       image,
			ReleaseDate: releaseDate,
		})
	})

	err := c.Visit(popularPageURL)
	if err != nil {
		return nil, err
	}

	return popularAnime, nil
}

func FetchEpisodeVideoURL(episodeURL string) (string, error) {
	var videoURL string

	c := colly.NewCollector()

	c.OnHTML("div.play-video iframe", func(e *colly.HTMLElement) {
		videoURL = e.Attr("src")
	})

	err := c.Visit(episodeURL)
	if err != nil {
		return "", err
	}

	if videoURL == "" {
		return "", errors.New("video URL not found")
	}

	return videoURL, nil
}
