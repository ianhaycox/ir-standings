// Package car to manage the car names and classes for an event
package car

import (
	"sort"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/data/cars"
)

type car struct {
	name  string
	carID model.CarID
}

type CarClass struct {
	name        string
	shortName   string
	carClassID  model.CarClassID
	carsInClass []car
}

type CarClasses struct {
	carClassNames map[model.CarClassID]CarClass
	carNames      map[model.CarID]car
}

func NewCarClasses(carClassIds []int, cars []cars.Car, classes []cars.CarClass) CarClasses {
	cc := CarClasses{
		carClassNames: make(map[model.CarClassID]CarClass),
		carNames:      make(map[model.CarID]car),
	}

	for i := range carClassIds {
		for j := range classes {
			if classes[j].CarClassID == carClassIds[i] {
				carClass := CarClass{
					name:       classes[j].Name,
					shortName:  classes[j].ShortName,
					carClassID: model.CarClassID(classes[j].CarClassID),
				}

				for k := range classes[j].CarsInClass {
					for l := range cars {
						if cars[l].CarID == classes[j].CarsInClass[k].CarID {
							carName := car{
								carID: model.CarID(classes[j].CarsInClass[k].CarID),
								name:  cars[l].CarName,
							}

							carClass.carsInClass = append(carClass.carsInClass, carName)
							cc.carNames[carName.carID] = carName

							break
						}
					}
				}

				cc.carClassNames[model.CarClassID(carClassIds[i])] = carClass

				break
			}
		}
	}

	return cc
}

func (cc *CarClasses) CarNames(carsDriven []model.CarID) []string {
	names := make([]string, 0)

	unique := make(map[model.CarID]string)

	for _, carID := range carsDriven {
		unique[carID] = cc.carNames[carID].name
	}

	for _, name := range unique {
		names = append(names, name)
	}

	sort.SliceStable(names, func(i, j int) bool { return names[i] < names[j] })

	return names
}

func (cc *CarClasses) Name(carClassID model.CarClassID) string {
	return cc.carClassNames[carClassID].name
}

func (cc *CarClasses) CarClassIDs() []int {
	carClassIDs := []int{}

	for k := range cc.carClassNames {
		carClassIDs = append(carClassIDs, int(k))
	}

	return carClassIDs
}
