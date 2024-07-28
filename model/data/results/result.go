// Package results response from CDN
package results

import "time"

type Result struct {
	SubsessionID            int               `json:"subsession_id"`
	AllowedLicenses         []AllowedLicenses `json:"allowed_licenses"`
	AssociatedSubsessionIds []int             `json:"associated_subsession_ids"`
	CanProtest              bool              `json:"can_protest"`
	CarClasses              []CarClasses      `json:"car_classes"`
	CautionType             int               `json:"caution_type"`
	CooldownMinutes         int               `json:"cooldown_minutes"`
	CornersPerLap           int               `json:"corners_per_lap"`
	DamageModel             int               `json:"damage_model"`
	DriverChangeParam1      int               `json:"driver_change_param1"`
	DriverChangeParam2      int               `json:"driver_change_param2"`
	DriverChangeRule        int               `json:"driver_change_rule"`
	DriverChanges           bool              `json:"driver_changes"`
	EndTime                 time.Time         `json:"end_time"`
	EventAverageLap         int               `json:"event_average_lap"`
	EventBestLapTime        int               `json:"event_best_lap_time"`
	EventLapsComplete       int               `json:"event_laps_complete"`
	EventStrengthOfField    int               `json:"event_strength_of_field"`
	EventType               int               `json:"event_type"`
	EventTypeName           string            `json:"event_type_name"`
	HeatInfoID              int               `json:"heat_info_id"`
	LicenseCategory         string            `json:"license_category"`
	LicenseCategoryID       int               `json:"license_category_id"`
	LimitMinutes            int               `json:"limit_minutes"`
	MaxTeamDrivers          int               `json:"max_team_drivers"`
	MaxWeeks                int               `json:"max_weeks"`
	MinTeamDrivers          int               `json:"min_team_drivers"`
	NumCautionLaps          int               `json:"num_caution_laps"`
	NumCautions             int               `json:"num_cautions"`
	NumDrivers              int               `json:"num_drivers"`
	NumLapsForQualAverage   int               `json:"num_laps_for_qual_average"`
	NumLapsForSoloAverage   int               `json:"num_laps_for_solo_average"`
	NumLeadChanges          int               `json:"num_lead_changes"`
	OfficialSession         bool              `json:"official_session"`
	PointsType              string            `json:"points_type"`
	PrivateSessionID        int               `json:"private_session_id"`
	RaceSummary             RaceSummary       `json:"race_summary"`
	RaceWeekNum             int               `json:"race_week_num"`
	ResultsRestricted       bool              `json:"results_restricted"`
	SeasonID                int               `json:"season_id"`
	SeasonName              string            `json:"season_name"`
	SeasonQuarter           int               `json:"season_quarter"`
	SeasonShortName         string            `json:"season_short_name"`
	SeasonYear              int               `json:"season_year"`
	SeriesID                int               `json:"series_id"`
	SeriesLogo              string            `json:"series_logo"`
	SeriesName              string            `json:"series_name"`
	SeriesShortName         string            `json:"series_short_name"`
	SessionID               int               `json:"session_id"`
	SessionResults          []SessionResults  `json:"session_results"`
	SessionSplits           []SessionSplits   `json:"session_splits"`
	SpecialEventType        int               `json:"special_event_type"`
	StartTime               time.Time         `json:"start_time"`
	Track                   ResultTrack       `json:"track"`
	TrackState              TrackState        `json:"track_state"`
	Weather                 Weather           `json:"weather"`
}
type AllowedLicenses struct {
	GroupName       string `json:"group_name"`
	LicenseGroup    int    `json:"license_group"`
	MaxLicenseLevel int    `json:"max_license_level"`
	MinLicenseLevel int    `json:"min_license_level"`
	ParentID        int    `json:"parent_id"`
}
type CarsInClass struct {
	CarID int `json:"car_id"`
}
type CarClasses struct {
	CarClassID      int           `json:"car_class_id"`
	ShortName       string        `json:"short_name"`
	Name            string        `json:"name"`
	StrengthOfField int           `json:"strength_of_field"`
	NumEntries      int           `json:"num_entries"`
	CarsInClass     []CarsInClass `json:"cars_in_class"`
}
type RaceSummary struct {
	SubsessionID         int    `json:"subsession_id"`
	AverageLap           int    `json:"average_lap"`
	LapsComplete         int    `json:"laps_complete"`
	NumCautions          int    `json:"num_cautions"`
	NumCautionLaps       int    `json:"num_caution_laps"`
	NumLeadChanges       int    `json:"num_lead_changes"`
	FieldStrength        int    `json:"field_strength"`
	NumOptLaps           int    `json:"num_opt_laps"`
	HasOptPath           bool   `json:"has_opt_path"`
	SpecialEventType     int    `json:"special_event_type"`
	SpecialEventTypeText string `json:"special_event_type_text"`
}
type WeatherResult struct {
	AvgSkies                 int     `json:"avg_skies"`
	AvgCloudCoverPct         float64 `json:"avg_cloud_cover_pct"`
	MinCloudCoverPct         float64 `json:"min_cloud_cover_pct"`
	MaxCloudCoverPct         float64 `json:"max_cloud_cover_pct"`
	TempUnits                int     `json:"temp_units"`
	AvgTemp                  float64 `json:"avg_temp"`
	MinTemp                  float64 `json:"min_temp"`
	MaxTemp                  float64 `json:"max_temp"`
	AvgRelHumidity           float64 `json:"avg_rel_humidity"`
	WindUnits                int     `json:"wind_units"`
	AvgWindSpeed             float64 `json:"avg_wind_speed"`
	MinWindSpeed             float64 `json:"min_wind_speed"`
	MaxWindSpeed             float64 `json:"max_wind_speed"`
	AvgWindDir               int     `json:"avg_wind_dir"`
	MaxFog                   int     `json:"max_fog"`
	FogTimePct               int     `json:"fog_time_pct"`
	PrecipTimePct            int     `json:"precip_time_pct"`
	PrecipMm                 int     `json:"precip_mm"`
	PrecipMm2HrBeforeSession int     `json:"precip_mm2hr_before_session"`
	SimulatedStartTime       string  `json:"simulated_start_time"`
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
type Suit struct {
	Pattern int    `json:"pattern"`
	Color1  string `json:"color1"`
	Color2  string `json:"color2"`
	Color3  string `json:"color3"`
}
type Helmet struct {
	Pattern    int    `json:"pattern"`
	Color1     string `json:"color1"`
	Color2     string `json:"color2"`
	Color3     string `json:"color3"`
	FaceType   int    `json:"face_type"`
	HelmetType int    `json:"helmet_type"`
}
type Results struct {
	CustID                  int       `json:"cust_id"`
	DisplayName             string    `json:"display_name"`
	FinishPosition          int       `json:"finish_position"`
	FinishPositionInClass   int       `json:"finish_position_in_class"`
	LapsLead                int       `json:"laps_lead"`
	LapsComplete            int       `json:"laps_complete"`
	OptLapsComplete         int       `json:"opt_laps_complete"`
	Interval                int       `json:"interval"`
	ClassInterval           int       `json:"class_interval"`
	AverageLap              int       `json:"average_lap"`
	BestLapNum              int       `json:"best_lap_num"`
	BestLapTime             int       `json:"best_lap_time"`
	BestNlapsNum            int       `json:"best_nlaps_num"`
	BestNlapsTime           int       `json:"best_nlaps_time"`
	BestQualLapAt           time.Time `json:"best_qual_lap_at"`
	BestQualLapNum          int       `json:"best_qual_lap_num"`
	BestQualLapTime         int       `json:"best_qual_lap_time"`
	ReasonOutID             int       `json:"reason_out_id"`
	ReasonOut               string    `json:"reason_out"`
	ChampPoints             int       `json:"champ_points"`
	DropRace                bool      `json:"drop_race"`
	ClubPoints              int       `json:"club_points"`
	Position                int       `json:"position"`
	QualLapTime             int       `json:"qual_lap_time"`
	StartingPosition        int       `json:"starting_position"`
	StartingPositionInClass int       `json:"starting_position_in_class"`
	CarClassID              int       `json:"car_class_id"`
	CarClassName            string    `json:"car_class_name"`
	CarClassShortName       string    `json:"car_class_short_name"`
	ClubID                  int       `json:"club_id"`
	ClubName                string    `json:"club_name"`
	ClubShortname           string    `json:"club_shortname"`
	Division                int       `json:"division"`
	DivisionName            string    `json:"division_name"`
	OldLicenseLevel         int       `json:"old_license_level"`
	OldSubLevel             int       `json:"old_sub_level"`
	OldCpi                  float64   `json:"old_cpi"`
	OldiRating              int       `json:"oldi_rating"`
	OldTtrating             int       `json:"old_ttrating"`
	NewLicenseLevel         int       `json:"new_license_level"`
	NewSubLevel             int       `json:"new_sub_level"`
	NewCpi                  float64   `json:"new_cpi"`
	NewiRating              int       `json:"newi_rating"`
	NewTtrating             int       `json:"new_ttrating"`
	Multiplier              int       `json:"multiplier"`
	LicenseChangeOval       int       `json:"license_change_oval"`
	LicenseChangeRoad       int       `json:"license_change_road"`
	Incidents               int       `json:"incidents"`
	MaxPctFuelFill          int       `json:"max_pct_fuel_fill"`
	WeightPenaltyKg         int       `json:"weight_penalty_kg"`
	LeaguePoints            int       `json:"league_points"`
	LeagueAggPoints         int       `json:"league_agg_points"`
	CarID                   int       `json:"car_id"`
	CarName                 string    `json:"car_name"`
	AggregateChampPoints    int       `json:"aggregate_champ_points"`
	Livery                  Livery    `json:"livery"`
	Suit                    Suit      `json:"suit"`
	Helmet                  Helmet    `json:"helmet"`
	Watched                 bool      `json:"watched"`
	Friend                  bool      `json:"friend"`
	Ai                      bool      `json:"ai"`
}
type SessionResults struct {
	SimsessionNumber   int    `json:"simsession_number"`
	SimsessionType     int    `json:"simsession_type"`
	SimsessionTypeName string `json:"simsession_type_name"`
	SimsessionSubtype  int    `json:"simsession_subtype"`
	SimsessionName     string `json:"simsession_name"`
	// WeatherResult      WeatherResult `json:"weather_result"`
	Results []Results `json:"results"`
}
type SessionSplits struct {
	SubsessionID         int `json:"subsession_id"`
	EventStrengthOfField int `json:"event_strength_of_field"`
}
type ResultTrack struct {
	Category   string `json:"category"`
	CategoryID int    `json:"category_id"`
	ConfigName string `json:"config_name"`
	TrackID    int    `json:"track_id"`
	TrackName  string `json:"track_name"`
}
type TrackState struct {
	LeaveMarbles         bool `json:"leave_marbles"`
	PracticeGripCompound int  `json:"practice_grip_compound"`
	PracticeRubber       int  `json:"practice_rubber"`
	QualifyGripCompound  int  `json:"qualify_grip_compound"`
	QualifyRubber        int  `json:"qualify_rubber"`
	RaceGripCompound     int  `json:"race_grip_compound"`
	RaceRubber           int  `json:"race_rubber"`
	WarmupGripCompound   int  `json:"warmup_grip_compound"`
	WarmupRubber         int  `json:"warmup_rubber"`
}
type Weather struct {
	AllowFog           bool   `json:"allow_fog"`
	Fog                int    `json:"fog"`
	PrecipOption       int    `json:"precip_option"`
	RelHumidity        int    `json:"rel_humidity"`
	SimulatedStartTime string `json:"simulated_start_time"`
	Skies              int    `json:"skies"`
	TempUnits          int    `json:"temp_units"`
	TempValue          int    `json:"temp_value"`
	TimeOfDay          int    `json:"time_of_day"`
	TrackWater         int    `json:"track_water"`
	Type               int    `json:"type"`
	Version            int    `json:"version"`
	WeatherVarInitial  int    `json:"weather_var_initial"`
	WeatherVarOngoing  int    `json:"weather_var_ongoing"`
	WindDir            int    `json:"wind_dir"`
	WindUnits          int    `json:"wind_units"`
	WindValue          int    `json:"wind_value"`
}
