package services

import (
	"context"
	"log"
	"os"

	"github.com/mide7/go_movie_scrapper/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetAllParams struct {
	Page  int64
	Limit int64
}

type Meta struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type IMovie interface {
	Create(ctx context.Context, movie *models.Movie) error
	GetAll(ctx context.Context, params *GetAllParams) ([]models.Movie, *Meta, error)
}

type MovieService struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

func getCollection(mongoClient *mongo.Client, collection string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return mongoClient.Database(dbName).Collection(collection)
}

func NewMovieService(mongoClient *mongo.Client) *MovieService {
	return &MovieService{
		mongoClient: mongoClient,
		collection:  getCollection(mongoClient, "movies"),
	}
}

func (service *MovieService) Create(ctx context.Context, movie *models.Movie) error {
	var existingMovie models.Movie

	result := service.collection.FindOne(ctx, bson.M{"url": movie.URL, "title": movie.Title})
	if result.Err() != nil {
		return result.Err()
	}

	result.Decode(&existingMovie)

	if existingMovie.ID != primitive.NilObjectID {
		return nil
	}

	_, err := service.collection.InsertOne(ctx, movie)
	return err
}

func (m *MovieService) GetAll(ctx context.Context, params *GetAllParams) ([]models.Movie, *Meta, error) {
	var movies []models.Movie
	var total int64

	if params == nil {
		params = &GetAllParams{
			Page:  1,
			Limit: 10,
		}
	}

	findOptions := options.Find().SetSkip((params.Page - 1) * params.Limit).SetLimit(params.Limit)
	cur, err := m.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Printf("failed to get all movies: %v", err)
		return nil, nil, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &movies); err != nil {
		log.Printf("failed to marshal movies: %v", err)
		return nil, nil, err
	}

	total, err = m.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("failed to get movies count: %v", err)
	}

	return movies, &Meta{
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}
