package service

import (
	"sync/atomic"

	"github.com/deweppro/go-app/application/ctx"
	e "github.com/deweppro/go-app/application/error"
	"github.com/deweppro/go-errors"
)

const (
	statusUp   uint32 = 1
	statusDown uint32 = 0
)

type (
	//Tree service initialization tree
	Tree struct {
		sequence *sequence
		status   uint32
	}
	sequence struct {
		Previous *sequence
		Current  interface{}
		Next     *sequence
	}
)

type (
	//IService interface for services
	IService interface {
		Up() error
		Down() error
	}
	//IServiceCtx interface for services with context
	IServiceCtx interface {
		Up(ctx ctx.Context) error
		Down(ctx ctx.Context) error
	}
)

//New tree constructor
func New() *Tree {
	return &Tree{
		sequence: nil,
		status:   statusDown,
	}
}

// IsUp - mark that all services have started
func (s *Tree) IsUp() bool {
	return atomic.LoadUint32(&s.status) == statusUp
}

// Add - add new service by interface
func (s *Tree) Add(v interface{}) error {
	if s.IsUp() {
		return e.ErrDepRunning
	}

	if !IsService(v) {
		return errors.WrapMessage(e.ErrServiceUnknown, "service <%T>", v)
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
func (s *Tree) Up(ctx ctx.Context) error {
	if !atomic.CompareAndSwapUint32(&s.status, statusDown, statusUp) {
		return e.ErrDepRunning
	}
	if s.sequence == nil {
		return nil
	}
	for s.sequence.Previous != nil {
		s.sequence = s.sequence.Previous
	}
	for {
		if vv, ok := s.sequence.Current.(IService); ok {
			if err := vv.Up(); err != nil {
				return err
			}
		} else if vv, ok := s.sequence.Current.(IServiceCtx); ok {
			if err := vv.Up(ctx); err != nil {
				return err
			}
		} else {
			return errors.WrapMessage(e.ErrServiceUnknown, "service <%T>", s.sequence.Current)
		}
		if s.sequence.Next == nil {
			break
		}
		s.sequence = s.sequence.Next
	}

	return nil
}

// Down - stop all services
func (s *Tree) Down(ctx ctx.Context) (er error) {
	if !atomic.CompareAndSwapUint32(&s.status, statusUp, statusDown) {
		return e.ErrDepNotRunning
	}
	if s.sequence == nil {
		return nil
	}
	for {
		if vv, ok := s.sequence.Current.(IService); ok {
			if err := vv.Down(); err != nil {
				er = errors.Wrap(er,
					errors.WrapMessage(err, "down <%T> service error", s.sequence.Current),
				)
			}
		} else if vv, ok := s.sequence.Current.(IServiceCtx); ok {
			if err := vv.Down(ctx); err != nil {
				er = errors.Wrap(er,
					errors.WrapMessage(err, "down <%T> service error", s.sequence.Current),
				)
			}
		} else {
			return errors.WrapMessage(e.ErrServiceUnknown, "service <%T>", s.sequence.Current)
		}
		if s.sequence.Previous == nil {
			break
		}
		s.sequence = s.sequence.Previous
	}
	for s.sequence.Next != nil {
		s.sequence = s.sequence.Next
	}
	return
}
