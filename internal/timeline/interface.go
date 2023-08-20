package timeline

import (
	"context"
	"time"
)

type Handler func()

type Defer interface {
	Defer(ctx context.Context, handler Handler, after time.Duration)
}

type Timeline interface {
	Defer
	Stop()
}
