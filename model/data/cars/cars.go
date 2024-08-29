// Package cars iRacing car data
package cars

type Car struct {
	AiEnabled          bool       `json:"ai_enabled"`
	CarDirpath         string     `json:"car_dirpath"`
	CarID              int        `json:"car_id"`
	CarName            string     `json:"car_name"`
	CarNameAbbreviated string     `json:"car_name_abbreviated"`
	CarTypes           []CarTypes `json:"car_types"`
	CarWeight          int        `json:"car_weight"`
	Categories         []string   `json:"categories"`
	CarMake            string     `json:"car_make,omitempty"`
	CarModel           string     `json:"car_model,omitempty"`
	RainEnabled        bool       `json:"rain_enabled"`
	Retired            bool       `json:"retired"`
}
type CarTypes struct {
	CarType string `json:"car_type"`
}

type CarClass struct {
	CarClassID  int           `json:"car_class_id"`
	CarsInClass []CarsInClass `json:"cars_in_class"`
	Name        string        `json:"name"`
	RainEnabled bool          `json:"rain_enabled"`
	ShortName   string        `json:"short_name"`
}
type CarsInClass struct {
	CarDirpath  string `json:"car_dirpath"`
	CarID       int    `json:"car_id"`
	RainEnabled bool   `json:"rain_enabled"`
	Retired     bool   `json:"retired"`
}
