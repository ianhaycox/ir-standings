// Package live predicted standings retrieved by the frontend
package live

import "github.com/ianhaycox/ir-standings/model"

// {CarClassID: 84, ShortName: "GTP", Name: "Nissan GTP ZX-T", CarsInClass: []results.CarsInClass{{CarID: 77}}},
// {CarClassID: 83, ShortName: "GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},

type PredictedStandings struct {
	Status         string                        `json:"status"` // iRacing connection status/session, Race, Qualifying,...
	TrackName      string                        `json:"track_name"`
	CountBestOf    int                           `json:"count_best_of"`
	SelfCarClassID int                           `json:"self_car_class_id"`
	CarClassIDs    []int                         `json:"car_class_ids"`
	Standings      map[model.CarClassID]Standing `json:"standings"`
}

type Standing struct {
	SoFByCarClass           int                 `json:"sof_by_car_class"`
	CarClassID              model.CarClassID    `json:"car_class_id"`
	CarClassName            string              `json:"car_class_name"`
	ClassLeaderLapsComplete model.LapsComplete  `json:"class_leader_laps_complete"`
	Items                   []PredictedStanding `json:"items"`
}

type PredictedStanding struct {
	Driving           bool                        `json:"driving"`            // In current session
	CustID            model.CustID                `json:"cust_id"`            // Key for React
	DriverName        string                      `json:"driver_name"`        // Driver
	CarNumber         string                      `json:"car_number"`         // May be blank if not in the session
	CurrentPosition   model.FinishPositionInClass `json:"current_position"`   // May be first race in the current session
	PredictedPosition model.FinishPositionInClass `json:"predicted_position"` // Championship position as is
	CurrentPoints     model.Point                 `json:"current_points"`     // Championship position before race
	PredictedPoints   model.Point                 `json:"predicted_points"`   // Championship points as is
	Change            int                         `json:"change"`             // +/- change from current position
	CarNames          []string                    `json:"car_names"`          // Cars driven in this class
}
