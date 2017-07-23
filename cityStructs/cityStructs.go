package cityStructs

// CityInfo is data about city
type CityInfo struct {
	City      string
	Geo       Geo
	Temp      [5]float64
	Rain      [5]float64
	Timestamp int64
}

// CitiesInfo is collection of cities
type CitiesInfo []CityInfo

// CoordinateAndTime is coordinates of city and time
type CoordinateAndTime struct {
	Lat       float64
	Lng       float64
	Timestamp int64
}

// Geo is coordinates of city
type Geo struct {
	Lat float64
	Lng float64
}

// Configuration file structure
type Configuration struct {
	Type            string
	Name  			string
	MySQL bool
	Database        Database
	ProcessorNumber int
	Port string
	Calculation Calculation
	FilteringCityData bool
	BalancedByDistance bool
}

type Database struct {
	Name string
	MySQL bool
	Username string
	Password string
}
type Calculation struct {
	FilteringCityData bool
	BalancedByDistance bool
}

//Out is necessary to not send back map because, if one goroutine is writing to a map, no other goroutine should be reading or writing the map concurrently. If the runtime detects this condition, it prints a diagnosis and crashes the program. (https://golang.org/doc/go1.6#runtime)
type Out struct {
	CityName string
	Distance float64
}
