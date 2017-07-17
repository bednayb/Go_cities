package main
// TODO a main go a főkönyvtárvban szokott lenni általában (ready)
import (
	"github.com/bednayb/Go_cities/citiesFunction"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"fmt"
)

func main() {

	// if you want to use mock data  run go run main.go -config development (or just go run main.go)
	// if you want to use test data  run go run main.go -config test
	// if you want to use production data  run go run main.go -config production

	var conf string
	citiesFunction.ConfigSettings(&conf)

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
	citiesFunction.Init(conf)

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
