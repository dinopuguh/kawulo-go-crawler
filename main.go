package main

import (
	"flag"
	"log"

	"github.com/dinopuguh/kawulo-go-crawler/database"
	"github.com/dinopuguh/kawulo-go-crawler/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	var name string
	flag.StringVar(&name, "data", "review", "Usage")

	flag.Parse()

	switch name {
	case "restaurant":
		getRestaurants(db)
		break
	case "review":
		getReviews(db)
		break
	default:
		getReviews(db)
		break
	}

	err = db.Client().Disconnect(database.Ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getRestaurants(db *mongo.Database) {
	locs := services.FindIndonesianLocations(db)

	for _, loc := range locs {
		log.Println("<--- Location ", loc.Name, "--->")
		services.InsertRestaurants(db, loc)
	}
}

func getReviews(db *mongo.Database) {
	locs := services.FindIndonesianLocations(db)

	for _, loc := range locs {
		log.Println("<--- Location ", loc.Name, "--->")

		rests := services.FindRestaurantByLocId(db, loc.LocationId)

		for _, rest := range rests {
			log.Println("Restaurant ", rest.Name)
			services.InsertReviews(db, rest)
		}
	}
}
