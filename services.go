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
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type services struct {
	sequence *sequence
	up       bool
}

type sequence struct {
	Previous *sequence
	Current  ServiceInterface
	Next     *sequence
}

type ServiceInterface interface {
	Up() error
	Down() error
}

var (
	srvType = reflect.TypeOf(new(ServiceInterface)).Elem()
)

func newServices() *services {
	return &services{
		sequence: nil,
		up:       false,
	}
}

// IsUp - mark that all services have started
func (s *services) IsUp() bool {
	return s.up
}

// Add - add new service by interface
func (s *services) Add(v ServiceInterface) error {
	if s.IsUp() {
		return ErrDepRunning
	}

	if s.sequence == nil {
		s.sequence = &sequence{
			Previous: nil,
			Current:  v,
			Next:     nil,
		}
	} else {
		n := &sequence{
			Previous: s.sequence,
			Current:  v,
			Next:     nil,
		}
		n.Previous.Next = n
		s.sequence = n
	}

	return nil
}

// Up - start all services
func (s *services) Up() error {
	if s.IsUp() {
		return ErrDepRunning
	}
	if s.sequence == nil {
		return ErrDepEmpty
	}
	s.up = true
	for s.sequence.Previous != nil {
		s.sequence = s.sequence.Previous
	}
	for {
		if er := s.sequence.Current.Up(); er != nil {
			return er
		}
		if s.sequence.Next == nil {
			break
		}
		s.sequence = s.sequence.Next
	}

	return nil
}

// Down - stop all services
func (s *services) Down() error {
	if !s.IsUp() {
		return ErrDepNotRunning
	}
	if s.sequence == nil {
		return ErrDepEmpty
	}
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("trace", err).Error("panic on services down")
		}
	}()
	var (
		e error
	)
	for {
		if err := s.sequence.Current.Down(); err != nil {
			e = errors.Wrapf(err,
				"down %s service error",
				reflect.TypeOf(s.sequence.Current).String(),
			)
		}
		if s.sequence.Previous == nil {
			break
		}
		s.sequence = s.sequence.Previous
	}
	for s.sequence.Next != nil {
		s.sequence = s.sequence.Next
	}
	s.up = false
	return e
}
