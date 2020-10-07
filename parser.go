/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
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
