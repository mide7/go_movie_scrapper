package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mide7/go_movie_scrapper/api"
	"github.com/mide7/go_movie_scrapper/internal/mongodb"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	ctx := context.Background()
	godotenv.Load("./config/.env")
	mongoUri := os.Getenv("MONGO_URI")
	mongoClient, err := mongodb.Connect(ctx, mongoUri)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer mongoClient.Disconnect(ctx)

	// go scrapper.StartScrapper(1, time.Hour*24, mongoClient)

	r := gin.Default()
	r = api.Setup(ctx, r, mongoClient)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
