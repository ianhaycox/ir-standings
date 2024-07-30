package live

// LiveStandings response to Windows Overlay
type LiveStandings struct {
	DriverName        string `json:"driver_name,omitempty"`
	CurrentPosition   int    `json:"current_position,omitempty"`
	PredictedPosition int    `json:"predicted_position,omitempty"`
	CurrentPoints     int    `json:"current_points,omitempty"`
	PredictedPoints   int    `json:"predicted_points,omitempty"`
	Change            int    `json:"change,omitempty"` // +/- change from current position
}
