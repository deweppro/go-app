package source

import (
	"io/ioutil"
	"path/filepath"

	e "github.com/deweppro/go-app/application/error"
	"github.com/deweppro/go-errors"
	"gopkg.in/yaml.v3"
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
	case ".yml", ".yaml":
		return s.unmarshal("yaml unmarshal", data, yaml.Unmarshal, configs...)
	}
	return e.ErrBadFileFormat
}

func (s Sources) unmarshal(title string, data []byte, call func([]byte, interface{}) error, configs ...interface{}) error {
	for _, conf := range configs {
		if err := call(data, conf); err != nil {
			return errors.WrapMessage(err, title)
		}
	}
	return nil
}
