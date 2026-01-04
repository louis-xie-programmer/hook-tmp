package hook

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Type int

const (
	Before Type = iota
	After
)

type Engine struct {
	mu    sync.RWMutex
	hooks map[Type][]*Hook
}

func NewEngine() *Engine {
	return &Engine{
		hooks: make(map[Type][]*Hook),
	}
}

// 注册 Hook（通常在 init 或启动阶段完成）
func (e *Engine) Register(t Type, h *Hook) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.hooks[t] = append(e.hooks[t], h)

	// 注册后立即按优先级排序（高优先级先执行）
	sort.Slice(e.hooks[t], func(i, j int) bool {
		return e.hooks[t][i].Priority > e.hooks[t][j].Priority
	})
}

// 执行 Hook
func (e *Engine) Execute(ctx context.Context, t Type, data any) error {
	e.mu.RLock()
	hooks := e.hooks[t]
	e.mu.RUnlock()

	for _, h := range hooks {
		run := func() error {
			start := time.Now()
			defer recordCost(h.Name, start)

			defer func() {
				if r := recover(); r != nil {
					recordPanic(h.Name, r)
				}
			}()

			return h.Fn(ctx, data)
		}

		// 异步 Hook：直接 goroutine
		if h.Mode == Async {
			go func(hook *Hook) {
				if err := run(); err != nil {
					recordError(hook.Name, err)
				}
			}(h)
			continue
		}

		// 同步 Hook
		if err := run(); err != nil {
			recordError(h.Name, err)
			if h.MustSucceed {
				return fmt.Errorf("critical hook failed [%s]: %w", h.Name, err)
			}
		}
	}

	return nil
}
