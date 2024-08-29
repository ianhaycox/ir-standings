package iryaml

type Frequencies struct {
	FrequencyNum  int    `yaml:"FrequencyNum"`
	FrequencyName string `yaml:"FrequencyName"`
	Priority      int    `yaml:"Priority"`
	CarIdx        int    `yaml:"CarIdx"`
	EntryIdx      int    `yaml:"EntryIdx"`
	ClubID        int    `yaml:"ClubID"`
	CanScan       int    `yaml:"CanScan"`
	CanSquawk     int    `yaml:"CanSquawk"`
	Muted         int    `yaml:"Muted"`
	IsMutable     int    `yaml:"IsMutable"`
	IsDeletable   int    `yaml:"IsDeletable"`
}

type Radios struct {
	RadioNum            int           `yaml:"RadioNum"`
	HopCount            int           `yaml:"HopCount"`
	NumFrequencies      int           `yaml:"NumFrequencies"`
	TunedToFrequencyNum int           `yaml:"TunedToFrequencyNum"`
	ScanningIsOn        int           `yaml:"ScanningIsOn"`
	Frequencies         []Frequencies `yaml:"Frequencies"`
}

type RadioInfo struct {
	SelectedRadioNum int      `yaml:"SelectedRadioNum"`
	Radios           []Radios `yaml:"Radios"`
}
