package cities_func

import (
	"github.com/bednayb/Go_cities/city_structs"
	"testing"
)

func TestCheck_distance(t *testing.T) {

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
	expectedThreeCityValue := []float64{5,10,29}

	if false == CompareSlices_General(resultThreeCity, expectedThreeCityValue) {
		t.Fatal("doesnt equal result:", resultThreeCity, "expected_value:", expectedThreeCityValue)
	}
}


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
