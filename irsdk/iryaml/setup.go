package iryaml

type CarSetup struct {
	UpdateCount int            `yaml:"UpdateCount"`
	TiresAero   SetupTiresAero `yaml:"TiresAero"`
	Chassis     Chassis        `yaml:"Chassis"`
}

type SetupTiresAero struct {
	LeftFront       SetupAeroFront       `yaml:"LeftFront"`
	RightFront      SetupAeroFront       `yaml:"RightFront"`
	LeftRear        SetupAeroRear        `yaml:"LeftRear"`
	RightRear       SetupAeroRear        `yaml:"RightRear"`
	AeroBalanceCalc SetupAeroBalanceCalc `yaml:"AeroBalanceCalc"`
}

type SetupAeroFront struct {
	StartingPressure string `yaml:"StartingPressure"`
	LastHotPressure  string `yaml:"LastHotPressure"`
	LastTempsOMI     string `yaml:"LastTempsOMI"`
	TreadRemaining   string `yaml:"TreadRemaining"`
}

type SetupAeroRear struct {
	StartingPressure string `yaml:"StartingPressure"`
	LastHotPressure  string `yaml:"LastHotPressure"`
	LastTempsIMO     string `yaml:"LastTempsIMO"`
	TreadRemaining   string `yaml:"TreadRemaining"`
}

type SetupAeroBalanceCalc struct {
	FrontRhAtSpeed string `yaml:"FrontRhAtSpeed"`
	RearRhAtSpeed  string `yaml:"RearRhAtSpeed"`
	RearWingAngle  string `yaml:"RearWingAngle"`
	FrontDownforce string `yaml:"FrontDownforce"`
}

type Chassis struct {
	Front      ChassisFront      `yaml:"Front"`
	LeftFront  ChassisFrontWheel `yaml:"LeftFront"`
	LeftRear   ChassisRearWheel  `yaml:"LeftRear"`
	RightFront ChassisFrontWheel `yaml:"RightFront"`
	RightRear  ChassisRearWheel  `yaml:"RightRear"`
	Rear       ChassisRear       `yaml:"Rear"`
	InCarDials InCarDials        `yaml:"InCarDials"`
}

type ChassisFront struct {
	ArbBlades          int    `yaml:"ArbBlades"`
	ToeIn              string `yaml:"ToeIn"`
	FrontMasterCyl     string `yaml:"FrontMasterCyl"`
	RearMasterCyl      string `yaml:"RearMasterCyl"`
	BrakePads          string `yaml:"BrakePads"`
	LeftNightLedStrip  string `yaml:"LeftNightLedStrip"`
	RightNightLedStrip string `yaml:"RightNightLedStrip"`
}

type ChassisRear struct {
	FuelLevel     string `yaml:"FuelLevel"`
	ArbBlades     int    `yaml:"ArbBlades"`
	GearStack     string `yaml:"GearStack"`
	FrictionFaces int    `yaml:"FrictionFaces"`
	DiffPreload   string `yaml:"DiffPreload"`
	RearWingAngle string `yaml:"RearWingAngle"`
}

type ChassisFrontWheel struct {
	CornerWeight      string `yaml:"CornerWeight"`
	RideHeight        string `yaml:"RideHeight"`
	SpringPerchOffset string `yaml:"SpringPerchOffset"`
	SpringRate        string `yaml:"SpringRate"`
	LsCompDamping     string `yaml:"LsCompDamping"`
	HsCompDamping     string `yaml:"HsCompDamping"`
	LsRbdDamping      string `yaml:"LsRbdDamping"`
	HsRbdDamping      string `yaml:"HsRbdDamping"`
	Camber            string `yaml:"Camber"`
	Caster            string `yaml:"Caster"`
}

type ChassisRearWheel struct {
	CornerWeight      string `yaml:"CornerWeight"`
	RideHeight        string `yaml:"RideHeight"`
	SpringPerchOffset string `yaml:"SpringPerchOffset"`
	SpringRate        string `yaml:"SpringRate"`
	LsCompDamping     string `yaml:"LsCompDamping"`
	HsCompDamping     string `yaml:"HsCompDamping"`
	LsRbdDamping      string `yaml:"LsRbdDamping"`
	HsRbdDamping      string `yaml:"HsRbdDamping"`
	Camber            string `yaml:"Camber"`
	ToeIn             string `yaml:"ToeIn"`
}

type InCarDials struct {
	BrakePressureBias      string `yaml:"BrakePressureBias"`
	AbsSetting             string `yaml:"AbsSetting"`
	TractionControlSetting string `yaml:"TractionControlSetting"`
	ThrottleShapeSetting   string `yaml:"ThrottleShapeSetting"`
	DisplayPage            string `yaml:"DisplayPage"`
	CrossWeight            string `yaml:"CrossWeight"`
}
