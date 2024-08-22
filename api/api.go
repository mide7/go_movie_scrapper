package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mide7/go/movie_scrapper/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(ctx context.Context, router *gin.Engine, mongoClient *mongo.Client) *gin.Engine {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	movieHandler := newMovieApiHandler(ctx, services.NewMovieService(mongoClient))

	router.GET("/movies", movieHandler.GetAllHandler)

	corsSetup(router)

	return router
}

func corsSetup(router *gin.Engine) {
	// allow swagger UI requests
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           time.Hour * 12,
		AllowWildcard:    true,
	}))
}
