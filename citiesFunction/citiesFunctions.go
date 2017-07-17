package citiesFunction

import (
	"flag"
	"github.com/bednayb/Go_cities/cityStructs"
	"github.com/bednayb/Go_cities/databases/mockDatabase"
	"github.com/bednayb/Go_cities/databases/testDatabase"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"github.com/bednayb/Go_cities/databases/productionDatabase"
)

// TODO ez nagyon úgy tűnik mintha a mock adatokat adnánk vissza minden esetben mikor a városokat lekérdezzük! (ready)
// TODO A mock adatokkal való tesztelést különítsük el a valós működéstől, csak akkor induljon mock adatokkal a program ha arra kértük (ready, test not works)
// TODO live/demo setupoláshoz vagy config file-t használjunk, vagy argumentumokat program indításkor (? ez full kodos :))
// Ricsi --> akkor hasznalj mock adatokat ha go run main.go --mock al hivod meg kul, (go run main.go) azzal ami el van mentve


// CitiesInfo is collection of cities
type CitiesInfo []cityStructs.CityInfo

// CityDatabase is collection of cities
var CityDatabase CitiesInfo

var mutex sync.Mutex

var wg sync.WaitGroup

var Counter = 0

//ConfigSettings here you can choose whice settings file will be used
func ConfigSettings(configFile *string) {

	var config = flag.String("config", "", "placeholder")
	flag.Parse()
	if *config == "production" {
		*configFile = "production"
	} else if *config == "test" {
		*configFile = "test"
	} else {
		*configFile = "development"
	}
}

func Init(conf string) {
	if conf == "development" {
		for i:=0; i < len(mockDatabase.Cities);i++ {
			CityDatabase = append(CityDatabase, mockDatabase.Cities[i])
		}
	} else if conf == "test" {
		for i := 0; i < len(testDatabase.Cities); i++ {
			CityDatabase = append(CityDatabase, testDatabase.Cities[i])
		}
	} else if conf == "production" {
		for i := 0; i < len(productionDatabase.Cities); i++ {
			CityDatabase = append(CityDatabase, productionDatabase.Cities[i])
		}
	}
}




// GetAllCity shows all cities
func GetAllCity(c *gin.Context) {
	cities := CityDatabase
	c.JSON(200, cities)
}

// GetCityByName shows every data where the city name is same (example: becs)
func GetCityByName(c *gin.Context) {

	cities := CityDatabase
	// find city's name from url
	name := c.Params.ByName("name")
	// bool for checking city is exist in our db
	redFlag := true

	// filtered cities order by timestamp (first the oldest)
	var filteredCitiesByTime CitiesInfo

	// filtering cities by name
	for _, v := range cities {
		if v.City == name {
			redFlag = false
			filteredCitiesByTime = append(filteredCitiesByTime, v)
		}
	}
	if redFlag {
		// response when city doesnt exist in our db
		content := gin.H{"error": "city with name " + name + " not found"}
		c.JSON(404, content)
		return
	}

	// sorting cities
	sort.Sort(filteredCitiesByTime)
	// response when city exist in our db
	c.JSON(200, gin.H{"filteredCitiesByTime": filteredCitiesByTime})
	return

	// TODO érdemes lenne mindkét if ágban egy return, hogy ide ne juthassunk el. (ready)
	// Ha itt bármilyen kód lenne független attól hogy not found volt e lefutna!
}

