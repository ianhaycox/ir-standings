package model

// Result trimmed version of iRacing results.Result to remove dependency
type Result struct {
	SessionID               SessionID    `json:"session_id,omitempty"`
	SubsessionID            SubsessionID `json:"subsession_id,omitempty"`
	CustID                  CustID       `json:"cust_id"`
	DisplayName             string       `json:"display_name"`
	FinishPosition          int          `json:"finish_position"`
	FinishPositionInClass   int          `json:"finish_position_in_class"`
	LapsLead                int          `json:"laps_lead"`
	LapsComplete            int          `json:"laps_complete"`
	Position                int          `json:"position"`
	QualLapTime             int          `json:"qual_lap_time"`
	StartingPosition        int          `json:"starting_position"`
	StartingPositionInClass int          `json:"starting_position_in_class"`
	CarClassID              CarClassID   `json:"car_class_id"`
	CarClassName            string       `json:"car_class_name"`
	CarClassShortName       string       `json:"car_class_short_name"`
	ClubID                  int          `json:"club_id"`
	ClubName                string       `json:"club_name"`
	ClubShortname           string       `json:"club_shortname"`
	Division                int          `json:"division"`
	DivisionName            string       `json:"division_name"`
	Incidents               int          `json:"incidents"`
	CarID                   int          `json:"car_id"`
	CarName                 string       `json:"car_name"`
	// Livery                  Livery `json:"livery"`
	// Helmet                  Helmet `json:"helmet"`
}

func (r *Result) IsClassified(winnerLapsComplete int) bool {
	return r.LapsComplete*3 <= winnerLapsComplete*4 // 75%
}

type Livery struct {
	CarID        int         `json:"car_id"`
	Pattern      int         `json:"pattern"`
	Color1       string      `json:"color1"`
	Color2       string      `json:"color2"`
	Color3       string      `json:"color3"`
	NumberFont   int         `json:"number_font"`
	NumberColor1 string      `json:"number_color1"`
	NumberColor2 string      `json:"number_color2"`
	NumberColor3 string      `json:"number_color3"`
	NumberSlant  int         `json:"number_slant"`
	Sponsor1     int         `json:"sponsor1"`
	Sponsor2     int         `json:"sponsor2"`
	CarNumber    string      `json:"car_number"`
	WheelColor   interface{} `json:"wheel_color"`
	RimType      int         `json:"rim_type"`
}

type Helmet struct {
	Pattern    int    `json:"pattern"`
	Color1     string `json:"color1"`
	Color2     string `json:"color2"`
	Color3     string `json:"color3"`
	FaceType   int    `json:"face_type"`
	HelmetType int    `json:"helmet_type"`
}
