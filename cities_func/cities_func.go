package cities_func

import (
	"flag"
	"github.com/bednayb/Go_cities/city_db"
	"github.com/bednayb/Go_cities/city_structs"
	"github.com/bednayb/Go_cities/mock_data"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
	"strconv"
	"strings"
)

// TODO ez nagyon úgy tűnik mintha a mock adatokat adnánk vissza minden esetben mikor a városokat lekérdezzük! (ready)
// TODO A mock adatokkal való tesztelést különítsük el a valós működéstől, csak akkor induljon mock adatokkal a program ha arra kértük (ready, test not works)
// TODO live/demo setupoláshoz vagy config file-t használjunk, vagy argumentumokat program indításkor (? add_config_file_branch))
// Ricsi --> akkor hasznalj mock adatokat ha go run main.go --mock al hivod meg kul, (go run main.go) azzal ami el van mentve
// Zoli -->  add config  https://github.com/spf13/viper

var Database []city_structs.CityInfo

func SelectDatabase() []city_structs.CityInfo {

	var mock = flag.String("mock", "", "placeholder")
	flag.Parse()
	if *mock == "true" {
		Database = mock_data.All_Cities
		return mock_data.All_Cities
	} else {
		Database = city_db.All_Cities
		return city_db.All_Cities
	}
}

func GetAllCity(c *gin.Context) {
	cities := Database
	c.JSON(200, cities)
}

func GetCityByName(c *gin.Context) {

	cities := Database
	// find city's name from url
	name := c.Params.ByName("name")
	// bool for checking city is exist in our db
	var redflag bool = true

	// filtered cities order by timestamp (first the oldest)
	var filteredCitiesByTime CitiesInfo

	// filtering cities by name
	for _, v := range cities {
		if v.City == name {
			redflag = false
			filteredCitiesByTime = append(filteredCitiesByTime, v)
		}
	}
	if redflag {
		// response when city doesnt exist in our db
		content := gin.H{"error": "city with name " + name + " not found"}
		c.JSON(404, content)
		return
	} else {
		// sorting cities
		sort.Sort(filteredCitiesByTime)
		// response when city exist in our db
		c.JSON(200, gin.H{"filtered_cities_by_time": filteredCitiesByTime})
		return
	}
	// TODO érdemes lenne mindkét if ágban egy return, hogy ide ne juthassunk el. (ready)
	// Ha itt bármilyen kód lenne független attól hogy not found volt e lefutna!
}

// TODO ennek a fügvénynek a neve nem tükrözi hogy valójában mit csinál  (ready)
func GetExpectedForecast(c *gin.Context) {

	if len(Database) == 0 {
		content := gin.H{"response": "sry we havnt had enough data for calculating yet"}
		c.JSON(200, content)
		return
	}

	// save data from URL
	lat := c.Query("lat")
	lng := c.Query("lng")
	timestamp := c.Query("timestamp")

	// TODO Hiba ellenőrzéskor értelmes hibaüzenetet szeretnénk adni pontosan arról ami a hibát okozta (ready)

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
		// TODO itt érdemes lenne egy return, hogy ne folytassuk a futást ha hiba volt  (ready)
		return
	}

	//Convert to float64/int
	var convertProblemMessage string
	latitudeFloat64, _ := strconv.ParseFloat(strings.TrimSpace(lat), 64)

	if lat != "0" && latitudeFloat64 == 0 {
		convertProblemMessage += "invalid lat data (not number), "
	}

	longitudeFloat64, _ := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if lat != "0" && longitudeFloat64 == 0 {
		convertProblemMessage += "invalid lng data (not number), "
	}

	timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
	if lat != "0" && timestampInt == 0 {
		convertProblemMessage += "invalid timestamp data (not number) "
	}

	if len(convertProblemMessage) > 0 {
		content := gin.H{"error_message ": convertProblemMessage}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk a futást ha hiba volt  (ready)
		return
	}

	if timestampInt < 0 {
		content := gin.H{"error_message ": "timestamp should be bigger than 0"}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk a futást ha hiba volt  (ready)
		return
	}

	//put data to struct
	// TODO a fenti parsolások mindegyikénél előfordulhat hiba, amit így teljesen figyelmen kívűl hagyunk (ready)
	// TODO a fentabbi ellenőrzési szisztémával adhatunk hibaüzenetet hogy melyikből nem sikerült számot kinyernünk. (ready)  (timestampnel minuszt nem fogadunk el -- Ricsi)
	// + Ellenőrizhető hogy a szám valós tartományban van e.
	// hiba esetén itt se menjünk tovább.

	// data from the URL
	var presentData = city_structs.Cordinate_and_time{latitudeFloat64, longitudeFloat64, timestampInt}

	// filter for the nearest data (by timestamp)
	var filteredCities = Nearest_city_data_in_time(Database, timestampInt)

	// count all distances
	var distances map[string]float64 = CountDistance(presentData, filteredCities)

	// balanced the distances
	var balance map[string]float64 = BalanceDistance(distances)

	// count the forecast data
	var forecastCelsius []float64 = CalculateTemps(balance, filteredCities)
	var forecastRain []float64 = CalculateRain(balance, filteredCities)

	// send data
	content := gin.H{"expected celsius next 5 days": forecastCelsius, "expected rainning chance next 5 days": forecastRain}
	c.JSON(200, content)
}

