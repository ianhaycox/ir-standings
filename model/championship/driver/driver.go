// Package driver in the championship
package driver

import "github.com/ianhaycox/ir-standings/model"

type Driver struct {
	custID      model.CustID
	displayName string
	iRating     int
}

func NewDriver(custID model.CustID, displayName string, iRating int) Driver {
	return Driver{
		custID:      custID,
		displayName: displayName,
		iRating:     iRating,
	}
}

func (d *Driver) CustID() model.CustID {
	return d.custID
}

func (d *Driver) DisplayName() string {
	return d.displayName
}

func (d *Driver) IRating() int {
	return d.iRating
}
