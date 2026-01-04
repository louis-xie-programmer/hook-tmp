package main

import (
	"context"
	"log"
	"time"

	"hook-tmp/example"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("=== Hook Framework Demo Start ===")

	orderService := example.NewOrderService()

	// 正常订单
	log.Println("\n--- Case 1: normal order ---")
	err := orderService.CreateOrder(context.Background(), &example.Order{
		ID:     "order_001",
		UserID: "user_1001",
		Amount: 500,
	})
	if err != nil {
		log.Printf("create order failed: %v", err)
	}

	// 触发风控失败（MustSucceed Hook）
	log.Println("\n--- Case 2: risk rejected order ---")
	err = orderService.CreateOrder(context.Background(), &example.Order{
		ID:     "order_002",
		UserID: "user_1002",
		Amount: 20_000,
	})
	if err != nil {
		log.Printf("create order failed: %v", err)
	}

	// 等待异步 Hook 执行完成
	log.Println("\n--- waiting async hooks ---")
	time.Sleep(1 * time.Second)

	log.Println("=== Hook Framework Demo End ===")
}
