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
// TODO live/demo setupoláshoz vagy config file-t használjunk, vagy argumentumokat program indításkor (? ez full kodos :))
// Ricsi --> akkor hasznalj mock adatokat ha go run main.go --mock al hivod meg kul, (go run main.go) azzal ami el van mentve
// Zoli -->  add config  https://github.com/spf13/viper
var City_Database []city_structs.CityInfo

func ConfigSettings(config_file *string) {

	var config = flag.String("config", "", "placeholder")
	flag.Parse()
	if *config == "development" {
		*config_file = "development"
	} else if *config == "test" {
		*config_file = "test"
	} else {
		*config_file = "production"
	}
}

func Init(conf string) {
	if conf == "development" {
		City_Database = mock_data.All_Cities
	} else if conf == "test" {
		City_Database = mock_data.All_Cities
	} else {
		City_Database = city_db.All_Cities
	}
}

func GetCities(c *gin.Context) {
	cities := City_Database
	c.JSON(200, cities)
}

func GetCityName(c *gin.Context) {

	cities := City_Database
	// find city's name from url
	name := c.Params.ByName("name")
	// bool for checking city is exist in our db
	var redflag bool = true

	// filtered cities order by timestamp (first the oldest)
	var filtered_cities_by_time CitiesInfo

	// filtering cities by name
	for _, v := range cities {
		if v.City == name {
			redflag = false
			filtered_cities_by_time = append(filtered_cities_by_time, v)
		}
	}
	if redflag {
		// response when city doesnt exist in our db
		content := gin.H{"error": "city with name " + name + " not found"}
		c.JSON(404, content)
		return
	} else {
		// sorting cities
		sort.Sort(filtered_cities_by_time)
		// response when city exist in our db
		c.JSON(200, gin.H{"filtered_cities_by_time": filtered_cities_by_time})
		return
	}
	// TODO érdemes lenne mindkét if ágban egy return, hogy ide ne juthassunk el. (ready)
	// Ha itt bármilyen kód lenne független attól hogy not found volt e lefutna!
}

