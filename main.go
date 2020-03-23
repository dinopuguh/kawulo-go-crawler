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
}

func getRestaurants(db *mongo.Database) {
	locs, err := services.FindIndonesianLocations(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, loc := range locs {
		log.Println("<--- Location ", loc.Name, "--->")
		err = services.InsertRestaurants(db, loc)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func getReviews(db *mongo.Database) {
	locs, err := services.FindIndonesianLocations(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, loc := range locs {
		log.Println("<--- Location ", loc.Name, "--->")

		rests, err := services.FindIndonesianRestaurants(db, loc.LocationId)
		if err != nil {
			log.Fatal(err)
		}

		for _, rest := range rests {
			log.Println("Restaurant ", rest.Name)
			err = services.InsertReviews(db, rest)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
