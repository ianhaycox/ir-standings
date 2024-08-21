// Package car to manage the car names and classes for an event
package car

import (
	"sort"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/data/results"
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

func (c CarClass) Name() string {
	return c.name
}

type CarClasses struct {
	carClassNames map[model.CarClassID]CarClass
	carNames      map[model.CarID]car
}

func NewCarClasses(carClasses []results.CarClasses) CarClasses {
	cc := CarClasses{
		carClassNames: make(map[model.CarClassID]CarClass),
		carNames:      make(map[model.CarID]car),
	}

	for i := range carClasses {
		carClassID := model.CarClassID(carClasses[i].CarClassID)

		carClass := CarClass{
			name:       carClasses[i].Name,
			shortName:  carClasses[i].ShortName,
			carClassID: carClassID,
		}

		// We only find out the car names later
		for j := range carClasses[i].CarsInClass {
			carID := model.CarID(carClasses[i].CarsInClass[j].CarID)
			car := car{carID: carID}

			carClass.carsInClass = append(carClass.carsInClass, car)
			cc.carClassNames[carClassID] = carClass
			cc.carNames[carID] = car
		}
	}

	return cc
}

func (cc *CarClasses) AddCarName(carID model.CarID, name string) {
	if name == "" {
		return
	}

	if carr, ok := cc.carNames[carID]; ok {
		carr.name = name
		cc.carNames[carID] = carr
	} else {
		cc.carNames[carID] = car{carID: carID, name: name}
	}

	for carClassID := range cc.carClassNames {
		for i := range cc.carClassNames[carClassID].carsInClass {
			if cc.carClassNames[carClassID].carsInClass[i].carID == carID {
				cc.carClassNames[carClassID].carsInClass[i].name = name

				break
			}
		}

		cic := cc.carClassNames[carClassID].carsInClass
		cic = append(cic, car{carID: carID, name: name})
		cc.carClassNames[carClassID] = CarClass{carsInClass: cic}
	}
}

func (cc *CarClasses) Names(carsDriven []model.CarID) []string {
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

func (cc *CarClasses) ClassNames() map[model.CarClassID]CarClass {
	return cc.carClassNames
}
