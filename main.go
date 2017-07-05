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
	m.Get("/", func() string {
		return "Hello world!"
	})


	m.Run()
}