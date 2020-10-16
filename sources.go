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
		return s.json(configs...)
	case ".yml", ".yaml":
		return s.yaml(configs...)
	}
	return ErrBadFileFormat
}

func (s *sources) yaml(configs ...interface{}) error {
	for _, conf := range configs {
		err := yaml.Unmarshal(s.data, conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sources) json(configs ...interface{}) error {
	for _, conf := range configs {
		err := json.Unmarshal(s.data, conf)
		if err != nil {
			return err
		}
	}
	return nil
}
