package services

import (
	"log"

	"github.com/dinopuguh/kawulo-go-crawler/database"
	"github.com/dinopuguh/kawulo-go-crawler/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllLocations(db *mongo.Database) ([]models.Location, error) {
	ctx := database.Ctx

	csr, err := db.Collection("location").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]models.Location, 0)
	for csr.Next(ctx) {
		var row models.Location
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	return result, nil
}

func FindIndonesianLocations(db *mongo.Database) ([]models.Location, error) {
	ctx := database.Ctx

	cities := []string{
		"Surabaya",
		"Banda Aceh",
		"Medan",
		"Padang",
		"Pekanbaru",
		"Palembang",
		"Bengkulu",
		"Bandar Lampung",
		"Pangkal Pinang",
		"Tanjung Pinang",
		"Serang",
		"Jakarta",
		"Bandung",
		"Semarang",
		"Yogyakarta Region",
		"Denpasar",
		"Mataram",
		"Kupang",
		"Pontianak",
		"Banjarmasin",
		"Samarinda",
		"Manado",
		"Palu",
		"Makassar",
		"Kendari",
		"Gorontalo",
		"Mamuju",
		"Ambon",
		"Jayapura",
		"Manokwari",
	}

	result := make([]models.Location, 0)

	for _, city := range cities {
		csr, err := db.Collection("location").Find(ctx, bson.M{"name": city})
		if err != nil {
			log.Fatal(err.Error())
		}
		defer csr.Close(ctx)

		for csr.Next(ctx) {
			var row models.Location
			err := csr.Decode(&row)
			if err != nil {
				log.Fatal(err.Error())
			}

			result = append(result, row)
		}
	}

	return result, nil
}

func InsertLocation(db *mongo.Database, loc models.Location) error {
	ctx := database.Ctx

	loc.ID = primitive.NewObjectID()
	crs, err := db.Collection("location").InsertOne(ctx, loc)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Insert location success -", crs.InsertedID)

	return nil
}
