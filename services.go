/**
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"fmt"
	"reflect"
)

type services struct {
	sequence *sequence
	up       bool
}

type sequence struct {
	Previous *sequence
	Current  Servicer
	Next     *sequence
}

//Servicer interface for services
type Servicer interface {
	Up() error
	Down() error
}

var (
	srvType = reflect.TypeOf(new(Servicer)).Elem()
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
func (s *services) Add(v Servicer) error {
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
func (s *services) Up() (er error) {
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
func (s *services) Down() (er error) {
	if !s.IsUp() {
		return ErrDepNotRunning
	}
	if s.sequence == nil {
		return ErrDepEmpty
	}
	for {
		if err := s.sequence.Current.Down(); err != nil {
			er = WrapErrors(er, err,
				fmt.Sprintf(
					"down %s service error",
					reflect.TypeOf(s.sequence.Current).String(),
				),
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
	return
}
