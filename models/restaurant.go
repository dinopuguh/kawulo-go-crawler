package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	ID         primitive.ObjectID `bson:"_id"`
	LocationId string             `bson:"location_id"`
	Name       string             `bson:"name"`
	Latitude   string             `bson:"latitude"`
	Longitude  string             `bson:"longitude"`
	NumReviews string             `bson:"num_reviews"`
	LocationID string             `bson:"locationID"`
	CreatedAt  time.Time          `bson:"created_at"`
}
