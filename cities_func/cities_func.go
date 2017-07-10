package cities_func

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"sort"
	"fmt"
	"github.com/bednayb/Go_cities/city_structs"
	"github.com/bednayb/Go_cities/mock_data"
	"math"
	"strconv"
	"strings"
	"sort"
)

// TODO ez nagyon úgy tűnik mintha a mock adatokat adnánk vissza minden esetben mikor a városokat lekérdezzük! (? ezekhez nem kene egy adatbazis)
// TODO A mock adatokkal való tesztelést különítsük el a valós működéstől, csak akkor induljon mock adatokkal a program ha arra kértük (? ezekhez nem kene egy adatbazis )
// TODO live/demo setupoláshoz vagy config file-t használjunk, vagy argumentumokat program indításkor (? ez full kodos :))
// (?)
func GetCities(c *gin.Context) {
	c.JSON(200, mock_data.All_Cities)
}

func GetCityName(c *gin.Context) {
	name := c.Params.ByName("name")

	var filtered_cities_by_time CitiesInfo

	var redflag bool = true
	for _, v := range mock_data.All_Cities {
		if v.City == name {
			redflag = false
			filtered_cities_by_time = append(filtered_cities_by_time, v)
		}
	}
	sort.Sort(filtered_cities_by_time)

	if redflag {
		content := gin.H{"error": "city with name " + name + " not found"}
		c.JSON(404, content)
		return
	} else {
		c.JSON(200, gin.H{"filtered cities by time": filtered_cities_by_time})
		return
	}
	// TODO érdemes lenne mindkét if ágban egy return, hogy ide ne juthassunk el. (ready)
	// Ha itt bármilyen kód lenne független attól hogy not found volt e lefutna!

}
// TODO ennek a fügvénynek a neve nem tükrözi hogy valójában mit csinál  (ready)
func GetExpectedForecast(c *gin.Context) {

	// save data from URL
	lat := c.Query("lat")
	lng := c.Query("lng")
	timestamp := c.Query("timestamp")

	// TODO Hiba ellenőrzéskor értelmes hibaüzenetet szeretnénk adni pontosan arról ami a hibát okozta (ready)
	var message string

	if lat == "" {
		message += "lat data must be exists, "
	}
	if lng == "" {
		message += "lng data must be exists, "
	}
	if timestamp == "" {
		message += "timestamp data must be exists"
	}
	if len(message) > 0 {
		content := gin.H{"error_message ": message}
		c.JSON(400, content)
		// TODO itt érdemes lenne egy return, hogy ne folytassuk a futást ha hiba volt  (ready)
		return
	}

	//Convert to float64/int
	var convert_problem string
	lat_float64, _ := strconv.ParseFloat(strings.TrimSpace(lat), 64)

	if lat != "0" && lat_float64 == 0{
		convert_problem += "invalid lat data (not number), "
	}

	lng_float64, _ := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if lat != "0" && lng_float64 == 0{
		convert_problem += "invalid lng data (not number), "
	}

	timestamp_int, _ := strconv.ParseInt(timestamp, 10, 64)
	if lat != "0" && timestamp_int == 0{
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

	// filter for the nearest data (by timestamp)
	var filtered_cities = Nearest_city_data_in_time(mock_data.All_Cities, timestamp_int)

	var present_data = city_structs.Cordinate_and_time{lat_float64, lng_float64, timestamp_int}

	// count all distances
	var distances []float64 = Check_distance(present_data, filtered_cities)

	// balanced the distances
	var balance []float64 = Balanced_distance(distances)

	// count the forecast data todo refactor to one function
	var forecast_celsius []float64 = Calculate_temps(balance, filtered_cities)
	var forecast_rain []float64 = Calculate_rain(balance, filtered_cities)

	// todo delete them
	fmt.Println(forecast_rain, forecast_celsius)
	// add description
	content := gin.H{"expected celsius next 5 days": forecast_celsius, "expected rainning chance next 5 days": forecast_rain}
	//content:= filtered_cities
	// send data
	c.JSON(200, content)
}

func Check_distance(cordinate city_structs.Cordinate_and_time, info []city_structs.CityInfo) []float64 {
	// container for distance
	var distances []float64

	for _, info := range info {
		//pitágoras
		dis_lat := cordinate.Lat - info.Geo.Lat
		dis_lng := cordinate.Lng - info.Geo.Lng

		var distance float64
		distance = math.Sqrt(math.Pow(dis_lat, 2) + math.Pow(dis_lng, 2))
		distances = append(distances, distance)
	}

	return distances
}

func Balanced_distance(distances []float64) []float64 {

	// container for balanced distance (return value)
	var balance_by_distance []float64
	// append number to return value
	var balance_number float64

	// todo find them in one for cicle
	// find smallest
	var permanent_smallest float64 = distances[0]
	var smallest float64 = distances[0]

	for _, v := range distances {
		if v < permanent_smallest {
			permanent_smallest = v
			smallest = permanent_smallest
		}
	}

	// find biggest
	var permanent_biggest float64
	var biggest float64 = distances[0]

	for _, v := range distances {
		if v > permanent_biggest {
			permanent_biggest = v
			biggest = permanent_biggest
		}
	}

	// calculate balance number
	for _, v := range distances {
		//rate
		balance_number = (v - smallest) / (biggest - smallest)
		// todo find good description for this :)
		balance_number -= 1
		balance_number *= -1
		//add data to cointainer
		balance_by_distance = append(balance_by_distance, balance_number)
	}

	return balance_by_distance
}

// todo refactor calculate_temps and calulate_rain to one function
func Calculate_temps(balance []float64, info []city_structs.CityInfo) []float64 {

	var forecast_celsius []float64
	var total_balance float64
	var total_temp float64

	// todo change to len
	for c := 0; c < 5; c++ {
		total_balance = 0
		total_temp = 0
		for i, v := range info {
			total_balance += balance[i]
			total_temp += v.Temp[c] * balance[i]
		}

		// cut off 2 decimal
		var untruncated float64 = total_temp/total_balance
		truncated := float64(int(untruncated * 100)) / 100

		forecast_celsius = append(forecast_celsius, truncated)
	}
	return forecast_celsius
}

func Calculate_rain(balance []float64, info []city_structs.CityInfo) []float64 {

	var forecast_rain []float64
	var total_balance float64
	var total_temp float64

	// todo change to len
	// c --> city's temps data number
	for c := 0; c < 5; c++ {
		total_balance = 0
		total_temp = 0
		// info --> every city
		for i, v := range info {
			total_balance += balance[i]
			total_temp += v.Rain[c] * balance[i]
		}
		// cut off 2 decimal
		var untruncated float64 = total_temp/total_balance
		truncated := float64(int(untruncated * 100)) / 100

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
	return slice[i].Timestamp < slice[j].Timestamp;
}

func (slice CitiesInfo) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}


func PostCity(c *gin.Context) {

	var json city_structs.CityInfo
	c.Bind(&json) // This will infer what binder to use depending on the content-type header.

	for _, v := range json.Rain {

		if v < 0 || v > 1 {
			c.JSON(400, gin.H{
				"result": "Failed, invalid temp data (should be beetween 0 and 1)",
			})
		}
	}
	mock_data.All_Cities = append(mock_data.All_Cities, json)
	content := gin.H{
		"result": "successful saving",
		//"title":  json,
	}
	c.JSON(201, content)

}
// TODO használjunk visszatérési érték változónevet is. (ready)
func Nearest_city_data_in_time(all_cities []city_structs.CityInfo, timestamp int64) (filtered_cities []city_structs.CityInfo) {
	// TODO én MAP ez használnék ahol a város neve a kulcs  (? miert)
	// és mindenhol az érték felülírása akkor történhet meg ha az infó frissebb.

	cities_distance := make(map[string]city_structs.CityInfo)

	for _, v := range all_cities {

		old_data_city_distance_time :=  cities_distance[v.City].Timestamp -timestamp
		if old_data_city_distance_time < 0 {
			old_data_city_distance_time *= -1
		}
		new_data_city_distance_time := v.Timestamp -timestamp
		if new_data_city_distance_time < 0 {
			new_data_city_distance_time*= -1
		}

		if old_data_city_distance_time > new_data_city_distance_time {
			cities_distance[v.City] = v
		}
	}

	fmt.Println("map data")
	fmt.Println(cities_distance)

	var order_by_time_cites CitiesInfo
	//var filtered_cities city_structs.CitiesInfo // TODO

	for _, v := range all_cities {
		order_by_time_cites = append(order_by_time_cites, v)
	}

	for i, _ := range order_by_time_cites {

		order_by_time_cites[i].Timestamp -= timestamp
		if order_by_time_cites[i].Timestamp < 0 {
			order_by_time_cites[i].Timestamp *= -1
		}
	}

	sort.Sort(order_by_time_cites)

	for i, v := range order_by_time_cites {

		if contains(filtered_cities, v) == false {
			order_by_time_cites[i].Timestamp += timestamp
			filtered_cities = append(filtered_cities, order_by_time_cites[i])
		}
	}

	return filtered_cities
}




func contains(intSlice city_structs.CitiesInfo, searchInt city_structs.CityInfo) bool {
	for _, value := range intSlice {
		if value.City == searchInt.City {
			return true
		}
	}
	return false
}

type CitiesInfo []city_structs.CityInfo

