package iryaml

type Cameras struct {
	CameraNum  int    `yaml:"CameraNum"`
	CameraName string `yaml:"CameraName"`
}

type CameraGroup struct {
	GroupNum  int       `yaml:"GroupNum"`
	GroupName string    `yaml:"GroupName"`
	Cameras   []Cameras `yaml:"Cameras"`
	IsScenic  bool      `yaml:"IsScenic,omitempty"`
}

type CameraInfo struct {
	Groups []CameraGroup `yaml:"Groups"`
}
