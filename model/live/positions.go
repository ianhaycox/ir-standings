package live

// LiveResults sent from Windows Overlay
type LiveResults struct {
	SeriesID     int               `json:"series_id,omitempty"`
	SessionID    int               `json:"session_id,omitempty"`
	SubsessionID int               `json:"subsession_id,omitempty"`
	Track        string            `json:"track,omitempty"`
	CountBestOf  int               `json:"count_best_of,omitempty"`
	CarClassID   int               `json:"car_class_id,omitempty"`
	TopN         int               `json:"top_n,omitempty"`
	Positions    []CurrentPosition `json:"positions,omitempty"`
}

type CurrentPosition struct {
	CustID                int `json:"cust_id,omitempty"`
	FinishPositionInClass int `json:"finish_position_in_class,omitempty"`
	LapsComplete          int `json:"laps_complete,omitempty"`
	CarID                 int `json:"car_id,omitempty"`
}
