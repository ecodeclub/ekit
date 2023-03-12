package pool

import (
	"context"
)

type TaskObserver interface {
	Observe(ctx context.Context, task Task, target int64)
}

type TaskPoolObserver interface {
	Observe(ctx context.Context, pool TaskPool, target int64)
}
