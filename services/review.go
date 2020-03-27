package services

import (
	"log"
	"time"

	"github.com/dinopuguh/kawulo-go-crawler/api"
	"github.com/dinopuguh/kawulo-go-crawler/database"
	"github.com/dinopuguh/kawulo-go-crawler/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllReviews(db *mongo.Database) ([]models.Review, error) {
	ctx := database.Ctx

	crs, err := db.Collection("review").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer crs.Close(ctx)

	result := make([]models.Review, 0)
	for crs.Next(ctx) {
		var row models.Review
		err := crs.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result, nil
}

func ReviewExist(db *mongo.Database, revId string) (bool, error) {
	ctx := database.Ctx

	var result models.Review

	err := db.Collection("review").FindOne(ctx, bson.D{primitive.E{Key: "id", Value: revId}}).Decode(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func InsertReview(db *mongo.Database, restId primitive.ObjectID, rev api.Review) error {
	ctx := database.Ctx

	var newRev models.Review
	var newSubratings []models.Subrating

	for _, s := range rev.Subratings {
		newSubratings = append(newSubratings, models.Subrating(s))
	}

	newRev.ID = primitive.NewObjectID()
	newRev.Id = rev.ReviewId
	newRev.Lang = rev.Lang
	newRev.LocationId = rev.LocationId
	newRev.PublishedDate = rev.PublishedDate
	newRev.Rating = rev.Rating
	newRev.RestaurantID = restId
	newRev.Subratings = newSubratings
	newRev.Text = rev.Text
	newRev.CreatedAt = time.Now()

	crs, err := db.Collection("review").InsertOne(ctx, newRev)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Insert review success -", crs.InsertedID)

	return nil
}

func InsertReviews(db *mongo.Database, rest models.Restaurant) error {
	url := api.LocationUrl + rest.LocationId + "/reviews"

	for {
		log.Println("<--- *** --->")
		res, err := api.FetchReviews(url)
		if err != nil {
			log.Fatal(err.Error())
		}

		revs := res.Data
		pag := res.Paging

		for _, rev := range revs {
			exist, _ := ReviewExist(db, rev.ReviewId)
			if exist == true {
				log.Println("Review with id", rev.ReviewId, "is already exist")
				continue
			}

			err = InsertReview(db, rest.ID, rev)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Println(rev.ReviewId)
		}

		if pag.Next != "" {
			url = pag.Next
		} else {
			log.Println(rest.LocationId, "done")
			break
		}
	}

	return nil
}
