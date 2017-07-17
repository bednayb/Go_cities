package citiesFunction

import (
	"github.com/bednayb/Go_cities/cityStructs"
	"testing"
	"reflect"
)

func TestCheckDistance(t *testing.T) {

	mockOneCity := make(map[string]cityStructs.CityInfo)
	mockOneCity["bp"] =	cityStructs.CityInfo{"bp", cityStructs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100}
	mockOneCity["becs"] =	cityStructs.CityInfo{"becs", cityStructs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100}
	mockOneCity["sopron"] =	cityStructs.CityInfo{"sopron", cityStructs.Geo{79, 80}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100}

	resultOneCity := Check_distance(cityStructs.CoordinateAndTime{100, 100, 3}, mockOneCity)
	expectedValue := make(map[string]float64)
	expectedValue["bp"] = 5
	expectedValue["becs"] = 10
	expectedValue["sopron"] = 29

	if false == reflect.DeepEqual(resultOneCity, expectedValue) {
		t.Fatal("doesnt equal result:", resultOneCity, "expected_value:", expectedValue)
	}
}

//func NearestCityDataInTime(all_cities []cityStructs.CityInfo, timestamp int64) (filtered_cities map[string]cityStructs.CityInfo)
func TestNearestCityDataInTime(t *testing.T) {

	var timestamp int64 = 1000

	mockFourCity := []cityStructs.CityInfo{
		cityStructs.CityInfo{"paris", cityStructs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
		cityStructs.CityInfo{"paris", cityStructs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500},
		cityStructs.CityInfo{"bp", cityStructs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100},
		cityStructs.CityInfo{"bp", cityStructs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 2100},
	}

	resultCities := NearestCityDataInTime(mockFourCity, timestamp)
	expectedValue := make(map[string]cityStructs.CityInfo)
	expectedValue["bp"] = cityStructs.CityInfo{"bp", cityStructs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	expectedValue["paris"] = cityStructs.CityInfo{"paris", cityStructs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500}

	if false == reflect.DeepEqual(resultCities, expectedValue) {
		t.Fatal("doesnt equal result:", resultCities, "expected_value:", expectedValue)
	}
}

func TestBalancedDistance(t *testing.T) {

	cityDistance := make(map[string]float64)
	cityDistance["bp"] = 0
	cityDistance["becs"] = 15
	cityDistance["sopron"] = 30

	expectedValue := make(map[string]float64)
	expectedValue["bp"] = 1
	expectedValue["becs"] = 0.5
	expectedValue["sopron"] = 0

	resultCities := BalanceDistance(cityDistance)

	if false == reflect.DeepEqual(resultCities, expectedValue) {
		t.Fatal("doesnt equal result:", resultCities, "expected_value:", expectedValue)
	}
}

func TestCalculateTemp(t *testing.T) {

	balance := make(map[string]float64)
	balance["bp"] = 0
	balance["becs"] = 0.5
	balance["sopron"] = 1

	expectedResult := []float64{20, 3, 17, 5, 6}

	cityInfo := make(map[string]cityStructs.CityInfo)

	cityInfo["bp"] = cityStructs.CityInfo{"bp", cityStructs.Geo{94, 92}, [5]float64{100, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	cityInfo["becs"] = cityStructs.CityInfo{"becs", cityStructs.Geo{94, 92}, [5]float64{10, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	cityInfo["sopron"] = cityStructs.CityInfo{"sopron", cityStructs.Geo{94, 92}, [5]float64{25, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}

	resultTemps := CalculateTemps(balance, cityInfo)

	if false == reflect.DeepEqual(resultTemps, expectedResult) {
		t.Fatal("doesnt equal result:", resultTemps, "expected_value:", expectedResult)
	}
}

func TestCalculateRain(t *testing.T) {

	balance := make(map[string]float64)
	balance["bp"] = 0
	balance["becs"] = 0.5
	balance["sopron"] = 1

	expectedResult := []float64{0.2, 0.6, 0.4, 0.5, 0.6}

	cityInfo := make(map[string]cityStructs.CityInfo)

	cityInfo["bp"] = cityStructs.CityInfo{"bp", cityStructs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	cityInfo["becs"] = cityStructs.CityInfo{"becs", cityStructs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	cityInfo["sopron"] = cityStructs.CityInfo{"sopron", cityStructs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}

	resultTemps := CalculateRain(balance, cityInfo)

	if false == reflect.DeepEqual(resultTemps, expectedResult) {
		t.Fatal("doesnt equal result:", resultTemps, "expected_value:", expectedResult)
	}
}
