// Package result trimmed version of iRacing results.Result to remove dependency
package result

import "github.com/ianhaycox/ir-standings/model"

type Result struct {
	SessionID             model.SessionID             `json:"session_id,omitempty"`
	SubsessionID          model.SubsessionID          `json:"subsession_id,omitempty"`
	CustID                model.CustID                `json:"cust_id"`
	DisplayName           string                      `json:"display_name"`
	FinishPositionInClass model.FinishPositionInClass `json:"finish_position_in_class"`
	LapsComplete          model.LapsComplete          `json:"laps_complete"`
	CarClassID            model.CarClassID            `json:"car_class_id"`
	CarID                 model.CarID                 `json:"car_id"`
	CarName               string                      `json:"car_name"`
	CarNumber             string                      `json:"car_number,omitempty"`
	// Livery                Livery `json:"livery"`
	// Helmet                Helmet `json:"helmet"`
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
