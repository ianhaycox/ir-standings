package iryaml

type QualifyingResultsInfo struct {
	Results []QualifyingResult `yaml:"Results"`
}

type QualifyingResult struct {
	Position      int     `yaml:"Position"`
	ClassPosition int     `yaml:"ClassPosition"`
	CarIdx        int     `yaml:"CarIdx"`
	FastestLap    int     `yaml:"FastestLap"`
	FastestTime   float64 `yaml:"FastestTime"`
}
