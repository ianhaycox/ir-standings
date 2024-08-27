package car

import (
	"encoding/json"
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/test/files"
	"github.com/stretchr/testify/assert"
)

func TestCarClasses(t *testing.T) {
	var (
		carData      []cars.Car
		carClassData []cars.CarClass
	)

	json.Unmarshal(files.ReadFile(t, "../../fixtures/car.json"), &carData)
	json.Unmarshal(files.ReadFile(t, "../../fixtures/carclass.json"), &carClassData)

	carClasses := NewCarClasses([]int{2268}, carData, carClassData)

	assert.Equal(t, []string{"BMW M4 GT4", "McLaren 570S GT4"}, carClasses.CarNames([]model.CarID{135, 122}))
	assert.Equal(t, []string{"Aston Martin Vantage GT4"}, carClasses.CarNames([]model.CarID{150}))
	assert.Equal(t, "GT4 Class", carClasses.Name(2268))
}
