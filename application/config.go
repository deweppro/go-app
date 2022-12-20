package application

// BaseConfig config model
type BaseConfig struct {
	Env     string `yaml:"env"`
	PidFile string `yaml:"pid"`
	//logger
	Level   uint32 `yaml:"level"`
	LogFile string `yaml:"log"`
}