// TODO ennek a fügvénynek a neve nem tükrözi hogy valójában mit csinál  (ready)
func GetExpectedForecast(c *gin.Context) {

	if len(CityDatabase) == 0 {
		content := gin.H{"response": "sry we havnt had enough data for calculating yet"}
		c.JSON(200, content)
		return
	}

	// save data from URL
	lat := c.Query("lat")
	lng := c.Query("lng")
	timestamp := c.Query("timestamp")

	// TODO Hiba ellenőrzéskor értelmes hibaüzenetet szeretnénk adni pontosan arról ami citiesDistance hibát okozta (ready)

	var dataDoenstExistsMessage string

	if lat == "" {
		dataDoenstExistsMessage += "lat data must be exists, "
	}
	if lng == "" {
		dataDoenstExistsMessage += "lng data must be exists, "
	}
	if timestamp == "" {
		dataDoenstExistsMessage += "timestamp data must be exists"
	}
	if len(dataDoenstExistsMessage) > 0 {
		content := gin.H{"error_message": dataDoenstExistsMessage}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk citiesDistance futást ha hiba volt  (ready)
		return
	}

	//Convert to float64/int
	var convertProblem string
	latitudeConvertToFloat64, _ := strconv.ParseFloat(strings.TrimSpace(lat), 64)

	if lat != "0" && latitudeConvertToFloat64 == 0 {
		convertProblem += "invalid lat data (not number), "
	}

	longitudeConvertToFloat64, _ := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if lat != "0" && longitudeConvertToFloat64 == 0 {
		convertProblem += "invalid lng data (not number), "
	}

	timestampConvertToInt, _ := strconv.ParseInt(timestamp, 10, 64)
	if lat != "0" && timestampConvertToInt == 0 {
		convertProblem += "invalid timestamp data (not number) "
	}

	if len(convertProblem) > 0 {
		content := gin.H{"error_message ": convertProblem}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk citiesDistance futást ha hiba volt  (ready)
		return
	}

	if timestampConvertToInt < 0 {
		content := gin.H{"error_message ": "timestamp should be bigger than 0"}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk citiesDistance futást ha hiba volt  (ready)
		return
	}

	//put data to struct
	// TODO a fenti parsolások mindegyikénél előfordulhat hiba, amit így teljesen figyelmen kívűl hagyunk (ready)
	// TODO a fentabbi ellenőrzési szisztémával adhatunk hibaüzenetet hogy melyikből nem sikerült számot kinyernünk. (ready)  (timestampnel minuszt nem fogadunk el -- Ricsi)
	// + Ellenőrizhető hogy citiesDistance szám valós tartományban van e.
	// hiba esetén itt se menjünk tovább.

	// data from the URL
	presentData := cityStructs.CoordinateAndTime{latitudeConvertToFloat64, longitudeConvertToFloat64, timestampConvertToInt}

	// filter for the nearest data (by timestamp)
	filteredCitiesbyTime := NearestCityDataInTime(CityDatabase, timestampConvertToInt)

	// count all distance with channels
	citiesDistance := DistanceCounter(15, presentData, filteredCitiesbyTime)

	// count all distances
	//distances := CountDistance(presentData, filteredCitiesbyTime)

	// balanced the distances
	balancedDistance := BalancedDistanceByLinearInterpolation(citiesDistance)

	// counting temps and raining data for next 5 days
	// Todo csatornaval
	wg.Add(2)
	var forecastRain []float64
	var forecastCelsius []float64
	go CalculateRain(balancedDistance, filteredCitiesbyTime, &forecastRain)
	go CalculateTemp(balancedDistance, filteredCitiesbyTime, &forecastCelsius)
	wg.Wait()

	// send data
	content := gin.H{"expected celsius next 5 days": forecastCelsius, "expected rainning chance next 5 days": forecastRain}
	c.JSON(200, content)
}

//func CountCitiesDistance(coordinate cityStructs.CoordinateAndTime, info map[string]cityStructs.CityInfo, a *map[string]float64) {
//
//	// container for distance  key --> city name, value --> distance
//	var citiesDistance = make(map[string]float64)
//
//	//count every distance of city (pitágoras)
//	var distance float64
//	for _, info := range info {
//
//		latitudeDistance := coordinate.Lat - info.Geo.Lat
//		longitudeDistance := coordinate.Lng - info.Geo.Lng
//
//		distance = math.Sqrt(math.Pow(latitudeDistance, 2) + math.Pow(longitudeDistance, 2))
//		citiesDistance[info.City] = distance
//	}
//
//	*a = mergeMaps(*a, citiesDistance)
//	wg.Done()
//}

