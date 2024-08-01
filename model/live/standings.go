package live

// PredictedStanding response to Windows Overlay
type PredictedStanding struct {
	DriverName        string `json:"driver_name,omitempty"`
	CarNumber         string `json:"car_number,omitempty"`
	CurrentPosition   int    `json:"current_position,omitempty"`
	PredictedPosition int    `json:"predicted_position,omitempty"`
	CurrentPoints     int    `json:"current_points,omitempty"`
	PredictedPoints   int    `json:"predicted_points,omitempty"`
	Change            int    `json:"change,omitempty"` // +/- change from current position
}
