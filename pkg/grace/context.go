package grace

import "context"

type (
	Context interface {
		context.Context

		Serve(Service)
		Start(Routine)
	}

	Service interface {
		Run() error
		Shutdown() error
	}

	Routine func() error
)
