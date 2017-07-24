package mockDatabase

import (
	"github.com/bednayb/Go_cities/cityStructs"
	"time"
)

// Cities contains every city's data
var Cities = []cityStructs.CityInfo{

	cityStructs.CityInfo{"paris", cityStructs.Geo{100, 100}, [5]float64{40, 5, 17, 5, 6}, [5]float64{0.2, 0.3, 0.4, 0.5, 0.6}, 1000},
	cityStructs.CityInfo{"becs", cityStructs.Geo{90, 90}, [5]float64{40, 3, 17, 5, 6}, [5]float64{0.5, 0.3, 0.4, 0.5, 0.6}, 1500},
	cityStructs.CityInfo{"bp", cityStructs.Geo{95, 95}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
	cityStructs.CityInfo{"london", cityStructs.Geo{98, 98}, [5]float64{20, 4, 17, 5, 6}, [5]float64{0.2, 0.3, 0.4, 0.5, 0.6}, 1000},
	cityStructs.CityInfo{"becs", cityStructs.Geo{96, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.5, 0.3, 0.4, 0.5, 0.6}, 1000},
	cityStructs.CityInfo{"111london", cityStructs.Geo{95, 95}, [5]float64{50, 3, 17, 5, 6}, [5]float64{0.5, 0.3, 0.4, 0.5, 0.6}, time.Now().Unix()},
}
