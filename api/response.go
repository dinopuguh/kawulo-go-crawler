package api

type Restaurant struct {
	LocationId         string        `json:"location_id"`
	Name               string        `json:"name"`
	Latitude           string        `json:"latitude"`
	Longitude          string        `json:"longitude"`
	NumReviews         string        `json:"num_reviews"`
	Photo              Photo         `json:"photo"`
	Rating             string        `json:"rating"`
	PriceLevel         string        `json:"price_level"`
	Price              string        `json:"price"`
	Address            string        `json:"address"`
	Phone              string        `json:"phone"`
	Website            string        `json:"website"`
	RawRanking         string        `json:"raw_ranking"`
	RankingGeo         string        `json:"ranking_geo"`
	RankingPosition    string        `json:"ranking_position"`
	RankingDenominator string        `json:"ranking_denominator"`
	RankingCategory    string        `json:"ranking_category"`
	Ranking            string        `json:"ranking"`
	SubCategory        []SubCategory `json:"subcategory"`
}

type Photo struct {
	Images Images `json:"images"`
}
type Images struct {
	Thumbnail Image `json:"thumbnail"`
	Original  Image `json:"original"`
	Medium    Image `json:"medium"`
	Large     Image `json:"large"`
}
type Image struct {
	Width  string `json:"width"`
	Url    string `json:"url"`
	Height string `json:"height"`
}

type SubCategory struct {
	Key  string `json:"key"`
	Name string `json:"name"`
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