// BalanceDistance ponderare by linear interpolation (nearest 1 weight, furthest 0)
func BalancedDistanceByLinearInterpolation(distances map[string]float64) (balanceByDistance map[string]float64) {

	//  balanced distance
	var balanceNumber float64

	//// find furthest (biggest number)
	var permanentBiggest float64
	var biggest float64

	for _, v := range distances {
		if v > permanentBiggest {
			permanentBiggest = v
			biggest = permanentBiggest
		}
	}
	//find nearest (smallest number)
	permanentSmallest := biggest
	smallest := biggest

	for _, v := range distances {
		if v < permanentSmallest {
			permanentSmallest = v
			smallest = permanentSmallest
		}
	}
	// calculate balanced numbers
	for i, v := range distances {
		balanceNumber = (v - smallest) / (biggest - smallest)
		balanceNumber--
		balanceNumber *= -1
		// overwrite distance with balanced distance
		distances[i] = balanceNumber
	}

	return distances
}
//CalculateRain where we count the expected raining chance for next five days
func CalculateRain(balance map[string]float64, cityInfo map[string]cityStructs.CityInfo, a *[]float64) {

	var totalBalance float64
	var totalTemp float64

	// count next five days
	for day := 0; day < 5; day++ {
		totalBalance = 0
		totalTemp = 0
		// v --> every city
		for _, v := range cityInfo {
			totalBalance += balance[v.City]
			totalTemp += v.Rain[day] * balance[v.City]
		}
		// cut off 2 decimal
		untruncated := totalTemp / totalBalance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		*a = append(*a, truncated)
	}
	wg.Done()

}
// BalanceDistance ponderare by linear interpolation (nearest 1 weight, furthest 0)
func CalculateTemp(balance map[string]float64, cityInfo map[string]cityStructs.CityInfo, a *[]float64) {

	var totalBalance float64
	var totalTemp float64

	// count next five days
	for day := 0; day < 5; day++ {
		totalBalance = 0
		totalTemp = 0
		// info --> every city
		for _, v := range cityInfo {
			totalBalance += balance[v.City]
			totalTemp += v.Temp[day] * balance[v.City]
		}
		// cut off 2 decimal
		untruncated := totalTemp / totalBalance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		*a = append(*a, truncated)

	}
	wg.Done()
}

// TODO az alábbi 3 fügvényt a tructok mellett tárolnám hogy (ready)
// egyben látszódjon egy egy adattípusról, hogy mik az elemei és mik a rá definiált fugvények  (?)
// order Cities by Timestamp
func (slice CitiesInfo) Len() int {
	return len(slice)
}

func (slice CitiesInfo) Less(i, j int) bool {
	return slice[i].Timestamp < slice[j].Timestamp
}

func (slice CitiesInfo) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// saving new city
func PostCity(c *gin.Context) {

	var json cityStructs.CityInfo
	c.Bind(&json) // This will infer what binder to use depending on the content-type header.

	// checking rain data
	for _, v := range json.Rain {
		if v < 0 || v > 1 {
			c.JSON(400, gin.H{
				"result": "Failed, invalid temp data (should be beetween 0 and 1)",
			})
			return
		}
	}

	CityDatabase = append(CityDatabase, json)

	content := gin.H{
		"result": "successful saving",
	}
	c.JSON(201, content)
}

// TODO használjunk visszatérési érték változónevet is. (ready)

