package cmd

import (
	"context"
	"github.com/lan143/metrology-master/pkg/cmd/cmdlog"
	"github.com/lan143/metrology-master/pkg/grace"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

type Grace struct {
	*grace.Group

	grace <-chan struct{}

	log *zap.Logger
}

func NewContext(ctx context.Context, log *zap.Logger) (*Grace, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	group := grace.NewContext(ctx)
	g := &Grace{
		Group: group,
		grace: group.Done(),
		log:   log,
	}

	return g, cancel
}

func (g *Grace) Shutdown() <-chan struct{} {
	return g.grace
}

func (g *Grace) Listen(sigs []os.Signal, delay time.Duration) {
	g.log.Debug(
		"register shutdown signal",
		cmdlog.Signals("signals", sigs),
	)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)

	stop := make(chan struct{})
	g.grace = stop

	go func() {
		defer signal.Stop(ch)

		select {
		case <-g.Done():
			close(stop)
		case sig := <-ch:
			g.log.Info("shutdown requested", cmdlog.Signal("signal", sig))
			close(stop)
			g.shutdown(delay)
		}
	}()
}

func (g *Grace) shutdown(d time.Duration) {
	err := g.coolDown(d)
	switch {
	case err != nil:
		g.log.Warn("force shutdown")
	default:
		g.log.Info("graceful shutdown")
	}

	g.Group.Shutdown()
}

func (g *Grace) coolDown(d time.Duration) error {
	if d <= 0 {
		return nil
	}

	g.log.Info(
		"waiting for grace period",
		zap.Duration("period", d),
	)

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-g.Done():
		return g.Err()
	case <-t.C:
	}

	return nil
}
