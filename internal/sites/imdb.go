package sites

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mide7/go/movie_scrapper/internal/models"
)

func ImdbGetUpcomingMovies() (*[]models.Movie, error) {
	movies, client := []models.Movie{}, &http.Client{}

	req, err := http.NewRequest("GET", "https://www.imdb.com/calendar/?region=US&type=MOVIE", nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.101.76 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to get: %v", err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("failed to get doc: %v", err)
		return nil, err
	}

	doc.Find(`[data-testid="coming-soon-entry"]`).Each(func(i int, s *goquery.Selection) {
		movie := models.Movie{}
		listItem := s.Find(".ipc-metadata-list-summary-item__t")

		movie.Title = strings.ToLower(strings.Trim(strings.Split(listItem.Text(), "(")[0], " "))

		movie.Year = strings.Trim(strings.Split(strings.Split(listItem.Text(), ")")[0], "(")[1], " ")

		href, ok := listItem.Attr("href")
		if ok {
			movie.URL = "https://www.imdb.com" + href
		}

		genreList := s.Find("ul.ipc-metadata-list-summary-item__tl")
		genreList.Find("li").Each(func(i int, s *goquery.Selection) {
			movie.Genres = append(movie.Genres, strings.Trim(s.Text(), " "))
		})

		listImage := s.Find("div>div>img.ipc-image")
		movieImg, ok := listImage.Attr("src")
		if ok {
			movie.Image = movieImg
		}

		movies = append(movies, movie)

	})

	return &movies, nil
}
