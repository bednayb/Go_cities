package citiesFunction

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bednayb/Go_cities/cityStructs"
	"github.com/bednayb/Go_cities/databases/mockDatabase"
	"github.com/bednayb/Go_cities/databases/productionDatabase"
	"github.com/bednayb/Go_cities/databases/testDatabase"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// CitiesInfo is collection of cities
type CitiesInfo []cityStructs.CityInfo

// CityDatabase is collection of cities
var CityDatabase CitiesInfo

// Configuration file structure
type Configuration struct {
	Type            string
	Database        string
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

	file, _ := os.Open("./config/" + configFile + ".json")
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

// GetAllCitySQL list all cities from SQL database
func GetAllCitySQL(c *gin.Context) {
	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM City")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	type City struct {
		CityID   int    `json:"id"`
		CityName string `json:"cityName"`
	}

	var cities []City

	for rows.Next() {
		var CityID int
		var CityName string

		rows.Scan(&CityID, &CityName)
		cities = append(cities, City{CityID, CityName})
	}

	c.JSON(200, cities)
}

// PostCitySQL add new city to SQL database
func PostCitySQL(c *gin.Context) {

	var json cityStructs.CityInfo
	c.Bind(&json) // This will infer what binder to use depending on the content-type header.

	// open the database
	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// if city not exist the
	var CheckCityID int

	db.QueryRow("SELECT ID FROM City WHERE CityName = ?", json.City).Scan(&CheckCityID)

	// if city not exist in our DB,id == 0 and the new city will be saved
	if CheckCityID == 0 {
		stmt, err := db.Prepare("INSERT City SET CityName = ?")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		res, err := stmt.Exec(json.City)
		if res == nil {
			fmt.Println(res)
		}
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	// cityId for cityInfo
	var cityID int
	db.QueryRow("SELECT ID FROM City WHERE CityName = ?", json.City).Scan(&cityID)

	// insert into Info
	stmt, err := db.Prepare("INSERT CityInfo Set CityId =?,Date=?,Temp=?,Rain=?,Latitude=?,Longitude=?")

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	tempData := json.Temp
	var tempDataToSQL string
	for i, v := range tempData {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if i != len(tempData)-1 {
			tempDataToSQL += s + ","
		} else {
			tempDataToSQL += s
		}
	}

	rainData := json.Rain
	var rainDataToSQL string
	for i, v := range rainData {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if i != len(tempData)-1 {
			rainDataToSQL += s + ","
		} else {
			rainDataToSQL += s
		}
	}

	res, err := stmt.Exec(cityID, json.Timestamp, tempDataToSQL, rainDataToSQL, json.Geo.Lat, json.Geo.Lng)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	go c.JSON(200, res)
}

//DeleteCitySQL delete city by id from SQL database (havnt worked perfectly yet)
func DeleteCitySQL(c *gin.Context) {

	CityID := c.Query("id")
	CityIDConvertToInt, _ := strconv.ParseInt(CityID, 10, 64)

	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
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

	if affect > 0 {
		c.JSON(200, "delete was successful")
		return
	}

	c.JSON(200, res)
}

// GetCityByIDSQL find every info from sql db about city by id
func GetCityByIDSQL(c *gin.Context) {
	id := c.Params.ByName("id")
	idConvertToInt, _ := strconv.ParseInt(id, 10, 64)

	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM CityInfo WHERE CityId = ? ORDER BY DATE", idConvertToInt)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var cities []cityStructs.CityData
	for rows.Next() {

		var InfoID int
		var CityID int
		var Date int
		var Temp string
		var Rain string
		var Latitude float64
		var Longitude float64

		rows.Scan(&InfoID, &CityID, &Date, &Temp, &Rain, &Latitude, &Longitude)
		// add every row to cities
		cities = append(cities, cityStructs.CityData{CityID, InfoID, Date, Temp, Rain, Latitude, Longitude})
	}
	c.JSON(200, cities)
}

// GetExpectedForecastSQL makes forecast for exact place
func GetExpectedForecastSQL(c *gin.Context) {
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

	var dataDoesntExistsMessage string

	if lat == "" {
		dataDoesntExistsMessage += "lat data must be exists, "
	}
	if lng == "" {
		dataDoesntExistsMessage += "lng data must be exists, "
	}
	if timestamp == "" {
		dataDoesntExistsMessage += "timestamp data must be exists"
	}
	if len(dataDoesntExistsMessage) > 0 {
		content := gin.H{"error_message": dataDoesntExistsMessage}
		c.JSON(400, content)
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
		return
	}

	if timestampConvertToInt < 0 {
		content := gin.H{"error_message ": "timestamp should be bigger than 0"}
		c.JSON(400, content)
		return
	}

	// data from the URL
	presentData := cityStructs.CoordinateAndTime{latitudeConvertToFloat64, longitudeConvertToFloat64, timestampConvertToInt}

	//filtered Cities
	cities := CitiesFromSQL(timestampConvertToInt)

	CitiesDataConvertToMap := DataToMap(cities)

	citiesDistance := DistanceCounter(ProcessorNumber, presentData, CitiesDataConvertToMap)

	// count all distance with channels
	//citiesDistance := DistanceCounter(ProcessorNumber, presentData, CitiesDataConvertToMap)

	// balanced the distances
	balancedDistance := BalancedDistanceByLinearInterpolation(citiesDistance)
	// counting temps and raining data for next 5 days
	wg.Add(2)
	var forecastRain []float64
	var forecastCelsius []float64
	go CalculateRain(balancedDistance, CitiesDataConvertToMap, &forecastRain, &wg)
	go CalculateTemp(balancedDistance, CitiesDataConvertToMap, &forecastCelsius, &wg)
	wg.Wait()

	// send data
	content := gin.H{"expected celsius next 5 days": forecastCelsius, "expected raining chance next 5 days": forecastRain}
	c.JSON(200, content)

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
}

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

	var dataDoesntExistsMessage string

	if lat == "" {
		dataDoesntExistsMessage += "lat data must be exists, "
	}
	if lng == "" {
		dataDoesntExistsMessage += "lng data must be exists, "
	}
	if timestamp == "" {
		dataDoesntExistsMessage += "timestamp data must be exists"
	}
	if len(dataDoesntExistsMessage) > 0 {
		content := gin.H{"error_message": dataDoesntExistsMessage}
		c.JSON(400, content)
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
		return
	}

	if timestampConvertToInt < 0 {
		content := gin.H{"error_message ": "timestamp should be bigger than 0"}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk citiesDistance futást ha hiba volt  (ready)
		return
	}

	// data from the URL
	presentData := cityStructs.CoordinateAndTime{latitudeConvertToFloat64, longitudeConvertToFloat64, timestampConvertToInt}

	// filter for the nearest data (by timestamp)
	filteredCitiesbyTime := NearestCityDataInTime(CityDatabase, timestampConvertToInt)

	// count all distance with channels
	citiesDistance := DistanceCounter(ProcessorNumber, presentData, filteredCitiesbyTime)

	// balanced the distances
	balancedDistance := BalancedDistanceByLinearInterpolation(citiesDistance)
	// counting temps and raining data for next 5 days

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

// NearestCityDataInTime is a filter where we get back just one city (exm if we have 3 becs back just one) which is the most relevant by time
func NearestCityDataInTime(allCities []cityStructs.CityInfo, timestamp int64) (filteredCities map[string]cityStructs.CityInfo) {

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

//CitiesFromSQL make a query to database for cities
func CitiesFromSQL(timestamp int64) (filteredCities map[string]cityStructs.CityData) {

	// open SQL
	db, err := sql.Open("mysql", "root:admin@/GoCities")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

	rows, err := db.Query("select * from CityInfo")

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// container of cities
	var cities []cityStructs.CityData

	for rows.Next() {

		var InfoID int
		var CityID int
		var Date int
		var Temp string
		var Rain string
		var Latitude float64
		var Longitude float64

		rows.Scan(&InfoID, &CityID, &Date, &Temp, &Rain, &Latitude, &Longitude)
		// add every row to cities
		cities = append(cities, cityStructs.CityData{CityID, InfoID, Date, Temp, Rain, Latitude, Longitude})

	}

	filteredCitiesbyDistance := make(map[string]cityStructs.CityData)
	timestampToInt := int(timestamp)

	for _, v := range cities {

		cityID := strconv.Itoa(v.CityID)

		var cityDistanceTimestamp = int(filteredCitiesbyDistance[cityID].Date)

		oldDataCityDistanceTime := cityDistanceTimestamp - timestampToInt
		if oldDataCityDistanceTime < 0 {
			oldDataCityDistanceTime *= -1
		}

		newDataCityDistanceTime := v.Date - timestampToInt
		if newDataCityDistanceTime < 0 {
			newDataCityDistanceTime *= -1
		}

		if oldDataCityDistanceTime > newDataCityDistanceTime || filteredCitiesbyDistance[cityID].Date == 0 {
			filteredCitiesbyDistance[cityID] = v
		}
	}
	return filteredCitiesbyDistance
}

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
	for i := 0; i < ProcessorNumber; i++ {
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
			longitudeDistance := coordinate.Lng - cityInfo.Geo.Lng
			distance = math.Sqrt(math.Pow(latitudeDistance, 2) + math.Pow(longitudeDistance, 2))

			//response data Out type and make map just at other side because with map can be gorutine problems (see at type Out)
			res := Out{cityInfo.City, distance}
			//send back
			out <- res
		}
	}
}

