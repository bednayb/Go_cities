package main

import (
	"github.com/go-martini/martini"

	//"fmt"
	//"math"
	//"net/http"

	//"github.com/martini-contrib/render"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"math"
)


//TODO put lat and lng to Geo
type CityInfo struct {
	City string `json:"City"`
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
	Temp [5]float64 `json:"Temp"`
	Rain [5]float64 `json:"Rain"`
	Timestamp int `json:"Timestamp"`
}

type Cordinate struct {
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
}



func main() {

	///////////////// DATA /////////////////
	a:= CityInfo{"bp",4,5,[5]float64{20,3,17,5,6},[5]float64{2,3,4,5,6},43}
	b:= CityInfo{"bp",4,5,[5]float64{20,3,4,5,6},[5]float64{2,3,4,5,6},43}
	//c:= CityInfo{"bp",2,1,[5]float64{20,3,4,5,6},[5]float64{2,3,4,5,6},43}
	//d:= CityInfo{"bp",10,1,[5]float64{20,3,4,5,6},[5]float64{2,3,4,5,6},43}
	//e:= CityInfo{"bp",10,1,[5]float64{20,3,4,5,6},[5]float64{2,3,4,5,6},43}


	present_cord:= Cordinate{1,1}
	///////////////// DATA /////////////////


    distances:= check_distance(present_cord,[]CityInfo{a,b})
	fmt.Println(distances)





	m := martini.Classic()

	// Here you can check a specific city  (not works yet)
	m.Get("/list", func() string {
		return " specific city data"
	})

	// Here you can check the forecast for a specific place (not works yet)
	m.Get("/avg", func() string {
		return "forecast"
	})

	// Here you can send your new data
	m.Post("/push", func(r render.Render, city CityInfo) {
		// create something
		var retData struct{
			City CityInfo
		}

		retData.City = city
		r.JSON(http.StatusOK, city)
	})

	m.Run()
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










