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

func FindAllRestaurants() ([]models.Restaurant, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	csr, err := db.Collection("restaurant").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer csr.Close(ctx)

	result := make([]models.Restaurant, 0)
	for csr.Next(ctx) {
		var row models.Restaurant
		err := csr.Decode(&row)
		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func FindIndonesianRestaurants(loc_id string) ([]models.Restaurant, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	csr, err := db.Collection("restaurant").Find(ctx, bson.M{"locationID": loc_id})
	if err != nil {
		return nil, err
	}
	defer csr.Close(ctx)

	result := make([]models.Restaurant, 0)
	for csr.Next(ctx) {
		var row models.Restaurant
		err := csr.Decode(&row)
		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func RestaurantExist(loc_id string) (bool, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return false, err
	}
	var result models.Restaurant

	err = db.Collection("restaurant").FindOne(ctx, bson.D{primitive.E{Key: "location_id", Value: loc_id}}).Decode(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func InsertRestaurant(rest api.Restaurant) error {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var newRest models.Restaurant

	newRest.ID = primitive.NewObjectID()
	newRest.LocationId = rest.LocationId
	newRest.Name = rest.Name
	newRest.Latitude = rest.Latitude
	newRest.Longitude = rest.Longitude
	newRest.NumReviews = rest.NumReviews
	newRest.LocationID = rest.LocationID
	newRest.CreatedAt = time.Now()

	crs, err := db.Collection("restaurant").InsertOne(ctx, newRest)
	if err != nil {
		return err
	}

	log.Println("Insert restaurant success -", crs.InsertedID)
	return nil
}

func InsertRestaurants(loc_id string) error {
	url := api.LocationUrl + loc_id + "/restaurants"

	for {
		log.Println("<--- *** --->")
		res, err := api.FetchRestaurants(url)
		if err != nil {
			return err
		}

		rests := res.Data
		pag := res.Paging

		for _, rest := range rests {
			exist, _ := RestaurantExist(rest.LocationId)
			if exist == true {
				log.Println("Restaurant with id", rest.LocationId, "is already exist")
				continue
			}

			rest.LocationID = loc_id

			err = InsertRestaurant(rest)
			if err != nil {
				return err
			}
		}

		if pag.Next != "" {
			url = pag.Next
		} else {
			log.Println(loc_id, "done")
			break
		}
	}

	return nil
}
