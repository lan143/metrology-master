package cmd

import "github.com/lan143/metrology-master/pkg/grace"

type Context interface {
	grace.Context

	Shutdown() <-chan struct{}
}
