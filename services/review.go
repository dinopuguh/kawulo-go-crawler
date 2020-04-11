package services

import (
	"log"
	"strconv"
	"time"

	"github.com/dinopuguh/kawulo-go-crawler/api"
	"github.com/dinopuguh/kawulo-go-crawler/database"
	"github.com/dinopuguh/kawulo-go-crawler/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllReviews(db *mongo.Database) []models.Review {
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

	return result
}

func reviewExist(db *mongo.Database, revId string) bool {
	ctx := database.Ctx

	var result models.Review

	err := db.Collection("review").FindOne(ctx, bson.D{{Key: "id", Value: revId}}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}

func countReviews(db *mongo.Database, restId string) int64 {
	ctx := database.Ctx

	count, err := db.Collection("review").CountDocuments(ctx, bson.M{"location_id": restId})
	if err != nil {
		log.Panic(err.Error())
	}

	return count
}

func InsertReview(db *mongo.Database, restId primitive.ObjectID, rev api.Review) {
	ctx := database.Ctx

	var newSubratings []models.Subrating

	for _, s := range rev.Subratings {
		newSubratings = append(newSubratings, models.Subrating(s))
	}

	newRev := models.Review{
		ID:            primitive.NewObjectID(),
		Id:            rev.ReviewId,
		Lang:          rev.Lang,
		LocationId:    rev.LocationId,
		PublishedDate: rev.PublishedDate,
		Rating:        rev.Rating,
		RestaurantID:  restId,
		Subratings:    newSubratings,
		Text:          rev.Text,
		CreatedAt:     time.Now(),
	}

	crs, err := db.Collection("review").InsertOne(ctx, newRev)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Insert review success -", crs.InsertedID)
}

func InsertReviews(db *mongo.Database, rest models.Restaurant) {
	reviewCount := countReviews(db, rest.LocationId)

	url := api.LocationUrl + rest.LocationId + "/reviews?offset=" + strconv.FormatInt(reviewCount, 10)

	for {
		log.Println("<--- *** --->")
		res, err := api.FetchReviews(url)
		if err != nil {
			log.Fatal(err.Error())
		}

		revs := res.Data
		pag := res.Paging
		results, err := strconv.Atoi(pag.Results)
		if err != nil {
			log.Panic(err.Error())
		}

		if results == 0 {
			log.Println("all reviews in", rest.LocationId, "is already obtained", reviewCount)
			break
		}

		for _, rev := range revs {
			InsertReview(db, rest.ID, rev)
			log.Println(rev.ReviewId)
		}

		if pag.Next != "" {
			url = pag.Next
		} else {
			log.Println(rest.LocationId, "done")
			break
		}
	}

}
