package main

import (
//"strconv"

"github.com/gin-gonic/gin"
"fmt"
	"math"
//"net/http"
	"strconv"
	"strings"
)


//TODO put lat and lng to Geo
type CityInfo struct {
	City string `gorm:"not null" form:"City" json:"City"`
	Lat float64 `gorm:"not null" form:"Lat" json:"Lat"`
	Lng float64 `gorm:"not null" form:"Lng"json:"Lng"`
	//Cordinate []float64
	Temp [5]float64 `gorm:"not null" form:"Temp"json:"Temp"`
	Rain [5]float64 `gorm:"not null" form:"Rain"json:"Rain"`
	Timestamp float64 `gorm:"not null" form:"Timestamp"json:"Timestamp"`
}

type Cordinate_and_time struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
	Timestamp int64 `gorm:"not null" form:"Timestamp"json:"Timestamp"`
}

var All_Cities = []CityInfo{
	CityInfo{"bp",  99 ,99,[5]float64{20,3,17,5,6},[5]float64{0.2,0.6,0.4,0.5,0.6},43},
	CityInfo{"becs",98,98,[5]float64{20,4,17,5,6},[5]float64{0.2,0.3,0.4,0.5,0.6},43},
	CityInfo{"becs",97,97,[5]float64{20,5,17,5,6},[5]float64{0.2,0.3,0.4,0.5,0.6},43},
	CityInfo{"becs",96,96,[5]float64{20,3,17,5,6},[5]float64{0.5,0.3,0.4,0.5,0.6},43},
}


func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{

		v1.POST("/push", PostCity)

		v1.GET("/cities", GetCities)
		v1.GET("/city/:name", GetCityName)
		v1.GET("/avg", GetCordinate)


	}

	r.Run(":8080")
}


func GetCities(c *gin.Context) {

	c.JSON(200, All_Cities)

}

func GetCityName(c *gin.Context) {
	name := c.Params.ByName("name")

	fmt.Println(name)
	var redflag bool = true
	for _,v := range All_Cities{
		if v.City == name{
			redflag = false
			content := gin.H{ "celsius next 5 days": v.Temp, "rainning chance next 5 days": v.Rain}
			c.JSON(200, content)
		}
	}

	if redflag{
		content := gin.H{"error": "city with name " + name + " not found"}
		c.JSON(404, content)
	}

	// curl -i http://localhost:8080/api/v1/cities/bp
}

func PostCity(c *gin.Context) {
	// The futur code…

}

func GetCordinate(c *gin.Context) {

	//Cordinate + time

		lat := c.Query("lat")
		lng := c.Query("lng")
		timestamp := c.Query("timestamp")
	//Convert to float64/int
		lat_float64, _ := strconv.ParseFloat(strings.TrimSpace(lat), 64)
		lng_float64, _ := strconv.ParseFloat(strings.TrimSpace(lng), 64)
		timestamp_int, _:= strconv.ParseInt(timestamp, 10,64)
	//put data to struct
		var present_data = Cordinate_and_time{lat_float64,lng_float64,timestamp_int}

	// count all distances
	var distances []float64 = check_distance(present_data,All_Cities)

	// check balance

	var balance []float64 = check_balance(distances)

	var forecast_celsius []float64  = calculate_temps(balance,All_Cities)
	var forecast_rain []float64  = calculate_rain(balance,All_Cities)

	var result [][]float64

	result = append(result,forecast_celsius)
	result = append(result, forecast_rain)

	content := gin.H{ "expected celsius next 5 days": result[0], "expected rainning chance next 5 days": result[1]}

	c.JSON(200, content)
}

func check_distance(cordinate Cordinate_and_time, info [] CityInfo) []float64{

	var distances []float64

	for _,info := range info {

		dis_lat := cordinate.Lat - info.Lat
		dis_lng := cordinate.Lng - info.Lng

		var distance float64
		distance = math.Sqrt( math.Pow(dis_lat,2) + math.Pow(dis_lng,2))
		distances = append(distances,distance)
	}

	return distances
}

func check_balance(distances []float64) []float64{

	//return value
	var balance_by_distance []float64
	// append number to return value
	var balance_number float64

	// find smallest
	var permanent_smallest float64 = distances[0]
	var smallest float64 = distances[0]

		for _,v:=range distances {
			if v < permanent_smallest {
				permanent_smallest = v
				smallest = permanent_smallest
			}
		}

	// find biggest
	var permanent_biggest float64
	var biggest float64 =  distances[0]

	for _,v:=range distances {
		if v > permanent_biggest {
			permanent_biggest = v
			biggest = permanent_biggest
		}
	}




	// calculate balance number
	for _,v := range distances{

		balance_number = (v - smallest) / (biggest-smallest)

		balance_number -=1
		balance_number *=-1

		balance_by_distance = append(balance_by_distance,balance_number)
	}

	return balance_by_distance
}

func calculate_temps(balance []float64, info [] CityInfo) []float64 {

	var forecast_celsius []float64
	var total_balance float64
	var total_temp float64




	// change to len
	for c:=0; c < 5; c++ {
		total_balance = 0
		total_temp = 0
		for i, v := range info {
			total_balance += balance[i]
			total_temp += v.Temp[c] * balance[i]
		}
		forecast_celsius = append(forecast_celsius, total_temp/total_balance)
	}
	return forecast_celsius
}

func calculate_rain(balance []float64, info [] CityInfo) []float64 {

	var forecast_rain []float64
	var total_balance float64
	var total_temp float64

	// change to len
	for c:=0; c < 5; c++ {
		total_balance = 0
		total_temp = 0
		for i, v := range info {
			total_balance += balance[i]
			total_temp += v.Rain[c] * balance[i]
		}
		forecast_rain = append(forecast_rain, total_temp/total_balance)
	}
	return forecast_rain
}










