package example

import (
	"context"
	"errors"
	"hook-tmp/hook"
	"log"
)

type Order struct {
	ID     string
	UserID string
	Amount int64
}

type OrderService struct {
	engine *hook.Engine
}

func NewOrderService() *OrderService {
	engine := hook.NewEngine()

	// 注册 Hook（启动阶段）
	engine.Register(hook.Before, &hook.Hook{
		Name:        "RiskCheck",
		Priority:    100,
		MustSucceed: true,
		Mode:        hook.Sync,
		Fn: func(ctx context.Context, data any) error {
			o := data.(*hook.OrderContext)
			if o.Amount > 10_000 {
				return errors.New("risk rejected")
			}
			return nil
		},
	})

	engine.Register(hook.After, &hook.Hook{
		Name:     "SendSms",
		Priority: 10,
		Mode:     hook.Async,
		Fn: func(ctx context.Context, data any) error {
			o := data.(*hook.OrderContext)
			log.Printf("send sms for order %s", o.OrderID)
			return nil
		},
	})

	engine.Register(hook.After, &hook.Hook{
		Name:     "AddPoints",
		Priority: 20,
		Mode:     hook.Async,
		Fn: func(ctx context.Context, data any) error {
			o := data.(*hook.OrderContext)
			log.Printf("add points for user %s", o.UserID)
			return nil
		},
	})

	return &OrderService{engine: engine}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error {
	octx := &hook.OrderContext{
		OrderID:  order.ID,
		UserID:   order.UserID,
		Amount:   order.Amount,
		Metadata: map[string]any{},
	}

	// 1. 创建前 Hook（校验 / 风控）
	if err := s.engine.Execute(ctx, hook.Before, octx); err != nil {
		return err
	}

	// 2. 核心业务逻辑（稳定、不轻易改动）
	log.Printf("order %s saved", order.ID)

	// 3. 创建后 Hook（通知 / 积分 / 埋点）
	_ = s.engine.Execute(ctx, hook.After, octx)

	return nil
}
