package iryaml

type Sector struct {
	SectorNum      int     `yaml:"SectorNum"`
	SectorStartPct float64 `yaml:"SectorStartPct"`
}

type SplitTimeInfo struct {
	Sectors []Sector `yaml:"Sectors"`
}
