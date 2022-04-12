package ctx

import "context"

type (
	ctx struct {
		ctx  context.Context
		cncl context.CancelFunc
	}
	//Context model for force close application
	Context interface {
		Close()
		Context() context.Context
		Done() <-chan struct{}
	}
)

//Close context close method
func (v *ctx) Close() {
	v.cncl()
}

//Context general context
func (v *ctx) Context() context.Context {
	return v.ctx
}

//Done context close wait channel
func (v *ctx) Done() <-chan struct{} {
	return v.ctx.Done()
}

//New init default context
func New() Context {
	wctx, cncl := context.WithCancel(context.Background())

	return &ctx{
		ctx:  wctx,
		cncl: cncl,
	}
}
