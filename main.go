package main

import (
//"strconv"

"github.com/gin-gonic/gin"
"fmt"
	"math"
//"net/http"
)


//TODO put lat and lng to Geo
type CityInfo struct {
	City string `gorm:"not null" form:"City" json:"City"`
	Lat float64 `gorm:"not null" form:"Lat" json:"Lat"`
	Lng float64 `gorm:"not null" form:"Lng"json:"Lng"`
	//Cordinate []float64
	Temp [5]float64 `gorm:"not null" form:"Temp"json:"Temp"`
	Rain [5]float64 `gorm:"not null" form:"Rain"json:"Rain"`
	Timestamp int `gorm:"not null" form:"Timestamp"json:"Timestamp"`
}

type Cordinate struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}

var All_Cities = []CityInfo{
	CityInfo{"bp",  2 ,3,[5]float64{20,3,17,5,6},[5]float64{2,3,4,5,6},43},
	CityInfo{"becs",44,5,[5]float64{20,3,17,5,6},[5]float64{2,3,4,5,6},43},
}


func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{

		v1.POST("/push", PostCity)
		v1.GET("/cities", GetCities)
		v1.GET("/cities/:name", GetCityName)

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
			content := gin.H{ "city": v.City, "lat": v.Lat}
			c.JSON(200, content)
		}
	}

	if redflag{
		content := gin.H{"error": "user with id#" + name + " not found"}
		c.JSON(404, content)
	}

	// curl -i http://localhost:8080/api/v1/cities/bp
}

func PostCity(c *gin.Context) {
	// The futur code…

}

func GetCordinate(c *gin.Context) {
	lat := c.Query("lat")
	lng:= c.Query("let")
	timestamp:= c.Query("timestamp")

	content := gin.H{  "lat": lat, "lng": lng, "timestamp":timestamp}
	c.JSON(200, content)
}







func check_distance(cordinate Cordinate, info [] CityInfo) float64{


	//var total_temp float64
	//var total_balance float64
	//var balanced_temp float64
	//var all_temp []float64


	var avg_temp float64 = 1
    distances:= []float64{}

	for _,info := range info {

		dis_lat := cordinate.Lat - info.Lat
		dis_lng := cordinate.Lng - info.Lng

		var distance float64
		distance = math.Sqrt( math.Pow(dis_lat,2) + math.Pow(dis_lng,2))
		distances = append(distances,distance)

	}

	//var balanced_numbers []float64
	//balanced_numbers = balance(distances)
	//
	//
	//fmt.Println(balanced_numbers)
	//for j:=0; j < len(info); j++ {
	//
	//	total_temp  = 0
	//	total_balance = 0
	//	for i:= 0; i < 5; i++ {
	//
	//		fmt.Println(info[j].Temp[i])
	//		balanced_temp = balanced_numbers[j] * info[i].Temp[j]
	//		total_temp += balanced_temp
	//		total_balance += balanced_numbers[j]
	//	}
	//	avg_temp = total_temp / total_balance
	//	fmt.Println("avg")
	//	fmt.Println(avg_temp)
	//	all_temp = append(all_temp, avg_temp)
	//}
	//fmt.Println(all_temp)
	////fmt.Println(total_temp)
	////fmt.Println(total_balance)
	//
	//avg_temp = total_temp / total_balance
	//fmt.Println(result)
	return avg_temp
}


func find_biggest_and_smallest(array []float64) (a[2]float64){
	var n,k, smallest,biggest float64

	n = array[0]
	for _,v:=range array {
		if v < n {
			n = v
			smallest = n
		}
	}

	for _,v:=range array {
		if v > k {
			k = v
			biggest = k
		}
	}

	var result[2]float64
	result[0] = biggest
	result[1] = smallest
	return result
}

func balance( a []float64) (result []float64){

	fmt.Println(a)

	var big_and_small[2] float64
	big_and_small = find_biggest_and_smallest(a)


	for _,v := range a {

		var balanced_number float64
		balanced_number = v-big_and_small[1]
		balanced_number /= big_and_small[0]

		result = append(result,balanced_number)

	}

	fmt.Println(result)




	return result
}










