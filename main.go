package main

import (
	"github.com/go-martini/martini"
	"time"

)


type CityInfo struct {
	City string
	Geo map[string]float64
	Temp [5]float64
	Rain [5]float64
    Timestamp time.Time
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