//Out is necessary to not send back map because, if one goroutine is writing to a map, no other goroutine should be reading or writing the map concurrently. If the runtime detects this condition, it prints a diagnosis and crashes the program. (https://golang.org/doc/go1.6#runtime)
type Out struct {
	CityName string
	Distance float64
}

//DataToMap change sql data format
func DataToMap(filteredCitiesFromSQLDb map[string]cityStructs.CityData) (filteredCitiesResult map[string]cityStructs.CityInfo) {

	filteredCities := make(map[string]cityStructs.CityInfo)

	for _, v := range filteredCitiesFromSQLDb {

		// cityID is a unique data, --> use for key value
		// convert int to string
		cityID := strconv.Itoa(v.CityID)

		// Rain and Temp data is in a string first split up,
		stringSliceTemp := strings.SplitN(v.Temp, ",", 5)
		stringSliceRain := strings.SplitN(v.Rain, ",", 5)

		var stringToFloatTemp = [5]float64{}
		var stringToFloatRain = [5]float64{}

		// convert Temp data to float and put into array
		for i, v := range stringSliceTemp {
			f, _ := strconv.ParseFloat(v, 64)
			stringToFloatTemp[i] = f
		}
		// convert Rain data to float and put into array
		for i, v := range stringSliceRain {
			f, _ := strconv.ParseFloat(v, 64)
			stringToFloatRain[i] = f
		}
		// cut off Geo data to decimal
		truncatedLatitude := float64(int(v.Latitude*100)) / 100
		truncatedLongtitude := float64(int(v.Longitude*100)) / 100

		date := int64(v.Date)

		filteredCities[cityID] = cityStructs.CityInfo{cityID, cityStructs.Geo{truncatedLatitude, truncatedLongtitude}, stringToFloatTemp, stringToFloatRain, date}
	}
	return filteredCities
}
