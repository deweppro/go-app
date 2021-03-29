package application

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

//Sources model
type Sources string

//Decode unmarshal file to model
func (s Sources) Decode(configs ...interface{}) error {
	data, err := ioutil.ReadFile(string(s))
	if err != nil {
		return err
	}
	ext := filepath.Ext(string(s))
	switch ext {
	case ".json":
		return s.unmarshal("json unmarshal", data, json.Unmarshal, configs...)
	case ".yml", ".yaml":
		return s.unmarshal("yaml unmarshal", data, yaml.Unmarshal, configs...)
	case ".toml":
		return s.unmarshal("toml unmarshal", data, toml.Unmarshal, configs...)
	}
	return ErrBadFileFormat
}

func (s Sources) unmarshal(title string, data []byte, call func([]byte, interface{}) error, configs ...interface{}) error {
	for _, conf := range configs {
		if err := call(data, conf); err != nil {
			return errors.Wrap(err, title)
		}
	}
	return nil
}
