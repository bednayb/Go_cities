package main

import (
	"github.com/bednayb/Go_cities/citiesFunction"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func main() {


	var configFile string
	citiesFunction.ConfigSettings(&configFile)

	//Set config file path including file name and extension
	viper.SetConfigFile("./config/"+ configFile +".json")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	//Settings data
	citiesFunction.Init(configFile)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)

		// Find and read the config file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		// Confirm which config file is used
		fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

		//Settings data
		citiesFunction.Init(configFile)
	})

	r := gin.Default()
	v1 := r.Group("/")
	{
		// list all cities
		v1.GET("/v1/cities", citiesFunction.GetAllCity)
		// make forecast for exact place
		v1.GET("/v1/avg", citiesFunction.GetExpectedForecast)
		// find specific city by name
		v1.GET("/v1/city/:name", citiesFunction.GetCityByName)
		// add new city
		v1.POST("/v1/push", citiesFunction.PostCity)

		// add new city SQL
		v1.DELETE("/sql/delete", citiesFunction.DeleteCitySQL)
	}
	r.Run(":"+citiesFunction.Config.Port)
}
