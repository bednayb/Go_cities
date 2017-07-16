package main
// TODO a main go a főkönyvtárvban szokott lenni általában (ready)
import (
	"github.com/bednayb/Go_cities/citiesFunction"
	"github.com/gin-gonic/gin"
)

func main() {

	// check use mocking db or not
	// if you want to use mock data  run your program with this codeline:
	//run go run main.go -mock true
	//else
	// go run main.go
	citiesFunction.SelectDatabase()

	r := gin.Default()
	v1 := r.Group("/")
	{
		// list all cities
		v1.GET("/cities", citiesFunction.GetAllCity)
		// find specific city by name
		v1.GET("/city/:name", citiesFunction.GetCityByName)
		// make forecast for exact place
		v1.GET("/avg", citiesFunction.GetExpectedForecast)
		// add new city
		v1.POST("/push", citiesFunction.PostCity)
	}

	r.Run(":8080")
}
