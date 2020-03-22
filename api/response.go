package api

type Restaurant struct {
	LocationId string `json:"location_id"`
	Name       string `json:"name"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	NumReviews string `json:"num_reviews"`
	LocationID string `json:"locationID"`
}

type Review struct {
	ReviewId      string      `json:"id"`
	Lang          string      `json:"lang"`
	LocationId    string      `json:"location_id"`
	PublishedDate string      `json:"published_date"`
	Rating        string      `json:"rating"`
	Text          string      `json:"text"`
	Subratings    []Subrating `json:"subratings"`
}

type Subrating struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
