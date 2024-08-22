package models

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Movie struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title  string             `bson:"title" json:"title"`
	URL    string             `bson:"url" json:"url"`
	Year   string             `bson:"year" json:"year"`
	Info   string             `bson:"info" json:"info"`
	Genres []string           `bson:"genres" json:"genres"`
	Image  string             `bson:"image" json:"image"`
}

func (m Movie) Save(client *mongo.Client) error {

	dbName, ctx := os.Getenv("DB_NAME"), context.Background()
	collection := client.Database(dbName).Collection("movies")

	var existingMovie Movie

	result := collection.FindOne(ctx, bson.M{"url": m.URL, "title": m.Title})
	if result.Err() != nil {
		return result.Err()
	}

	result.Decode(&existingMovie)

	if existingMovie.ID != primitive.NilObjectID {
		return nil
	}

	_, err := collection.InsertOne(ctx, m)
	return err
}
