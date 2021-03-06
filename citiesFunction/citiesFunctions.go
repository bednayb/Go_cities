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
var CityDatabase2 CitiesInfo

// ProcessorNumber declaration
var Config cityStructs.Configuration

//ConfigSettings here you can choose which settings file will be used (default is development)
func ConfigSettings(configFile *string) {

	var config = flag.String("config", "", "placeholder")
	flag.Parse()

	*configFile = *config
}

// Init before run the program settings config contents
func Init(configFile string) {

	file, _ := os.Open("./config/" + configFile + ".json")
	decoder := json.NewDecoder(file)

	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("error:", err)
	}

	//Todo find better solution than switch
	switch Config.Name {
	case "productionDatabase":
			CityDatabase = productionDatabase.Cities
	case "testDatabase":
			CityDatabase = testDatabase.Cities
	default:
			CityDatabase =  mockDatabase.Cities
	}
}

// GetAllCitySQL list all cities from SQL database
func GetAllCity(c *gin.Context) {

	if Config.Database.MySQL {
		cities := CitiesFromSQL()
		c.JSON(200, cities)
	} else {
		cities := CityDatabase
		c.JSON(200, cities)
	}

}

// PostCitySQL add new city to SQL database
func PostCity(c *gin.Context) {
	if Config.Database.MySQL {
		var json cityStructs.CityInfo
		c.Bind(&json) // This will infer what binder to use depending on the content-type header.

		// checking rain data
		for _, v := range json.Rain {
			if v < 0 || v > 1 {
				c.JSON(400, gin.H{
					"result": "Failed, invalid Rain data (should be beetween 0 and 1)",
				})
				return
			}
		}

		// open the database
		db, err := sql.Open("mysql",Config.Database.Username+":"+Config.Database.Password+"@/"+Config.Database.Name )
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

		tempDataToSQL := FloatArrayConvertToString(json.Temp)
		rainDataToSQL := FloatArrayConvertToString(json.Rain)

		res, err := stmt.Exec(cityID, json.Timestamp, tempDataToSQL, rainDataToSQL, json.Geo.Lat, json.Geo.Lng)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		go c.JSON(200, res)
	} else {
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
}

//DeleteCitySQL delete city by id from SQL database (havnt worked perfectly yet)
func DeleteCitySQL(c *gin.Context) {

	CityID := c.Query("id")
	CityIDConvertToInt, _ := strconv.ParseInt(CityID, 10, 64)

	db, err := sql.Open("mysql",Config.Database.Username+":"+Config.Database.Password+"@/"+Config.Database.Name )
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM City WHERE ID=?" )
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
}

// GetCityByName shows every data where the city name is same (example: becs)
func GetCityByName(c *gin.Context) {

	if Config.Database.MySQL {
		name := c.Params.ByName("name")

		db, err := sql.Open("mysql",Config.Database.Username+":"+Config.Database.Password+"@/"+Config.Database.Name )
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		defer db.Close()
		rows, err := db.Query("SELECT CityName,Latitude,Longitude,Temp,Rain,Date FROM City INNER JOIN CityInfo ON City.ID = CityInfo.CityID where CityName =? ORDER BY Date ", name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// cities container
		var cities cityStructs.CitiesInfo

		for rows.Next() {
			var CityName string
			var Latitude float64
			var Longitude float64
			var Temp string
			var Rain string
			var Date int64

			rows.Scan(&CityName, &Latitude, &Longitude, &Temp, &Rain, &Date)

			RainData := StringToFloatArray(Rain)
			TempData := StringToFloatArray(Temp)

			cities = append(cities, cityStructs.CityInfo{CityName, cityStructs.Geo{Latitude, Longitude}, TempData, RainData, Date})
		}
		if len(cities) == 0{
			content := gin.H{"error": "city with name " + name + " not found"}
			c.JSON(404, content)
			return
		}
		c.JSON(200, cities)
		return
	}else {
		cities := CityDatabase
		// find city's name from url
		name := c.Params.ByName("name")
		// bool for checking city is exist in our db

		filteredCitiesByName := FilterCitiesByName(cities, name)

		if len(filteredCitiesByName) == 0 {
			// response when city doesnt exist in our db
			content := gin.H{"error": "city with name " + name + " not found"}
			c.JSON(404, content)
			return
		}

		// sorting cities
		sort.Sort(filteredCitiesByName)
		// response when city exist in our db
		c.JSON(200, gin.H{"filteredCitiesByTime": filteredCitiesByName})
		return
	}
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

	dataDoesntExistsMessage := CheckDataExist(lat, lng, timestamp)

	if len(dataDoesntExistsMessage) > 0 {
		content := gin.H{"error_message": dataDoesntExistsMessage}
		c.JSON(400, content)
		return
	}

	//Convert to float64/int
	latitudeConvertToFloat64, _ := strconv.ParseFloat(strings.TrimSpace(lat), 64)
	longitudeConvertToFloat64, _ := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	timestampConvertToInt, _ := strconv.ParseInt(timestamp, 10, 64)

	convertProblem := CheckConverting(lat, latitudeConvertToFloat64, lng, longitudeConvertToFloat64, timestamp, timestampConvertToInt)

	if len(convertProblem) > 0 {
		content := gin.H{"error_message ": convertProblem}
		c.JSON(400, content)
		return
	}

	if timestampConvertToInt <= 0 {
		content := gin.H{"error_message ": "timestamp should be bigger than 0"}
		c.JSON(400, content)
		return
	}

	// data from the URL
	presentData := cityStructs.CoordinateAndTime{latitudeConvertToFloat64, longitudeConvertToFloat64, timestampConvertToInt}

	var citiesDistance map[string]float64
	var filteredCitiesByTimeForCalculate map[string]cityStructs.CityInfo

	// Using sql database (depends on cingig file)
	if Config.Database.MySQL {
		CityDatabase = CitiesFromSQL()
		CityDatabase2 = NearestCities(c,timestampConvertToInt)
		fmt.Println(CityDatabase2)
	}




	// filter for the nearest data (by timestamp)
	filteredCitiesByTimeForCalculate = NearestCityDataInTime(CityDatabase, timestampConvertToInt)

	// count all distance with channels
	citiesDistance = DistanceCounter(Config.ProcessorNumber, presentData, filteredCitiesByTimeForCalculate)
	// balanced the distances
	var balancedDistanceByLinearInterpolation map[string]float64

	if Config.Calculation.BalancedByDistance {
		balancedDistanceByLinearInterpolation = BalancedDistanceByLinearInterpolation(citiesDistance)
	} else {
		balancedDistanceByLinearInterpolation = citiesDistance
	}

	// counting temps and raining data for next 5 days
	wg.Add(2)
	var forecastRain []float64
	var forecastCelsius []float64
	go CalculateRain(balancedDistanceByLinearInterpolation, filteredCitiesByTimeForCalculate, &forecastRain, &wg)
	go CalculateTemp(balancedDistanceByLinearInterpolation, filteredCitiesByTimeForCalculate, &forecastCelsius, &wg)
	wg.Wait()

	// send data
	content := gin.H{"expected celsius next 5 days": forecastCelsius, "expected rainning chance next 5 days": forecastRain}
	c.JSON(200, content)
}

// BalancedDistanceByLinearInterpolation ponderare by linear interpolation (nearest 1 weight, furthest 0)
func BalancedDistanceByLinearInterpolation(distances map[string]float64) (balanceByDistance map[string]float64) {

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
		balanceNumber := (v - smallest) / (biggest - smallest)
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


// NearestCityDataInTime is a filter where we get back just one city (exm if we have 3 becs back just one) which is the most relevant by time
func NearestCityDataInTime(allCities []cityStructs.CityInfo, timestamp int64) (filteredCities map[string]cityStructs.CityInfo) {

	citiesDistance := make(map[string]cityStructs.CityInfo)
	makeDifferenceBetweenCities := 0

	for _, v := range allCities {

		if Config.Calculation.FilteringCityData {
			oldDataCityDistanceTime := citiesDistance[v.City].Timestamp - timestamp
			if oldDataCityDistanceTime < 0 {
				oldDataCityDistanceTime *= -1
			}

			newDataCityDistanceTime := v.Timestamp - timestamp
			if newDataCityDistanceTime < 0 {
				newDataCityDistanceTime *= -1
			}

			if oldDataCityDistanceTime > newDataCityDistanceTime || citiesDistance[v.City].Timestamp == 0 {
				citiesDistance[v.City] = v
			}

		} else {
			makeDifferenceBetweenCities++
			cityID := strconv.Itoa(makeDifferenceBetweenCities)
			citiesDistance[cityID] = v
		}
	}
	return citiesDistance
}

//CitiesFromSQL make a query to database for cities
func CitiesFromSQL() (cities []cityStructs.CityInfo) {

		// open SQL
		db, err := sql.Open("mysql",Config.Database.Username+":"+Config.Database.Password+"@/"+Config.Database.Name )
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		defer db.Close()
		rows, err := db.Query("SELECT CityName,Latitude,Longitude,Temp,Rain,Date FROM City INNER JOIN CityInfo ON City.ID = CityInfo.CityID ")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// cities container
		for rows.Next() {
			var CityName string
			var Latitude float64
			var Longitude float64
			var Temp string
			var Rain string
			var Date int64

			rows.Scan(&CityName, &Latitude, &Longitude, &Temp, &Rain, &Date)

			RainData := StringToFloatArray(Rain)
			TempData := StringToFloatArray(Temp)

			cities = append(cities, cityStructs.CityInfo{CityName, cityStructs.Geo{Latitude, Longitude}, TempData, RainData, Date})
		}

	return cities
}

// DistanceCounter where we count every city's distance from an exact place
func DistanceCounter(ProcessorNumber int, coordinate cityStructs.CoordinateAndTime, filteredCities map[string]cityStructs.CityInfo) (distanceCities map[string]float64) {

	var wg sync.WaitGroup

	// because of the append we need to declare here by make
	result := make(map[string]float64)

	//make channel
	in := make(chan cityStructs.CityInfo, len(filteredCities))
	out := make(chan cityStructs.Out, len(filteredCities))

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
func DistanceCounterProcess(in chan cityStructs.CityInfo, coordinate cityStructs.CoordinateAndTime, out chan cityStructs.Out, wg *sync.WaitGroup) {

	for {
		select {
		case cityInfo := <-in:
			defer wg.Done()
			// count distance
			latitudeDistance := coordinate.Lat - cityInfo.Geo.Lat
			longitudeDistance := coordinate.Lng - cityInfo.Geo.Lng
			distance := math.Sqrt(math.Pow(latitudeDistance, 2) + math.Pow(longitudeDistance, 2))

			//response data Out type and make map just at other side because with map can be gorutine problems (see at type Out)
			res := cityStructs.Out{cityInfo.City, distance}
			//send back
			out <- res
		}
	}
}

func CheckDataExist(lat string, lng string, timestamp string) (dataDoesntExistsMessage string) {

	if lat == "" {
		dataDoesntExistsMessage += "lat data must be exists, "
	}
	if lng == "" {
		dataDoesntExistsMessage += "lng data must be exists, "
	}
	if timestamp == "" {
		dataDoesntExistsMessage += "timestamp data must be exists"
	}
	return dataDoesntExistsMessage
}

func CheckConverting(lat string, latitudeConvertToFloat64 float64, lng string, longitudeConvertToFloat64 float64, timestamp string, timestampConvertToInt int64) (convertProblem string) {

	if lat != "0" && latitudeConvertToFloat64 == 0 {
		convertProblem += "invalid lat data (not number), "
	}

	if lng != "0" && longitudeConvertToFloat64 == 0 {
		convertProblem += "invalid lng data (not number), "
	}

	if timestamp != "0" && timestampConvertToInt == 0 {
		convertProblem += "invalid timestamp data (not number) "
	}
	return convertProblem
}


func StringToFloatArray(string string)(result [5]float64){

	// Rain and Temp data is in a string first split up,
	SplitString := strings.SplitN(string, ",", 5)

	var stringToFloat = [5]float64{}

	// convert Temp data to float and put into array
	for i, v := range SplitString {
		f, _ := strconv.ParseFloat(v, 64)
		stringToFloat[i] = f
	}

	return stringToFloat
}

func FloatArrayConvertToString(jsonData [5]float64)(dataForSQL string){

	for i, v := range  jsonData {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if i != len(jsonData)-1 {
			dataForSQL += s + ","
		} else {
			dataForSQL += s
		}
	}
	return dataForSQL
}

func FilterCitiesByName(cities []cityStructs.CityInfo, name string)(filteredCities CitiesInfo){

	for _, v := range cities {
		if v.City == name {
			filteredCities = append(filteredCities, v)
		}
	}
	return filteredCities
}


func NearestCities(c *gin.Context, timestamp_int int64)(cities []cityStructs.CityInfo){

		db, err := sql.Open("mysql",Config.Database.Username+":"+Config.Database.Password+"@/"+Config.Database.Name )
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		defer db.Close()
		rows1, err := db.Query("select CityId,Max(Date) from CityInfo where Date < 2400  Group by CityId")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		rows2, err := db.Query("select CityId,Min(Date) from CityInfo where Date > 2400  Group by CityId")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		for rows1.Next() {
			var CityId int
			var Date int

			rows1.Scan(&CityId, &Date)
			fmt.Println(CityId, Date)

		}
	fmt.Println("***")
		for rows2.Next() {
			var CityId int
			var Date int

			rows2.Scan(&CityId, &Date)
			fmt.Println(CityId, Date)

		}

		fmt.Println(cities)

		if len(cities) == 0{
			content := gin.H{"error": "city with name not found"}
			c.JSON(404, content)
			return
		}

	return
}

func CreateNewDatabase(name string) *sql.DB {

	db, err := sql.Open("mysql", Config.Database.Username+":"+Config.Database.Password+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_,err = db.Exec("CREATE DATABASE IF NOT EXISTS "+name)
	if err != nil {
		panic(err)
	}

	_,err = db.Exec("USE "+name)
	if err != nil {
		panic(err)
	}

	_,err = db.Exec("CREATE TABLE City ( ID int not null Auto_Increment, CityName varchar(255) not NUll, Primary Key(ID), UNIQUE KEY (CityName))")
	if err != nil {
		panic(err)
	}

	_,err = db.Exec("CREATE TABLE CityInfo ( InfoId int not null Auto_Increment, CityId int not NUll, Date int, Temp varchar(255), Rain varchar(255), Latitude float, Longitude float, Primary Key(InfoId), FOREIGN KEY (CityId) REFERENCES City(ID) on DELETE CASCADE )")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	return db
}