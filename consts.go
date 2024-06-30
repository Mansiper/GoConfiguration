package configuration

const (
	ftUnknown formatType = iota
	ftEnvironment
	FtEnv
	FtIni
	FtJson
	// FtYaml
)

const (
	sepDefault  = ","
	sep2Default = ":"
	ignoreField = "-"
	nilDefault  = "*nil"
	nowTime     = "now"
)

const (
	vtEmpty valueType = iota
	vtAny
	vtString
	vtNumber
	vtBool
	vtNull
)
