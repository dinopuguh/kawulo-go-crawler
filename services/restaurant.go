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

func FindAllRestaurants(db *mongo.Database) []models.Restaurant {
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

	return result
}

func FindIndonesianRestaurants(db *mongo.Database, locId string) []models.Restaurant {
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

	return result
}

func RestaurantExist(db *mongo.Database, locId string) bool {
	ctx := database.Ctx

	var result models.Restaurant

	err := db.Collection("restaurant").FindOne(ctx, bson.D{primitive.E{Key: "location_id", Value: locId}}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}

func InsertRestaurant(db *mongo.Database, loc models.Location, rest api.Restaurant) {
	ctx := database.Ctx

	var newSubCategory []models.SubCategory

	for _, s := range rest.SubCategory {
		newSubCategory = append(newSubCategory, models.SubCategory(s))
	}

	newRest := models.Restaurant{
		ID:         primitive.NewObjectID(),
		LocationId: rest.LocationId,
		Name:       rest.Name,
		Latitude:   rest.Latitude,
		Longitude:  rest.Longitude,
		NumReviews: rest.NumReviews,
		Photo: models.Photo{
			Images: models.Images{
				Thumbnail: models.Image(rest.Photo.Images.Thumbnail),
				Original:  models.Image(rest.Photo.Images.Original),
				Medium:    models.Image(rest.Photo.Images.Medium),
				Large:     models.Image(rest.Photo.Images.Thumbnail),
			},
		},
		Rating:             rest.Rating,
		Price:              rest.Price,
		Address:            rest.Address,
		Phone:              rest.Phone,
		Website:            rest.Website,
		RawRanking:         rest.RawRanking,
		RankingGeo:         rest.RankingGeo,
		RankingPosition:    rest.RankingPosition,
		RankingDenominator: rest.RankingDenominator,
		RankingCategory:    rest.RankingCategory,
		Ranking:            rest.Ranking,
		PriceLevel:         rest.PriceLevel,
		SubCategory:        newSubCategory,
		LocationID:         loc.LocationId,
		LocationObjectID:   loc.ID,
		CreatedAt:          time.Now(),
	}

	crs, err := db.Collection("restaurant").InsertOne(ctx, newRest)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Insert restaurant success -", crs.InsertedID)
}

func InsertRestaurants(db *mongo.Database, loc models.Location) {
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
			exist := RestaurantExist(db, rest.LocationId)
			if exist == true {
				log.Println("Restaurant with id", rest.LocationId, "is already exist")
				continue
			}

			InsertRestaurant(db, loc, rest)
		}

		if pag.Next != "" {
			url = pag.Next
		} else {
			log.Println(loc.LocationId, "done")
			break
		}
	}
}
