package live

// PredictedStanding response to Windows Overlay
type PredictedStanding struct {
	DriverName        string `json:"driver_name"`
	CarNumber         string `json:"car_number"`
	CurrentPosition   int    `json:"current_position"`
	PredictedPosition int    `json:"predicted_position"`
	CurrentPoints     int    `json:"current_points"`
	PredictedPoints   int    `json:"predicted_points"`
	Change            int    `json:"change"` // +/- change from current position
}
