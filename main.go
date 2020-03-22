package main

import (
	"flag"
	"log"

	"github.com/dinopuguh/tripadvisor-crawler/services"
)

func main() {
	var name string
	flag.StringVar(&name, "data", "review", "Usage")

	flag.Parse()

	switch name {
	case "restaurant":
		getRestaurants()
		break
	case "review":
		getReviews()
		break
	default:
		getReviews()
		break
	}
}

func getRestaurants() {
	locs, err := services.FindIndonesianLocations()
	if err != nil {
		log.Fatal(err)
	}

	for _, loc := range locs {
		log.Println("<--- Location ", loc.Name, "--->")
		err = services.InsertRestaurants(loc.LocationId)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func getReviews() {
	locs, err := services.FindIndonesianLocations()
	if err != nil {
		log.Fatal(err)
	}

	for _, loc := range locs {
		log.Println("<--- Location ", loc.Name, "--->")

		rests, err := services.FindIndonesianRestaurants(loc.LocationId)
		if err != nil {
			log.Fatal(err)
		}

		for _, rest := range rests {
			log.Println("Restaurant ", rest.Name)
			err = services.InsertReviews(rest)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}
