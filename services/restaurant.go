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

func FindAllRestaurants(db *mongo.Database) ([]models.Restaurant, error) {
	ctx := database.Ctx

	csr, err := db.Collection("restaurant").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Restaurant, 0)
	for csr.Next(ctx) {
		var row models.Restaurant
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result, nil
}

func FindIndonesianRestaurants(db *mongo.Database, locId string) ([]models.Restaurant, error) {
	ctx := database.Ctx

	csr, err := db.Collection("restaurant").Find(ctx, bson.M{"locationID": locId})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Restaurant, 0)
	for csr.Next(ctx) {
		var row models.Restaurant
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result, nil
}

func RestaurantExist(db *mongo.Database, locId string) (bool, error) {
	ctx := database.Ctx

	var result models.Restaurant

	err := db.Collection("restaurant").FindOne(ctx, bson.D{primitive.E{Key: "location_id", Value: locId}}).Decode(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func InsertRestaurant(db *mongo.Database, loc models.Location, rest api.Restaurant) error {
	ctx := database.Ctx

	var newRest models.Restaurant
	var newSubCategory []models.SubCategory

	for _, s := range rest.SubCategory {
		newSubCategory = append(newSubCategory, models.SubCategory(s))
	}

	newRest.ID = primitive.NewObjectID()
	newRest.LocationId = rest.LocationId
	newRest.Name = rest.Name
	newRest.Latitude = rest.Latitude
	newRest.Longitude = rest.Longitude
	newRest.NumReviews = rest.NumReviews
	newRest.Photo.Images.Thumbnail = models.Image(rest.Photo.Images.Thumbnail)
	newRest.Photo.Images.Original = models.Image(rest.Photo.Images.Original)
	newRest.Photo.Images.Medium = models.Image(rest.Photo.Images.Medium)
	newRest.Photo.Images.Large = models.Image(rest.Photo.Images.Large)
	newRest.Rating = rest.Rating
	newRest.Price = rest.Price
	newRest.Address = rest.Address
	newRest.Phone = rest.Phone
	newRest.Website = rest.Website
	newRest.RawRanking = rest.RawRanking
	newRest.RankingGeo = rest.RankingGeo
	newRest.RankingPosition = rest.RankingPosition
	newRest.RankingDenominator = rest.RankingDenominator
	newRest.RankingCategory = rest.RankingCategory
	newRest.Ranking = rest.Ranking
	newRest.PriceLevel = rest.PriceLevel
	newRest.SubCategory = newSubCategory
	newRest.LocationID = loc.LocationId
	newRest.LocationObjectID = loc.ID
	newRest.CreatedAt = time.Now()

	crs, err := db.Collection("restaurant").InsertOne(ctx, newRest)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Insert restaurant success -", crs.InsertedID)

	return nil
}

func InsertRestaurants(db *mongo.Database, loc models.Location) error {
	url := api.LocationUrl + loc.LocationId + "/restaurants"

	for {
		log.Println("<--- *** --->")
		res, err := api.FetchRestaurants(url)
		if err != nil {
			log.Fatal(err.Error())
		}

		rests := res.Data
		pag := res.Paging

		for _, rest := range rests {
			exist, _ := RestaurantExist(db, rest.LocationId)
			if exist == true {
				log.Println("Restaurant with id", rest.LocationId, "is already exist")
				continue
			}

			err = InsertRestaurant(db, loc, rest)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		if pag.Next != "" {
			url = pag.Next
		} else {
			log.Println(loc.LocationId, "done")
			break
		}
	}

	return nil
}
