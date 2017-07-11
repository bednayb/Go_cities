package mock_data

import (
	"github.com/bednayb/Go_cities/city_structs"
	//"time"
	//"time"
)

var All_Cities = []city_structs.CityInfo{


	city_structs.CityInfo{"paris", city_structs.Geo{100, 100}, [5]float64{40, 5, 17, 5, 6}, [5]float64{0.2, 0.3, 0.4, 0.5, 0.6}, 1000},
	city_structs.CityInfo{"becs", city_structs.Geo{90, 90}, [5]float64{40, 3, 17, 5, 6}, [5]float64{0.5, 0.3, 0.4, 0.5, 0.6}, 1500},

	city_structs.CityInfo{"bp", city_structs.Geo{95, 95}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
	//city_structs.CityInfo{"london", city_structs.Geo{98, 98}, [5]float64{20, 4, 17, 5, 6}, [5]float64{0.2, 0.3, 0.4, 0.5, 0.6}, 1000},
	//city_structs.CityInfo{"becs", city_structs.Geo{96, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.5, 0.3, 0.4, 0.5, 0.6}, 1000},
	//city_structs.CityInfo{"london", city_structs.Geo{95, 95}, [5]float64{50, 3, 17, 5, 6}, [5]float64{0.5, 0.3, 0.4, 0.5, 0.6}, time.Now().Unix()},
}
