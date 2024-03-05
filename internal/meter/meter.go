package meter

import "context"

type Meter interface {
	Init(ctx context.Context) error
	GetParams() Params
}
