package application

import "context"

//ForceClose model for force close application
type ForceClose struct {
	C     context.Context
	Close context.CancelFunc
}

func newForceClose() *ForceClose {
	ctx, cncl := context.WithCancel(context.Background())

	return &ForceClose{
		C:     ctx,
		Close: cncl,
	}
}
