package api

import (
	"encoding/json"
	"net"
	"net/http"
	"time"
)

var (
	BaseUrl     = "https://api.tripadvisor.com/api/internal/1.14/"
	LocationUrl = BaseUrl + "location/"
	transport   = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   0,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
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
	client := &http.Client{
		Transport: transport,
	}

	var data RestaurantResponse

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("X-TripAdvisor-API-Key", "3c7beec8-846d-4377-be03-71cae6145fdc")
	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func FetchReviews(url string) (ReviewResponse, error) {
	client := &http.Client{
		Transport: transport,
	}

	var data ReviewResponse

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("X-TripAdvisor-API-Key", "3c7beec8-846d-4377-be03-71cae6145fdc")
	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
