// Package searchseries response from CDN
package searchseries

import "time"

type SearchSeriesResult struct {
	SessionID            int       `json:"session_id"`
	SubsessionID         int       `json:"subsession_id"`
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	LicenseCategoryID    int       `json:"license_category_id"`
	LicenseCategory      string    `json:"license_category"`
	NumDrivers           int       `json:"num_drivers"`
	NumCautions          int       `json:"num_cautions"`
	NumCautionLaps       int       `json:"num_caution_laps"`
	NumLeadChanges       int       `json:"num_lead_changes"`
	EventLapsComplete    int       `json:"event_laps_complete"`
	DriverChanges        bool      `json:"driver_changes"`
	WinnerGroupID        int       `json:"winner_group_id"`
	WinnerName           string    `json:"winner_name,omitempty"`
	WinnerAi             bool      `json:"winner_ai"`
	Track                Track     `json:"track"`
	OfficialSession      bool      `json:"official_session"`
	SeasonID             int       `json:"season_id"`
	SeasonYear           int       `json:"season_year"`
	SeasonQuarter        int       `json:"season_quarter"`
	EventType            int       `json:"event_type"`
	EventTypeName        string    `json:"event_type_name"`
	SeriesID             int       `json:"series_id"`
	SeriesName           string    `json:"series_name"`
	SeriesShortName      string    `json:"series_short_name"`
	RaceWeekNum          int       `json:"race_week_num"`
	EventStrengthOfField int       `json:"event_strength_of_field"`
	EventAverageLap      int       `json:"event_average_lap"`
	EventBestLapTime     int       `json:"event_best_lap_time"`
}
type Track struct {
	ConfigName string `json:"config_name"`
	TrackID    int    `json:"track_id"`
	TrackName  string `json:"track_name"`
}

func (ssr *SearchSeriesResult) IsBroadcast() bool {
	const broadcastRaceTime = 17

	return ssr.StartTime.Weekday() == time.Saturday && ssr.StartTime.Hour() == broadcastRaceTime
}
