package iryaml

type WeekendInfo struct {
	TrackName              string           `yaml:"TrackName"`
	TrackID                int              `yaml:"TrackID"`
	TrackLength            string           `yaml:"TrackLength"`
	TrackLengthOfficial    string           `yaml:"TrackLengthOfficial"`
	TrackDisplayName       string           `yaml:"TrackDisplayName"`
	TrackDisplayShortName  string           `yaml:"TrackDisplayShortName"`
	TrackConfigName        string           `yaml:"TrackConfigName"`
	TrackCity              string           `yaml:"TrackCity"`
	TrackCountry           string           `yaml:"TrackCountry"`
	TrackAltitude          string           `yaml:"TrackAltitude"`
	TrackLatitude          string           `yaml:"TrackLatitude"`
	TrackLongitude         string           `yaml:"TrackLongitude"`
	TrackNorthOffset       string           `yaml:"TrackNorthOffset"`
	TrackNumTurns          int              `yaml:"TrackNumTurns"`
	TrackPitSpeedLimit     string           `yaml:"TrackPitSpeedLimit"`
	TrackType              string           `yaml:"TrackType"`
	TrackDirection         string           `yaml:"TrackDirection"`
	TrackWeatherType       string           `yaml:"TrackWeatherType"`
	TrackSkies             string           `yaml:"TrackSkies"`
	TrackSurfaceTemp       string           `yaml:"TrackSurfaceTemp"`
	TrackAirTemp           string           `yaml:"TrackAirTemp"`
	TrackAirPressure       string           `yaml:"TrackAirPressure"`
	TrackWindVel           string           `yaml:"TrackWindVel"`
	TrackWindDir           string           `yaml:"TrackWindDir"`
	TrackRelativeHumidity  string           `yaml:"TrackRelativeHumidity"`
	TrackFogLevel          string           `yaml:"TrackFogLevel"`
	TrackCleanup           int              `yaml:"TrackCleanup"`
	TrackDynamicTrack      int              `yaml:"TrackDynamicTrack"`
	TrackVersion           string           `yaml:"TrackVersion"`
	SeriesID               int              `yaml:"SeriesID"`
	SeasonID               int              `yaml:"SeasonID"`
	SessionID              int              `yaml:"SessionID"`
	SubSessionID           int              `yaml:"SubSessionID"`
	LeagueID               int              `yaml:"LeagueID"`
	Official               int              `yaml:"Official"`
	RaceWeek               int              `yaml:"RaceWeek"`
	EventType              string           `yaml:"EventType"`
	Category               string           `yaml:"Category"`
	SimMode                string           `yaml:"SimMode"`
	TeamRacing             int              `yaml:"TeamRacing"`
	MinDrivers             int              `yaml:"MinDrivers"`
	MaxDrivers             int              `yaml:"MaxDrivers"`
	DCRuleSet              string           `yaml:"DCRuleSet"`
	QualifierMustStartRace int              `yaml:"QualifierMustStartRace"`
	NumCarClasses          int              `yaml:"NumCarClasses"`
	NumCarTypes            int              `yaml:"NumCarTypes"`
	HeatRacing             int              `yaml:"HeatRacing"`
	BuildType              string           `yaml:"BuildType"`
	BuildTarget            string           `yaml:"BuildTarget"`
	BuildVersion           string           `yaml:"BuildVersion"`
	WeekendOptions         WeekendOptions   `yaml:"WeekendOptions"`
	TelemetryOptions       TelemetryOptions `yaml:"TelemetryOptions"`
}

type WeekendOptions struct {
	NumStarters                int    `yaml:"NumStarters"`
	StartingGrid               string `yaml:"StartingGrid"`
	QualifyScoring             string `yaml:"QualifyScoring"`
	CourseCautions             string `yaml:"CourseCautions"`
	StandingStart              int    `yaml:"StandingStart"`
	ShortParadeLap             int    `yaml:"ShortParadeLap"`
	Restarts                   string `yaml:"Restarts"`
	WeatherType                string `yaml:"WeatherType"`
	Skies                      string `yaml:"Skies"`
	WindDirection              string `yaml:"WindDirection"`
	WindSpeed                  string `yaml:"WindSpeed"`
	WeatherTemp                string `yaml:"WeatherTemp"`
	RelativeHumidity           string `yaml:"RelativeHumidity"`
	FogLevel                   string `yaml:"FogLevel"`
	TimeOfDay                  string `yaml:"TimeOfDay"`
	Date                       string `yaml:"Date"`
	EarthRotationSpeedupFactor int    `yaml:"EarthRotationSpeedupFactor"`
	Unofficial                 int    `yaml:"Unofficial"`
	CommercialMode             string `yaml:"CommercialMode"`
	NightMode                  string `yaml:"NightMode"`
	IsFixedSetup               int    `yaml:"IsFixedSetup"`
	StrictLapsChecking         string `yaml:"StrictLapsChecking"`
	HasOpenRegistration        int    `yaml:"HasOpenRegistration"`
	HardcoreLevel              int    `yaml:"HardcoreLevel"`
	NumJokerLaps               int    `yaml:"NumJokerLaps"`
	IncidentLimit              string `yaml:"IncidentLimit"`
	FastRepairsLimit           string `yaml:"FastRepairsLimit"`
	GreenWhiteCheckeredLimit   int    `yaml:"GreenWhiteCheckeredLimit"`
}

type TelemetryOptions struct {
	TelemetryDiskFile string `yaml:"TelemetryDiskFile"`
}
