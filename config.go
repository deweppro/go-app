package app

//go:generate easyjson

//BaseConfig config model
//easyjson:json
type BaseConfig struct {
	Env     string `yaml:"env" json:"env" toml:"env"`
	PidFile string `yaml:"pid" json:"pid" toml:"pid"`
	//logger
	Level   uint32 `yaml:"level" json:"level" toml:"level"`
	LogFile string `yaml:"log" json:"log" toml:"log"`
}
