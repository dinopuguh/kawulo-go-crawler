package services

import (
	"log"
	"time"

	"github.com/dinopuguh/tripadvisor-crawler/api"
	"github.com/dinopuguh/tripadvisor-crawler/database"
	"github.com/dinopuguh/tripadvisor-crawler/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindAllReviews() ([]models.Review, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	crs, err := db.Collection("review").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer crs.Close(ctx)

	result := make([]models.Review, 0)
	for crs.Next(ctx) {
		var row models.Review
		err := crs.Decode(&row)
		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func ReviewExist(rev_id string) (bool, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return false, err
	}
	var result models.Review

	err = db.Collection("review").FindOne(ctx, bson.D{primitive.E{Key: "id", Value: rev_id}}).Decode(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func InsertReview(rest_id primitive.ObjectID, rev api.Review) error {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var newRev models.Review
	var newSubratings []models.Subrating

	for _, s := range rev.Subratings {
		var newSubrating models.Subrating

		newSubrating.Name = s.Name
		newSubrating.Value = s.Value

		newSubratings = append(newSubratings, newSubrating)
	}

	newRev.ID = primitive.NewObjectID()
	newRev.Id = rev.ReviewId
	newRev.Lang = rev.Lang
	newRev.LocationId = rev.LocationId
	newRev.PublishedDate = rev.PublishedDate
	newRev.Rating = rev.Rating
	newRev.RestaurantID = rest_id
	newRev.Subratings = newSubratings
	newRev.Text = rev.Text
	newRev.CreatedAt = time.Now()

	crs, err := db.Collection("review").InsertOne(ctx, newRev)
	if err != nil {
		return err
	}

	log.Println("Insert review success -", crs.InsertedID)
	return nil
}

func InsertReviews(rest models.Restaurant) error {
	url := api.LocationUrl + rest.LocationId + "/reviews"

	for {
		log.Println("<--- *** --->")
		res, err := api.FetchReviews(url)
		if err != nil {
			return err
		}

		revs := res.Data
		pag := res.Paging

		for _, rev := range revs {
			exist, _ := ReviewExist(rev.ReviewId)
			if exist == true {
				log.Println("Review with id", rev.ReviewId, "is already exist")
				continue
			}

			err = InsertReview(rest.ID, rev)
			if err != nil {
				return err
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
