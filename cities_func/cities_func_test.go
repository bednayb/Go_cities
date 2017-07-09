package cities_func

import (
	"github.com/bednayb/Go_cities/city_structs"
	"testing"
)

func TestCheck_distance(t *testing.T) {
	// todo not dry make one outside
	mockOneCity := []city_structs.CityInfo{
		city_structs.CityInfo{"bp", city_structs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
	}
	mockThreeCity := []city_structs.CityInfo{

		city_structs.CityInfo{"bp", city_structs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
		city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
		city_structs.CityInfo{"bp", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
	}

	resultOneCity := Check_distance(city_structs.Cordinate_and_time{100, 100, 3}, mockOneCity)
	expectedValue := []float64{5}

	// works but dont prefer !true
	if false == CompareSlices_General(resultOneCity, expectedValue) {
		t.Fatal("doesnt equal result:", resultOneCity, "expected_value:", expectedValue)
	}

	resultThreeCity := Check_distance(city_structs.Cordinate_and_time{100, 100, 3}, mockThreeCity)
	expectedThreeCityValue := []float64{5, 10, 29}

	if false == CompareSlices_General(resultThreeCity, expectedThreeCityValue) {
		t.Fatal("doesnt equal result:", resultThreeCity, "expected_value:", expectedThreeCityValue)
	}
}

func TestNearest_city_data_in_time(t *testing.T) {
	// todo not dry make one outside

	mockFourCity := []city_structs.CityInfo{

		city_structs.CityInfo{"paris", city_structs.Geo{97, 96}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
		city_structs.CityInfo{"paris", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500},
		city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100},
		city_structs.CityInfo{"bp", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 2100},
	}

	resultCities := Nearest_city_data_in_time(mockFourCity, 1000)
	expectedValue := []city_structs.CityInfo{
		city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100},
		city_structs.CityInfo{"paris", city_structs.Geo{80, 79}, [5]float64{20, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500},
	}

	if false == CompareSlices_General_city(resultCities, expectedValue) {
		t.Fatal("doesnt equal result:", resultCities, "expected_value:", expectedValue)
	}
}

func TestBalanced_distance(t *testing.T) {

	distances_two_data := []float64{10, 29}
	expected_value_two_data := []float64{1, 0}

	// two data (should be 0 and 1)
	resultCities := Balanced_distance(distances_two_data)

	if false == CompareSlices_General(resultCities, expected_value_two_data) {
		t.Fatal("doesnt equal result:", resultCities, "expected_value:", expected_value_two_data)
	}

	distances_three_data := []float64{10, 20, 30}
	expected_value_three_data := []float64{1, 0.5, 0}

	// two data (should be 0  0.5 and 1)
	resultCities_three_data := Balanced_distance(distances_three_data)

	if false == CompareSlices_General(resultCities_three_data, expected_value_three_data) {
		t.Fatal("doesnt equal result:", resultCities_three_data, "expected_value:", expected_value_three_data)
	}

	// not sorted, more data
	distances_five_data := []float64{10, 40, 60, 110, 20}
	expected_value_five_data := []float64{1, 0.7, 0.5, 0, 0.9}

	// two data (should be 0 and 1)
	resultCities_five_data := Balanced_distance(distances_five_data)

	if false == CompareSlices_General(resultCities_five_data, expected_value_five_data) {
		t.Fatal("doesnt equal result:", resultCities_five_data, "expected_value:", expected_value_five_data)
	}

}

func TestCalculate_temps(t *testing.T) {

	balance := []float64{0, 0.5, 1}

	mockCities := []city_structs.CityInfo{
		city_structs.CityInfo{"becs", city_structs.Geo{97, 96}, [5]float64{100, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 100},
		city_structs.CityInfo{"paris", city_structs.Geo{80, 79}, [5]float64{40, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1500},
		city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{10, 3, 17, 5, 6}, [5]float64{0.2, 0.6, 0.4, 0.5, 0.6}, 1100},
	}

	resultTemps := Calculate_temps(balance, mockCities)
	expected_result := []float64{20, 3, 17, 5, 6}

	if false == CompareSlices_General(resultTemps, expected_result) {
		t.Fatal("doesnt equal result:", resultTemps, "expected_value:", expected_result)
	}
}

func TestCalculate_rains(t *testing.T) {

	balance := []float64{0, 0.5, 1}

	mockCities := []city_structs.CityInfo{
		city_structs.CityInfo{"becs", city_structs.Geo{97, 96}, [5]float64{100, 3, 17, 5, 6}, [5]float64{1, 0.6, 0.4, 0.5, 0.6}, 100},
		city_structs.CityInfo{"paris", city_structs.Geo{80, 79}, [5]float64{40, 3, 17, 5, 6}, [5]float64{0.4, 0.6, 0.4, 0.5, 0.6}, 1500},
		city_structs.CityInfo{"bp", city_structs.Geo{94, 92}, [5]float64{10, 3, 17, 5, 6}, [5]float64{0.1, 0.6, 0.4, 0.5, 0.6}, 1100},
	}

	resultRain := Calculate_rain(balance, mockCities)
	expected_result := []float64{0.2, 0.6, 0.4, 0.5, 0.6}

	if false == CompareSlices_General(resultRain, expected_result) {
		t.Fatal("doesnt equal result:", resultRain, "expected_value:", expected_result)
	}
}

// helper functions
func CompareSlices_General(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func CompareSlices_General_city(a, b []city_structs.CityInfo) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
