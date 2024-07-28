package car

import (
	"testing"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/stretchr/testify/assert"
)

func TestCarClasses(t *testing.T) {
	resultCarClasses := []results.CarClasses{
		{
			CarClassID: 81,
			Name:       "GTE cars",
			ShortName:  "GTE",
			CarsInClass: []results.CarsInClass{
				{CarID: 811},
				{CarID: 812},
			},
		},
		{
			CarClassID: 82,
			Name:       "GTO cars",
			ShortName:  "GTO",
			CarsInClass: []results.CarsInClass{
				{CarID: 821},
				{CarID: 822},
			},
		},
		{
			CarClassID: 83,
			Name:       "GTP cars",
			ShortName:  "GTP",
			CarsInClass: []results.CarsInClass{
				{CarID: 831},
			},
		},
	}

	carClasses := NewCarClasses(resultCarClasses)

	carClasses.AddCarName(811, "Ferrari")
	carClasses.AddCarName(811, "Ferrari") // Dup ignored
	carClasses.AddCarName(812, "Ford")
	carClasses.AddCarName(821, "Audi GTO")
	carClasses.AddCarName(822, "Nissan Skyline")
	carClasses.AddCarName(831, "Porsche 962")
	carClasses.AddCarName(999, "Lambo")

	assert.Equal(t, []string{"Audi GTO"}, carClasses.Names([]model.CarID{821}))
	assert.Equal(t, []string{"Audi GTO", "Ferrari", "Porsche 962"}, carClasses.Names([]model.CarID{831, 811, 821}))
	assert.Equal(t, []string{"Lambo"}, carClasses.Names([]model.CarID{999}))
}
