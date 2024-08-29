// Package iryaml session data
package iryaml

type SessionInfo struct {
	Sessions []Session `yaml:"Sessions"`
}

type Session struct {
	SessionNum                       int                 `yaml:"SessionNum"`
	SessionLaps                      string              `yaml:"SessionLaps"`
	SessionTime                      string              `yaml:"SessionTime"`
	SessionNumLapsToAvg              int                 `yaml:"SessionNumLapsToAvg"`
	SessionType                      string              `yaml:"SessionType"`
	SessionTrackRubberState          string              `yaml:"SessionTrackRubberState"`
	SessionName                      string              `yaml:"SessionName"`
	SessionSubType                   interface{}         `yaml:"SessionSubType"`
	SessionSkipped                   int                 `yaml:"SessionSkipped"`
	SessionRunGroupsUsed             int                 `yaml:"SessionRunGroupsUsed"`
	SessionEnforceTireCompoundChange int                 `yaml:"SessionEnforceTireCompoundChange"`
	ResultsPositions                 []ResultsPosition   `yaml:"ResultsPositions"`
	ResultsFastestLap                []ResultsFastestLap `yaml:"ResultsFastestLap"`
	ResultsAverageLapTime            float64             `yaml:"ResultsAverageLapTime"`
	ResultsNumCautionFlags           int                 `yaml:"ResultsNumCautionFlags"`
	ResultsNumCautionLaps            int                 `yaml:"ResultsNumCautionLaps"`
	ResultsNumLeadChanges            int                 `yaml:"ResultsNumLeadChanges"`
	ResultsLapsComplete              int                 `yaml:"ResultsLapsComplete"`
	ResultsOfficial                  int                 `yaml:"ResultsOfficial"`
}

type ResultsPosition struct {
	Position          int     `yaml:"Position"`
	ClassPosition     int     `yaml:"ClassPosition"`
	CarIdx            int     `yaml:"CarIdx"`
	Lap               int     `yaml:"Lap"`
	Time              float64 `yaml:"Time"`
	FastestLap        int     `yaml:"FastestLap"`
	FastestTime       float64 `yaml:"FastestTime"`
	LastTime          float64 `yaml:"LastTime"`
	LapsLed           int     `yaml:"LapsLed"`
	LapsComplete      int     `yaml:"LapsComplete"`
	JokerLapsComplete int     `yaml:"JokerLapsComplete"`
	LapsDriven        float64 `yaml:"LapsDriven"`
	Incidents         int     `yaml:"Incidents"`
	ReasonOutID       int     `yaml:"ReasonOutID"`
	ReasonOutStr      string  `yaml:"ReasonOutStr"`
}

type ResultsFastestLap struct {
	CarIdx      int     `yaml:"CarIdx"`
	FastestLap  int     `yaml:"FastestLap"`
	FastestTime float64 `yaml:"FastestTime"`
}
