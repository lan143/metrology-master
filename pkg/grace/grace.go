package grace

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Group struct {
	context.Context

	force context.CancelFunc
	grace chan struct{}

	errg *errgroup.Group
	serv sync.WaitGroup
}

func NewContext(ctx context.Context) *Group {
	errg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)

	return &Group{
		Context: ctx,
		force:   cancel,
		grace:   make(chan struct{}),
		errg:    errg,
	}
}

func (g *Group) Start(fn Routine) {
	g.errg.Go(fn)
}

func (g *Group) Serve(svc Service) {
	g.serv.Add(1)
	g.errg.Go(func() error {
		defer g.serv.Done()
		return g.serve(svc)
	})
}

func (g *Group) serve(svc Service) error {
	ch := make(chan error, 1)
	go func() {
		ch <- svc.Run()
	}()

	select {
	case err := <-ch:
		return err
	case <-g.Done():
	case <-g.grace:
	}

	err := svc.Shutdown()
	if err != nil {
		return err
	}

	return <-ch
}

func (g *Group) Wait() error {
	return g.errg.Wait()
}

func (g *Group) Shutdown() {
	close(g.grace)
	g.serv.Wait()
	g.force()
}
