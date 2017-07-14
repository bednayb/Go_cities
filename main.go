package main
// TODO a main go a főkönyvtárvban szokott lenni általában (ready)
import (
	"github.com/bednayb/Go_cities/cities_func"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"fmt"
)

func main() {

	// mock data  go run main.go -config development
	// test_db data  go run main.go -config test
	// real data  go run main.go

	var conf string
	cities_func.ConfigSettings(&conf)

	//Set config file path including file name and extension
	viper.SetConfigFile("./config/"+conf+".json")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	if !viper.IsSet(conf+".database") {
		log.Fatal("missing database")
	}

	//Settings data
	cities_func.Init(conf)


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
