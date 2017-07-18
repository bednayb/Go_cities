package citiesFunction

import (
	"flag"
	"github.com/bednayb/Go_cities/cityStructs"
	"github.com/bednayb/Go_cities/databases/mockDatabase"
	"github.com/bednayb/Go_cities/databases/productionDatabase"
	"github.com/bednayb/Go_cities/databases/testDatabase"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"fmt"
	"encoding/json"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)



// TODO ez nagyon úgy tűnik mintha a mock adatokat adnánk vissza minden esetben mikor a városokat lekérdezzük! (ready)
// TODO A mock adatokkal való tesztelést különítsük el a valós működéstől, csak akkor induljon mock adatokkal a program ha arra kértük (ready, test not works)
// TODO live/demo setupoláshoz vagy config file-t használjunk, vagy argumentumokat program indításkor (? ez full kodos :))
// Ricsi --> akkor hasznalj mock adatokat ha go run main.go --mock al hivod meg kul, (go run main.go) azzal ami el van mentve

// CitiesInfo is collection of cities
type CitiesInfo []cityStructs.CityInfo

// CityDatabase is collection of cities
var CityDatabase CitiesInfo


// configuration file structure
type Configuration struct {
	Type     string
	Database string
	ProcessorNumber int
}
// ProcessorNumber declaration
var ProcessorNumber int

var db *sql.DB


//ConfigSettings here you can choose which settings file will be used (default is development)
func ConfigSettings(configFile *string) {

	var config = flag.String("config", "", "placeholder")
	flag.Parse()
	switch *config {
	case "production":
		*configFile = "production"
	case "test":
		*configFile = "test"
	default:
		*configFile = "development"
	}
}

// Init before run the program settings config contents
func Init(configFile string) {

	file, _ := os.Open("./config/"+ configFile +".json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	ProcessorNumber = configuration.ProcessorNumber

	switch configuration.Database {
	case "productionDatabase":
		for i := 0; i < len(productionDatabase.Cities); i++ {
			CityDatabase = append(CityDatabase, productionDatabase.Cities[i])
		}
	case "testDatabase":
		for i := 0; i < len(testDatabase.Cities); i++ {
			CityDatabase = append(CityDatabase, testDatabase.Cities[i])
		}
	default:
		for i := 0; i < len(mockDatabase.Cities); i++ {
			CityDatabase = append(CityDatabase, mockDatabase.Cities[i])

		}
	}
}

// SQL MAGIC
func GetAllCitySQL(c *gin.Context) {
	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM City")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	type City struct {
		CityId int `json:"id"`
		CityName string `json:"cityName"`
	}

	var cities []City

	for rows.Next() {
		var CityId  int
		var CityName string

		rows.Scan(&CityId ,&CityName)
		cities = append(cities, City{CityId, CityName})
	}

	fmt.Println("hey",cities)
	c.JSON(200, cities)
}

