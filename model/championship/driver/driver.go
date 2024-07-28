// Package driver in the championship
package driver

import "github.com/ianhaycox/ir-standings/model"

type Driver struct {
	custID      model.CustID
	displayName string
}

func NewDriver(custID model.CustID, displayName string) Driver {
	return Driver{
		custID:      custID,
		displayName: displayName,
	}
}

func (d *Driver) CustID() model.CustID {
	return d.custID
}

func (d *Driver) DisplayName() string {
	return d.displayName
}
