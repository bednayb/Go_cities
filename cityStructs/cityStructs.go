package cityStructs
// TODO gorm? nem szükséges egyelőre sehová, ezeket tisztítsuk ki (ready)
// TODO form ot használjuk valamire? (ready)
// ebből a sorból töröltem a, amit már végrehajtottál

type CityInfo struct {
	City      string `json:"City"`
	Geo       Geo
	Temp      [5]float64 `json:"Temp"`
	Rain      [5]float64 `json:"Rain"`
	Timestamp int64      `json:"Timestamp"` // TODO olvashatóság kedvéért legyen space a form és a json között (ready)
}

type CitiesInfo []CityInfo

type CoordinateAndTime struct {
	Lat       float64 `json:"Lat"`
	Lng       float64 `json:"Lng"`
	Timestamp int64   `json:"Timestamp"`
}

type Geo struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}
