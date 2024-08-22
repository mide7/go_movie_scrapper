package scrapper

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mide7/go_movie_scrapper/internal/models"
	"github.com/mide7/go_movie_scrapper/internal/services"
	"github.com/mide7/go_movie_scrapper/internal/sites"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartScrapper(concurrency int, timeBetweenRequest time.Duration, mongoClient *mongo.Client) {
	ticker := time.NewTicker(timeBetweenRequest)

	sites := []string{"imdb", "netflix", "prime"}

	for ; ; <-ticker.C {
		var wg sync.WaitGroup

		for _, site := range sites {

			if site == "imdb" {
				wg.Add(1)
				go ImdbScrapper(&wg, mongoClient)
			}
		}

		wg.Wait()
	}

}

func ImdbScrapper(wg *sync.WaitGroup, mongoClient *mongo.Client) {
	movies, err := sites.ImdbGetUpcomingMovies()
	if err != nil {
		log.Printf("IMDB: failed to get upcoming movies: %v", err)
	}

	fmt.Print("IMDB: ", len(*movies), " movies found\n")

	for _, movieData := range *movies {
		movie := models.Movie{Title: movieData.Title, Year: movieData.Year, URL: movieData.URL, Image: movieData.Image, Genres: movieData.Genres}
		services := services.NewMovieService(mongoClient)
		err = services.Create(context.Background(), &movie)
		if err != nil {
			log.Printf("IMDB: failed to save movie: %v", err)
		}
	}

	wg.Done()
}