func CountDistance(cordinate city_structs.Cordinate_and_time, info map[string]city_structs.CityInfo) (city_distance map[string]float64) {

	// container for distance  key --> city name, value --> distance
	var citiesDistance = make(map[string]float64)

	//count every distance of city (pitágoras)
	var distance float64
	for _, info := range info {

		latitudeDistance := cordinate.Lat - info.Geo.Lat
		longitudeDistance := cordinate.Lng - info.Geo.Lng

		distance = math.Sqrt(math.Pow(latitudeDistance, 2) + math.Pow(longitudeDistance, 2))
		citiesDistance[info.City] = distance
	}
	return citiesDistance
}

// linear interpolation (nearest 1 weight, furthest 0)
func BalanceDistance(distances map[string]float64) (balance_by_distance map[string]float64) {

	//  balanced distance
	var balanceNumber float64

	//// find furthest (biggest number)
	var permanentBiggest float64
	var biggest float64 = 0

	for _, v := range distances {
		if v > permanentBiggest {
			permanentBiggest = v
			biggest = permanentBiggest
		}
	}
	//find nearest (smallest number)
	var permanentSmallest float64 = biggest
	var smallest float64 = biggest

	for _, v := range distances {
		if v < permanentSmallest {
			permanentSmallest = v
			smallest = permanentSmallest
		}
	}
	// calculate balanced numbers
	for i, v := range distances {
		balanceNumber = (v - smallest) / (biggest - smallest)
		balanceNumber -= 1
		balanceNumber *= -1
		// overwrite distance with balanced distance
		distances[i] = balanceNumber
	}

	return distances
}

// todo refactor calculate_temps and calulate_rain to one function
func CalculateTemps(balance map[string]float64, cityInfo map[string]city_structs.CityInfo) (forecastCelsius []float64) {

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
		var untruncated float64 = totalTemp / totalBalance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		forecastCelsius = append(forecastCelsius, truncated)
	}
	return forecastCelsius
}

func CalculateRain(balance map[string]float64, city_info map[string]city_structs.CityInfo) (forecastRain []float64) {

	var total_balance float64
	var total_temp float64

	// count next five days
	for day := 0; day < 5; day++ {
		total_balance = 0
		total_temp = 0
		// info --> every city
		for _, v := range city_info {
			total_balance += balance[v.City]
			total_temp += v.Rain[day] * balance[v.City]
		}
		// cut off 2 decimal
		var untruncated float64 = total_temp / total_balance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		forecastRain = append(forecastRain, truncated)
	}
	return forecastRain
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

	var json city_structs.CityInfo
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

	Database = append(Database, json)

	content := gin.H{
		"result": "successful saving",
	}
	c.JSON(201, content)
}

// TODO használjunk visszatérési érték változónevet is. (ready)
func Nearest_city_data_in_time(all_cities []city_structs.CityInfo, timestamp int64) (filtered_Cities map[string]city_structs.CityInfo) {
	// TODO én MAP ez használnék ahol a város neve a kulcs  (ready)
	// és mindenhol az érték felülírása akkor történhet meg ha az infó frissebb.

	filteredCities := make(map[string]city_structs.CityInfo)

	for _, v := range all_cities {

		oldDataCityDistanceTime := filteredCities[v.City].Timestamp - timestamp
		if oldDataCityDistanceTime < 0 {
			oldDataCityDistanceTime *= -1
		}

		newDataCityDistanceTime := v.Timestamp - timestamp
		if newDataCityDistanceTime < 0 {
			newDataCityDistanceTime *= -1
		}

		if oldDataCityDistanceTime > newDataCityDistanceTime {
			filteredCities[v.City] = v
		}
	}
	return filteredCities
}

type CitiesInfo []city_structs.CityInfo
