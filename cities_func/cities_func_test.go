package cities_func

import (
	"github.com/bednayb/Go_cities/city_structs"
	"testing"
	"reflect"
)

func TestCheck_distance(t *testing.T) {

	mockOneCity := make(map[string]city_structs.CityInfo)
	mockOneCity["bp"] =	city_structs.CityInfo{"bp", city_structs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100}
	mockOneCity["becs"] =	city_structs.CityInfo{"becs", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100}
	mockOneCity["sopron"] =	city_structs.CityInfo{"sopron", city_structs.Geo{79, 80}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100}

	resultOneCity := Check_distance(city_structs.Cordinate_and_time{100, 100, 3}, mockOneCity)
	expectedValue := make(map[string]float64)
	expectedValue["bp"] = 5
	expectedValue["becs"] = 10
	expectedValue["sopron"] = 29

	if false == reflect.DeepEqual(resultOneCity, expectedValue) {
		t.Fatal("doesnt equal result:", resultOneCity, "expected_value:", expectedValue)
	}
}

//func Nearest_city_data_in_time(all_cities []city_structs.CityInfo, timestamp int64) (filtered_cities map[string]city_structs.CityInfo)
func TestNearest_city_data_in_time(t *testing.T) {

	var timestamp int64 = 1000

	mockFourCity := []city_structs.CityInfo{
		city_structs.CityInfo{"paris", city_structs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
		city_structs.CityInfo{"paris", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500},
		city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100},
		city_structs.CityInfo{"bp", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 2100},
	}

	resultCities := Nearest_city_data_in_time(mockFourCity, timestamp)
	expectedValue := make(map[string]city_structs.CityInfo)
	expectedValue["bp"] = city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	expectedValue["paris"] = city_structs.CityInfo{"paris", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500}

	if false == reflect.DeepEqual(resultCities, expectedValue) {
		t.Fatal("doesnt equal result:", resultCities, "expected_value:", expectedValue)
	}
}

func TestBalanced_distance(t *testing.T) {

	cityDistance := make(map[string]float64)
	cityDistance["bp"] = 0
	cityDistance["becs"] = 15
	cityDistance["sopron"] = 30

	expectedValue := make(map[string]float64)
	expectedValue["bp"] = 1
	expectedValue["becs"] = 0.5
	expectedValue["sopron"] = 0

	resultCities := Balanced_distance(cityDistance)

	if false == reflect.DeepEqual(resultCities, expectedValue) {
		t.Fatal("doesnt equal result:", resultCities, "expected_value:", expectedValue)
	}
}

func TestCalculate_temps(t *testing.T) {

	balance := make(map[string]float64)
	balance["bp"] = 0
	balance["becs"] = 0.5
	balance["sopron"] = 1

	expected_result := []float64{20, 3, 17, 5, 6}

	city_info := make(map[string]city_structs.CityInfo)

	city_info["bp"] = city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{100, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	city_info["becs"] = city_structs.CityInfo{"becs", city_structs.Geo{94, 92}, [5]float64{10, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	city_info["sopron"] = city_structs.CityInfo{"sopron", city_structs.Geo{94, 92}, [5]float64{25, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}

	resultTemps := Calculate_temps(balance, city_info)

	if false == reflect.DeepEqual(resultTemps, expected_result) {
		t.Fatal("doesnt equal result:", resultTemps, "expected_value:", expected_result)
	}
}

func TestCalculate_rains(t *testing.T) {

	balance := make(map[string]float64)
	balance["bp"] = 0
	balance["becs"] = 0.5
	balance["sopron"] = 1

	expected_result := []float64{0.2, 0.6, 0.4, 0.5, 0.6}

	city_info := make(map[string]city_structs.CityInfo)

	city_info["bp"] = city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	city_info["becs"] = city_structs.CityInfo{"becs", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}
	city_info["sopron"] = city_structs.CityInfo{"sopron", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100}

	resultTemps := Calculate_rain(balance, city_info)

	if false == reflect.DeepEqual(resultTemps, expected_result) {
		t.Fatal("doesnt equal result:", resultTemps, "expected_value:", expected_result)
	}
}
