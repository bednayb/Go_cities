package cityStructs

// CityInfo is data about city
type CityInfo struct {
	City      string `json:"City"`
	Geo       Geo
	Temp      [5]float64 `json:"Temp"`
	Rain      [5]float64 `json:"Rain"`
	Timestamp int64      `json:"Timestamp"`
}

// CitiesInfo is collection of cities
type CitiesInfo []CityInfo

// CoordinateAndTime is coordinates of city and time
type CoordinateAndTime struct {
	Lat       float64 `json:"Lat"`
	Lng       float64 `json:"Lng"`
	Timestamp int64   `json:"Timestamp"`
}

// Geo is coordinates of city
type Geo struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}

// CityData is contains cityInfo data from slq db
type CityData struct {
	CityID    int
	InfoID    int
	Date      int
	Temp      string
	Rain      string
	Latitude  float64
	Longitude float64
}