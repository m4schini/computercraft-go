package computer

import (
	"context"
	"time"
)

type GPS interface {
	Locate(ctx context.Context) (int, int, int, error)
	LocateWithTimeout(ctx context.Context, timeout time.Duration) (int, int, int, error)
}
