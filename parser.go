/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type sources struct {
	filename string
	data     []byte
}

func newSources(filename string) (*sources, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &sources{
		filename: filename,
		data:     data,
	}, nil
}

func (s *sources) Filename() string {
	return s.filename
}

func (s *sources) Data() []byte {
	return s.data
}

func (s *sources) YAML(configs ...interface{}) error {
	for _, conf := range configs {
		err := yaml.Unmarshal(s.data, conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sources) JSON(configs ...interface{}) error {
	for _, conf := range configs {
		err := json.Unmarshal(s.data, conf)
		if err != nil {
			return err
		}
	}
	return nil
}