// NearestCityDataInTime is a filter where we get back just one city (exm if we have 3 becs back just one) which is the most relevant by time
func NearestCityDataInTime(allCities []cityStructs.CityInfo, timestamp int64) (filteredCities map[string]cityStructs.CityInfo) {
	// TODO én MAP ez használnék ahol a város neve a kulcs  (ready)
	// és mindenhol az érték felülírása akkor történhet meg ha az infó frissebb.

	citiesDistance := make(map[string]cityStructs.CityInfo)

	for _, v := range allCities {

		oldDataCityDistanceTime := citiesDistance[v.City].Timestamp - timestamp
		if oldDataCityDistanceTime < 0 {
			oldDataCityDistanceTime *= -1
		}

		newDataCityDistanceTime := v.Timestamp - timestamp
		if newDataCityDistanceTime < 0 {
			newDataCityDistanceTime *= -1
		}

		if oldDataCityDistanceTime > newDataCityDistanceTime || filteredCities[v.City].Timestamp == 0 {
			citiesDistance[v.City] = v
		}
	}
	return citiesDistance
}

// linear interpolation (nearest 1 weight, furthest 0)
//func mergeMaps(x map[string]float64, y map[string]float64) map[string]float64 {
//	for k, v := range x {
//		y[k] = v
//	}
//	return y
//}

//func CountDistance(currentPlaceAndTime cityStructs.CoordinateAndTime, filteredCities map[string]cityStructs.CityInfo) map[string]float64 {
//
//	wg.Add(2)
//	var distances map[string]float64
//
//	cityHalf1 := make(map[string]cityStructs.CityInfo)
//	cityHalf2:= make(map[string]cityStructs.CityInfo)
//
//	cutter := 0
//	for key, val := range filteredCities {
//		if cutter%2 == 0 {
//			cityHalf1[key] = val
//		} else {
//			cityHalf2[key] = val
//		}
//		cutter ++
//	}
//
//	go CountCitiesDistance(currentPlaceAndTime, cityHalf1, &distances)
//	go CountCitiesDistance(currentPlaceAndTime, cityHalf2, &distances)
//	wg.Wait()
//	return distances
//}

//Todo pointer helyett channeleket irj,
//Todo 1. feldolgozo Process ( StartDatabaseWritingNode)
//Todo 2. feldolgozando elemeket tartalmazo csatorna letrehozasa
// Todo 3. response elemeket tartalmazo csatorna letrehozasa
//Todo 4.   eleinditasz barmennyit
// Todo 5. ciklus ami a valaszcsatornat dolgozza fel

// DistanceCounter where we count every city's distance from an exact place
func DistanceCounter(procNumber int, cordinate cityStructs.CoordinateAndTime, filteredCities map[string]cityStructs.CityInfo) (distanceCities map[string]float64) {

	var wg sync.WaitGroup

	// because of the append we need to declare here by make
	result := make(map[string]float64)

	// contains every filtered city's name
	var names []string
	for _, v := range filteredCities {
		names = append(names, v.City)
	}

	//channels
	in := make(chan chan map[string]float64)

	//make processor
	for i := 0; i < procNumber; i++ {
		go DistanceCounterProcess(in, cordinate, filteredCities, names)
	}

	// Send data until left
	for n := 0; n < len(filteredCities); n++ {

		mutex.Lock()
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := make(chan map[string]float64)
			in <- c
			z := <-c

			// merge maps
			for k, v := range z {
				result[k] = v
			}
			mutex.Unlock()
		}()
	}

	wg.Wait()
	// its necceseary to Counter equal to 0, because we can have more query and can be out of range
	Counter = 0
	return result

}

func DistanceCounterProcess(in chan chan map[string]float64, cordinate cityStructs.CoordinateAndTime, filteredCities map[string]cityStructs.CityInfo, names []string) {

	var distance float64
	result := make(map[string]float64)

	for in := range in {

		latitudeDistance := cordinate.Lat - filteredCities[names[Counter]].Geo.Lat
		longitudeDistance := cordinate.Lng - filteredCities[names[Counter]].Geo.Lat

		distance = math.Sqrt(math.Pow(latitudeDistance, 2) + math.Pow(longitudeDistance, 2))
		result[filteredCities[names[Counter]].City] = distance
		Counter++
		in <- result

	}
}
