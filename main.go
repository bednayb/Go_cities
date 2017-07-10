package main
// TODO a main go a főkönyvtárvban szokott lenni általában (ready)
import (
	"github.com/bednayb/Go_cities/cities_func"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	v1 := r.Group("/")
	{
		// list all cities
		v1.GET("/cities", cities_func.GetCities)
		// find specific city by name
		v1.GET("/city/:name", cities_func.GetCityName)
		// make forecast for exact place
		v1.GET("/avg", cities_func.GetExpectedForecast)
		// add new city
		v1.POST("/push", cities_func.PostCity)
	}

	r.Run(":8080")
}
