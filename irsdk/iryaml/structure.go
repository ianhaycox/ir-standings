package iryaml

type IRSession struct {
	WeekendInfo        WeekendInfo           `yaml:"WeekendInfo"`
	SessionInfo        SessionInfo           `yaml:"SessionInfo"`
	QualifyResultsInfo QualifyingResultsInfo `yaml:"QualifyResultsInfo"`
	CameraInfo         CameraInfo            `yaml:"CameraInfo"`
	RadioInfo          RadioInfo             `yaml:"RadioInfo"`
	DriverInfo         DriverInfo            `yaml:"DriverInfo"`
	SplitTimeInfo      SplitTimeInfo         `yaml:"SplitTimeInfo"`
	CarSetup           CarSetup              `yaml:"CarSetup"`
}
