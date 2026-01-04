# Go Hook Framework Demo

这是一个轻量级的 Go 语言 Hook 框架示例，展示了如何通过钩子机制实现业务逻辑的解耦和扩展。该项目演示了如何在核心业务流程（如订单创建）的前后插入自定义逻辑（如风控、通知、积分等），而无需修改核心代码。

详细的内容介绍全在微信公众号中。干货持续更新，敬请关注「代码扳手」微信公众号：

<img width="430" height="430" alt="image" src="wx.jpg" />


## ✨ 核心特性

- **多阶段拦截**：支持 `Before` (前置) 和 `After` (后置) 两个执行阶段。
- **灵活的执行模式**：
  - **同步 (Sync)**：阻塞执行，适合参数校验、权限控制、风控等强依赖场景。
  - **异步 (Async)**：并发执行，适合发送短信、邮件、埋点、积分发放等非核心链路，不增加主流程耗时。
- **优先级控制**：支持通过 `Priority` 数值控制执行顺序（数值越大优先级越高）。
- **关键路径保护**：支持 `MustSucceed` 标志。若关键的同步 Hook 执行失败，将直接阻断业务流程并返回错误。
- **健壮性设计**：内置 Panic 捕获机制，防止 Hook 内部崩溃导致主程序退出。
- **可观测性**：内置基础的耗时记录、错误日志和 Panic 记录。

## 📂 目录结构

```text
├── hook/           # 核心框架代码
│   ├── engine.go   # 引擎核心逻辑 (注册、执行、排序、并发控制)
│   ├── hook.go     # Hook 结构体定义及枚举
│   ├── metrics.go  # 简单的监控日志实现
│   └── context.go  # 业务上下文数据结构定义
├── example/        # 业务使用示例
│   └── order.go    # 订单服务示例 (展示如何集成 Hook)
├── main.go         # 程序入口 (包含正常下单和风控拦截的测试用例)
└── README.md       # 项目说明文档
```

## 🚀 快速开始

### 1. 定义 Hook

Hook 结构体包含名称、优先级、模式（同步/异步）以及具体的执行函数。

```go
// 示例：定义一个风控 Hook
riskHook := &hook.Hook{
    Name:        "RiskCheck",
    Priority:    100,
    MustSucceed: true, // 如果失败，阻断流程
    Mode:        hook.Sync,
    Fn: func(ctx context.Context, data any) error {
        // 在这里编写具体的业务逻辑
        return nil
    },
}
```

### 2. 初始化引擎并注册

通常在服务初始化阶段（如 `NewService`）将 Hook 注册到引擎中。

```go
engine := hook.NewEngine()
engine.Register(hook.Before, riskHook)
```

### 3. 在业务流程中埋点

在业务逻辑的关键节点调用 `engine.Execute`。

```go
// 1. 执行前置 Hook (如：参数校验、风控)
if err := s.engine.Execute(ctx, hook.Before, data); err != nil {
    return err
}

// ... 执行核心业务逻辑 ...

// 2. 执行后置 Hook (如：发通知、加积分)
// 通常后置 Hook 的错误不影响主流程，或者由 Hook 内部处理
_ = s.engine.Execute(ctx, hook.After, data)
```

## 💡 示例运行

本项目包含一个订单服务的完整示例 (`example/order.go`)，模拟了以下流程：

1.  **创建订单前 (Before)**：执行风控检查 (Sync, MustSucceed)。如果金额 > 10,000 则拒绝。
2.  **核心逻辑**：保存订单信息。
3.  **创建订单后 (After)**：
    *   发送短信通知 (Async, Priority 10)
    *   增加用户积分 (Async, Priority 20)

### 运行代码

```bash
go run main.go
```

### 预期输出

```text
=== Hook Framework Demo Start ===

--- Case 1: normal order ---
[METRIC] hook=RiskCheck cost=0s
order order_001 saved
[METRIC] hook=AddPoints cost=0s
add points for user user_1001
[METRIC] hook=SendSms cost=0s
send sms for order order_001

--- Case 2: risk rejected order ---
[METRIC] hook=RiskCheck cost=0s
[ERROR] hook=RiskCheck err=risk rejected
create order failed: critical hook failed [RiskCheck]: risk rejected

--- waiting async hooks ---
=== Hook Framework Demo End ===
```

## 🛠️ 设计细节

*   **Engine**: 负责管理所有的 Hook，内部使用 `map[Type][]*Hook` 存储，并在注册时自动按优先级排序。
*   **Context**: 框架设计为通用型，使用 `any` 类型传递业务数据 (`data`)。在具体实现中（如 `hook/context.go`），定义具体的结构体来承载业务上下文。
*   **Concurrency**: 异步 Hook 使用 `go` 关键字启动协程执行，利用闭包捕获上下文，确保并发安全。
