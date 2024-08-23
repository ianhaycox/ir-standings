// Package seasons info
package seasons

type Season struct {
	SeasonID                int         `json:"season_id"`
	SeasonName              string      `json:"season_name"`
	Active                  bool        `json:"active"`
	CarClassIds             []int       `json:"car_class_ids"`
	CarTypes                []CarTypes  `json:"car_types"`
	Complete                bool        `json:"complete"`
	Drops                   int         `json:"drops"`
	MaxWeeks                int         `json:"max_weeks"`
	Multiclass              bool        `json:"multiclass"`
	Official                bool        `json:"official"`
	RaceWeek                int         `json:"race_week"`
	RaceWeekToMakeDivisions int         `json:"race_week_to_make_divisions"`
	ScheduleDescription     string      `json:"schedule_description"`
	Schedules               []Schedules `json:"schedules"`
	SeasonQuarter           int         `json:"season_quarter"`
	SeasonShortName         string      `json:"season_short_name"`
	SeasonYear              int         `json:"season_year"`
	SeriesID                int         `json:"series_id"`
	RacePoints              int         `json:"race_points"`
}

type CarTypes struct {
	CarType string `json:"car_type"`
}

type RaceTimeDescriptors struct {
	DayOffset        []int  `json:"day_offset"`
	FirstSessionTime string `json:"first_session_time"`
	RepeatMinutes    int    `json:"repeat_minutes"`
	Repeating        bool   `json:"repeating"`
	SessionMinutes   int    `json:"session_minutes"`
	StartDate        string `json:"start_date"`
	SuperSession     bool   `json:"super_session"`
}

type Schedules struct {
	SeasonID            int                   `json:"season_id"`
	RaceTimeDescriptors []RaceTimeDescriptors `json:"race_time_descriptors"`
	RaceWeekNum         int                   `json:"race_week_num"`
	ScheduleName        string                `json:"schedule_name"`
	SeasonName          string                `json:"season_name"`
	SeriesID            int                   `json:"series_id"`
	SeriesName          string                `json:"series_name"`
	StartDate           string                `json:"start_date"`
}
