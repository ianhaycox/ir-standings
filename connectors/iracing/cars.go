package iracing

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/model/data/results"
)

func (ir *IracingAPI) Cars(ctx context.Context) ([]cars.Car, error) {
	var link results.ResultLink

	queryParams := url.Values{}

	err := ir.data.Get(ctx, &link, Endpoint+"/data/car/get", queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get cars, err:%w", err)
	}

	var cars []cars.Car

	err = ir.data.CDN(ctx, link.Link, &cars)
	if err != nil {
		return nil, fmt.Errorf("can not get cars result:%s, err:%w", link.Link, err)
	}

	return cars, nil
}

func (ir *IracingAPI) CarClasses(ctx context.Context) ([]cars.CarClass, error) {
	var link results.ResultLink

	queryParams := url.Values{}

	err := ir.data.Get(ctx, &link, Endpoint+"/data/carclass/get", queryParams)
	if err != nil {
		return nil, fmt.Errorf("could not get car classes, err:%w", err)
	}

	var carclasses []cars.CarClass

	err = ir.data.CDN(ctx, link.Link, &carclasses)
	if err != nil {
		return nil, fmt.Errorf("can not get car classes result:%s, err:%w", link.Link, err)
	}

	return carclasses, nil
}