func PostCitySQL(c *gin.Context) {


	var json cityStructs.CityInfo
	c.Bind(&json) // This will infer what binder to use depending on the content-type header.


	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()


	// insert into City

	stmt, err := db.Prepare("INSERT City SET CityName=? " )
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	res, err := stmt.Exec(json.City)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// insert into Info
	stmt, err = db.Prepare("INSERT CityInfo Set CityId =?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	res, err = stmt.Exec(9)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// insert into Geo
	stmt, err = db.Prepare("INSERT Geo Set Longitude=?, Latitude=?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	res, err = stmt.Exec(json.Geo.Lng,json.Geo.Lat)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}



	fmt.Println("heydd",res)
	c.JSON(200, res)
}

func DeleteCitySQL(c *gin.Context) {

	CityID := c.Query("id")
	CityIDConvertToInt, _ := strconv.ParseInt(CityID, 10, 64)

	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// delete
	stmt, err := db.Prepare("DELETE FROM City WHERE ID=?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	res, err := stmt.Exec(CityIDConvertToInt)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	affect, err := res.RowsAffected()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if affect > 0{
		c.JSON(200, "delete was successful")
		return
	}

	c.JSON(200, res)
}


// SQL MAGIC

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

//GetExpectedForecast count the expected celsius and raining change for next five days
func GetExpectedForecast(c *gin.Context) {
	var wg sync.WaitGroup
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
	citiesDistance := DistanceCounter(ProcessorNumber,presentData, filteredCitiesbyTime)

	// count all distances
	//distances := CountDistance(presentData, filteredCitiesbyTime)

	// balanced the distances
	balancedDistance := BalancedDistanceByLinearInterpolation(citiesDistance)
	// counting temps and raining data for next 5 days
	// Todo csatornaval
	wg.Add(2)
	var forecastRain []float64
	var forecastCelsius []float64
	go CalculateRain(balancedDistance, filteredCitiesbyTime, &forecastRain, &wg)
	go CalculateTemp(balancedDistance, filteredCitiesbyTime, &forecastCelsius, &wg)
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

// BalancedDistanceByLinearInterpolation ponderare by linear interpolation (nearest 1 weight, furthest 0)
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
func CalculateRain(balancedCityDistance map[string]float64, cityInfo map[string]cityStructs.CityInfo, ForecastRainingChance *[]float64, databaseWaitGroup *sync.WaitGroup) {

	var totalBalance float64
	var totalTemp float64

	// count next five days
	for day := 0; day < 5; day++ {
		totalBalance = 0
		totalTemp = 0
		// v --> every city
		for _, v := range cityInfo {
			totalBalance += balancedCityDistance[v.City]
			totalTemp += v.Rain[day] * balancedCityDistance[v.City]
		}
		// cut off 2 decimal
		untruncated := totalTemp / totalBalance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		*ForecastRainingChance = append(*ForecastRainingChance, truncated)
	}
	databaseWaitGroup.Done()

}

//CalculateTemp where we count the expected Celsius chance for next five days
func CalculateTemp(balancedCityDistance map[string]float64, cityInfo map[string]cityStructs.CityInfo, ForecastTemps *[]float64, databaseWaitGroup *sync.WaitGroup) {

	var totalBalance float64
	var totalTemp float64

	// count next five days
	for day := 0; day < 5; day++ {
		totalBalance = 0
		totalTemp = 0
		// info --> every city
		for _, v := range cityInfo {
			totalBalance += balancedCityDistance[v.City]
			totalTemp += v.Temp[day] * balancedCityDistance[v.City]
		}
		// cut off 2 decimal
		untruncated := totalTemp / totalBalance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		*ForecastTemps = append(*ForecastTemps, truncated)

	}
	databaseWaitGroup.Done()
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

// PostCity saving new city
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

//Todo pointer helyett channeleket irj,
//Todo 1. feldolgozo Process ( StartDatabaseWritingNode)
//Todo 2. feldolgozando elemeket tartalmazo csatorna letrehozasa
// Todo 3. response elemeket tartalmazo csatorna letrehozasa
//Todo 4.   eleinditasz barmennyit
// Todo 5. ciklus ami a valaszcsatornat dolgozza fel

// DistanceCounter where we count every city's distance from an exact place
func DistanceCounter(ProcessorNumber int, coordinate cityStructs.CoordinateAndTime, filteredCities map[string]cityStructs.CityInfo) (distanceCities map[string]float64) {

	var wg sync.WaitGroup


	// because of the append we need to declare here by make
	result := make(map[string]float64)

	///////Todo buffer annyi legyen mint ahany csatorna van (ready)
	//make channel
	in := make(chan cityStructs.CityInfo, len(filteredCities))
	out := make(chan Out, len(filteredCities))


	//////Todo megadhato proc szam (configbol) (ready)
	//make processor
	for i:= 0 ; i < ProcessorNumber; i++{
		go DistanceCounterProcess(in, coordinate, out, &wg)
	}
	// Send data until left
	for _, cityInfo := range filteredCities {
		wg.Add(1)
		in <- cityInfo
	}
	go func() {
		for {
			select {
			case res := <-out:
				result[res.CityName] = res.Distance
				wg.Done()
			}
		}
	}()
	wg.Wait()
	return result
}

//DistanceCounterProcess count the distances of city and send back to DistanceCounter
func DistanceCounterProcess(in chan cityStructs.CityInfo, coordinate cityStructs.CoordinateAndTime, out chan Out, wg *sync.WaitGroup) {

	var distance float64

	for {
		select {
		case cityInfo := <-in:
			defer wg.Done()
			// count distance
			latitudeDistance := coordinate.Lat - cityInfo.Geo.Lat
			longitudeDistance := coordinate.Lng - cityInfo.Geo.Lat
			distance = math.Sqrt(math.Pow(latitudeDistance, 2) + math.Pow(longitudeDistance, 2))

			//response data Out type and make map just at other side because with map can be gorutine problems (see at type Out)
			res := Out{cityInfo.City, distance}
			//send back
			out <- res
		}
	}
}

// necesseary to not send back map because
// if one goroutine is writing to a map, no other goroutine should be reading or writing the map concurrently. If the runtime detects this condition, it prints a diagnosis and crashes the program. (https://golang.org/doc/go1.6#runtime)
type Out struct {
	CityName string
	Distance float64
}
