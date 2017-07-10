package city_structs
// TODO gorm? nem szükséges egyelőre sehová, ezeket tisztítsuk ki
// TODO form ot használjuk valamire?
// ebből a sorból töröltem a, amit már végrehajtottál

type CityInfo struct {
	City      string `gorm:"not null" form:"City" json:"City"`
	Geo       Geo
	Temp      [5]float64 `gorm:"not null" form:"Temp"json:"Temp"`
	Rain      [5]float64 `gorm:"not null" form:"Rain"json:"Rain"`
	Timestamp int64      `gorm:"not null" form:"Timestamp"json:"Timestamp"` // TODO olvashatóság kedvéért legyen space a form és a json között
}

type CitiesInfo []CityInfo

type Cordinate_and_time struct {
	Lat       float64 `json:"Lat"`
	Lng       float64 `json:"Lng"`
	Timestamp int64   `gorm:"not null" form:"Timestamp"json:"Timestamp"`
}

type Geo struct {
	Lat float64 `gorm:"not null" form:"Lat" json:"Lat"`
	Lng float64 `gorm:"not null" form:"Lng"json:"Lng"`
}