// TODO ennek a fügvénynek a neve nem tükrözi hogy valójában mit csinál  (ready)
func GetExpectedForecast(c *gin.Context) {

	if len(City_Database) == 0 {
		content := gin.H{"response": "sry we havnt had enough data for calculating yet"}
		c.JSON(200, content)
		return
	}

	// save data from URL
	lat := c.Query("lat")
	lng := c.Query("lng")
	timestamp := c.Query("timestamp")

	// TODO Hiba ellenőrzéskor értelmes hibaüzenetet szeretnénk adni pontosan arról ami a hibát okozta (ready)

	var data_doenst_exists_message string

	if lat == "" {
		data_doenst_exists_message += "lat data must be exists, "
	}
	if lng == "" {
		data_doenst_exists_message += "lng data must be exists, "
	}
	if timestamp == "" {
		data_doenst_exists_message += "timestamp data must be exists"
	}
	if len(data_doenst_exists_message) > 0 {
		content := gin.H{"error_message": data_doenst_exists_message}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk a futást ha hiba volt  (ready)
		return
	}

	//Convert to float64/int
	var convert_problem string
	lat_float64, _ := strconv.ParseFloat(strings.TrimSpace(lat), 64)

	if lat != "0" && lat_float64 == 0 {
		convert_problem += "invalid lat data (not number), "
	}

	lng_float64, _ := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if lat != "0" && lng_float64 == 0 {
		convert_problem += "invalid lng data (not number), "
	}

	timestamp_int, _ := strconv.ParseInt(timestamp, 10, 64)
	if lat != "0" && timestamp_int == 0 {
		convert_problem += "invalid timestamp data (not number) "
	}

	if len(convert_problem) > 0 {
		content := gin.H{"error_message ": convert_problem}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk a futást ha hiba volt  (ready)
		return
	}

	if timestamp_int < 0 {
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
	var present_data = city_structs.Cordinate_and_time{lat_float64, lng_float64, timestamp_int}

	// filter for the nearest data (by timestamp)
	var filtered_cities = Nearest_city_data_in_time(City_Database, timestamp_int)

	// count all distances
	var distances map[string]float64 = Check_distance(present_data, filtered_cities)

	// balanced the distances
	var balance map[string]float64 = Balanced_distance(distances)

	// count the forecast data
	var forecast_celsius []float64 = Calculate_temps(balance, filtered_cities)
	var forecast_rain []float64 = Calculate_rain(balance, filtered_cities)

	// send data
	content := gin.H{"expected celsius next 5 days": forecast_celsius, "expected rainning chance next 5 days": forecast_rain}
	c.JSON(200, content)
}

func Check_distance(cordinate city_structs.Cordinate_and_time, info map[string]city_structs.CityInfo) (city_distance map[string]float64) {

	// container for distance  key --> city name, value --> distance
	var cities_distance = make(map[string]float64)

	//count every distance of city (pitágoras)
	var distance float64
	for _, info := range info {

		dis_lat := cordinate.Lat - info.Geo.Lat
		dis_lng := cordinate.Lng - info.Geo.Lng

		distance = math.Sqrt(math.Pow(dis_lat, 2) + math.Pow(dis_lng, 2))
		cities_distance[info.City] = distance
	}
	return cities_distance
}

// linear interpolation (nearest 1 weight, furthest 0)
func Balanced_distance(distances map[string]float64) (balance_by_distance map[string]float64) {

	//  balanced distance
	var balance_number float64

	//// find furthest (biggest number)
	var permanent_biggest float64
	var biggest float64 = 0

	for _, v := range distances {
		if v > permanent_biggest {
			permanent_biggest = v
			biggest = permanent_biggest
		}
	}
	//find nearest (smallest number)
	var permanent_smallest float64 = biggest
	var smallest float64 = biggest

	for _, v := range distances {
		if v < permanent_smallest {
			permanent_smallest = v
			smallest = permanent_smallest
		}
	}
	// calculate balanced numbers
	for i, v := range distances {
		balance_number = (v - smallest) / (biggest - smallest)
		balance_number -= 1
		balance_number *= -1
		// overwrite distance with balanced distance
		distances[i] = balance_number
	}

	return distances
}

// todo refactor calculate_temps and calulate_rain to one function
func Calculate_temps(balance map[string]float64, city_info map[string]city_structs.CityInfo) (forecast_temp []float64) {

	// container for temps
	var forecast_celsius []float64

	var total_balance float64
	var total_temp float64

	// count next five days
	for day := 0; day < 5; day++ {
		total_balance = 0
		total_temp = 0
		// info --> every city
		for _, v := range city_info {
			total_balance += balance[v.City]
			total_temp += v.Temp[day] * balance[v.City]
		}
		// cut off 2 decimal
		var untruncated float64 = total_temp / total_balance
		truncated := float64(int(untruncated*100)) / 100
		// put data to container
		forecast_celsius = append(forecast_celsius, truncated)
	}

	return forecast_celsius
}

func Calculate_rain(balance map[string]float64, city_info map[string]city_structs.CityInfo) (forecast_temp []float64) {

	// container for temps
	var forecast_rain []float64

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
		forecast_rain = append(forecast_rain, truncated)
	}
	return forecast_rain
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

	City_Database = append(City_Database, json)

	content := gin.H{
		"result": "successful saving",
	}
	c.JSON(201, content)
}

// TODO használjunk visszatérési érték változónevet is. (ready)
func Nearest_city_data_in_time(all_cities []city_structs.CityInfo, timestamp int64) (filtered_cities map[string]city_structs.CityInfo) {
	// TODO én MAP ez használnék ahol a város neve a kulcs  (ready)
	// és mindenhol az érték felülírása akkor történhet meg ha az infó frissebb.

	cities_distance := make(map[string]city_structs.CityInfo)

	for _, v := range all_cities {

		old_data_city_distance_time := cities_distance[v.City].Timestamp - timestamp
		if old_data_city_distance_time < 0 {
			old_data_city_distance_time *= -1
		}

		new_data_city_distance_time := v.Timestamp - timestamp
		if new_data_city_distance_time < 0 {
			new_data_city_distance_time *= -1
		}

		if old_data_city_distance_time > new_data_city_distance_time || filtered_cities[v.City].Timestamp == 0{
			cities_distance[v.City] = v
		}
	}
	return cities_distance
}

type CitiesInfo []city_structs.CityInfo
