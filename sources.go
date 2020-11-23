/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type sources struct {
	filename string
	data     []byte
}

func NewSources(filename string) (*sources, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &sources{
		filename: filename,
		data:     data,
	}, nil
}

func (s *sources) Decode(configs ...interface{}) error {
	ext := filepath.Ext(s.filename)
	switch ext {
	case ".json":
		return s.unmarshal("json unmarshal", json.Unmarshal, configs...)
	case ".yml", ".yaml":
		return s.unmarshal("yaml unmarshal", yaml.Unmarshal, configs...)
	case ".toml":
		return s.unmarshal("toml unmarshal", toml.Unmarshal, configs...)
	}
	return ErrBadFileFormat
}

func (s *sources) unmarshal(title string, call func([]byte, interface{}) error, configs ...interface{}) error {
	for _, conf := range configs {
		if err := call(s.data, conf); err != nil {
			return errors.Wrap(err, title)
		}
	}
	return nil
}
