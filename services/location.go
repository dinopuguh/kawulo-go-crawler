package services

import (
	"log"

	"github.com/dinopuguh/tripadvisor-crawler/database"
	"github.com/dinopuguh/tripadvisor-crawler/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindAllLocations() ([]models.Location, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	csr, err := db.Collection("location").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer csr.Close(ctx)

	result := make([]models.Location, 0)
	for csr.Next(ctx) {
		var row models.Location
		err := csr.Decode(&row)
		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func FindIndonesianLocations() ([]models.Location, error) {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	cities := []string{
		// "Banda Aceh",
		// "Medan",
		// "Padang",
		// "Pekanbaru",
		// "Palembang",
		// "Bengkulu",
		// "Bandar Lampung",
		// "Pangkal Pinang",
		// "Tanjung Pinang",
		"Jakarta",
		"Bandung",
		"Semarang",
		"Yogyakarta Region",
		"Surabaya",
		"Serang",
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
			return nil, err
		}
		defer csr.Close(ctx)

		for csr.Next(ctx) {
			var row models.Location
			err := csr.Decode(&row)
			if err != nil {
				return nil, err
			}

			result = append(result, row)
		}
	}

	return result, nil
}

func InsertLocation(loc models.Location) error {
	ctx := database.Ctx
	db, err := database.Connect()
	if err != nil {
		return err
	}

	loc.ID = primitive.NewObjectID()
	crs, err := db.Collection("location").InsertOne(ctx, loc)
	if err != nil {
		return err
	}

	log.Println("Insert location success -", crs.InsertedID)
	return nil
}
