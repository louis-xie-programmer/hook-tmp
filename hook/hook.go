package hook

import "context"

type Mode int

const (
	Sync Mode = iota
	Async
)

type Hook struct {
	Name        string
	Priority    int
	MustSucceed bool
	Mode        Mode
	Fn          func(ctx context.Context, data any) error
}
