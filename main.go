package main

import (
	"github.com/go-martini/martini"

	"fmt"
	"math"
)


//TODO put lat and lng to Geo
type CityInfo struct {
	City string `json:"City"`
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
	Temp []float64 `json:"Temp"`
	Rain []float64 `json:"Rain"`
	Timestamp int `json:"Timestamp"`
}



func main() {
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
	m.Post("/push", func() {
		// create something
	})



	m.Run()
}

func check_distance(a float64, b float64 ) float64{
	fmt.Println(math.Sqrt(a * a + b * b))
	return math.Sqrt(a * a + b * b)
}

func find_biggest(array []float64) float64{
	var n, biggest float64

	for _,v:=range array {
		if v>n {
			n = v
			biggest = n
		}
	}
	return biggest
}






