package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	BaseUrl     = "https://api.tripadvisor.com/api/internal/1.14/"
	LocationUrl = BaseUrl + "location/"
	transport   = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	res     *http.Response
	retries int
)

type RestaurantResponse struct {
	Data   []Restaurant `json:"data"`
	Paging Paging       `json:"paging"`
}

type ReviewResponse struct {
	Data   []Review `json:"data"`
	Paging Paging   `json:"paging"`
}

type Paging struct {
	Previous     string `json:"previous"`
	Next         string `json:"next"`
	Skipped      string `json:"skipped"`
	Results      string `json:"results"`
	TotalResults string `json:"total_results"`
}

func FetchRestaurants(url string) (RestaurantResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := &http.Client{
		Transport: transport,
	}

	var data RestaurantResponse
	retries = 3

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("X-TripAdvisor-API-Key", os.Getenv("TRIPADVISOR_API_KEY"))
	req.Close = true

	for retries > 0 {
		res, err = client.Do(req)

		if err != nil {
			retries -= 1
			log.Println("Retrying...")
		} else {
			defer res.Body.Close()
			break
		}
	}

	if err != nil {
		return data, err
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func FetchReviews(url string) (ReviewResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := &http.Client{
		Transport: transport,
	}

	var data ReviewResponse
	retries = 3

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("X-TripAdvisor-API-Key", os.Getenv("TRIPADVISOR_API_KEY"))
	req.Close = true

	for retries > 0 {
		res, err = client.Do(req)

		if err != nil {
			retries -= 1
			log.Println("Retrying...")
		} else {
			defer res.Body.Close()
			break
		}
	}

	if err != nil {
		return data, err
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
