package job

import "context"

type Job interface {
	Execute(ctx context.Context) error
}
