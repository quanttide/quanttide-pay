# quanttide-pay-toolkit Go 基础库

量潮支付工程工具箱的 Go 基础库，金额值对象基于 [go-money](https://github.com/Rhymond/go-money)（零传递依赖），提供支付领域通用基础能力。遵循 Go 社区标准布局：`pkg/` 为可被外部导入的公共包，`internal/` 为库内部实现，`examples/` 为使用示例。

## 布局

| 路径 | 说明 |
|------|------|
| `pkg/money` | 金额值对象（go-money 薄封装）：整数分 + ISO 4217 币种，JSON 为 `{"amount": 整数分, "currency": 币种}`（严格整数校验） |
| `pkg/ledger` | 账本交易类型契约（recharge/refund/consume/issue/redeem）与余额影响语义；`Transaction` 交易记录契约（无存储绑定） |
| `pkg/status` | 状态契约：`PaymentStatus` / `RefundStatus` 支付、退款单状态（渠道码解析 + 存库前校验）；`CouponStatus`（issued/used/expired）、`OrderStatus`（created/settled）账本侧状态 |
| `pkg/idempotency` | 幂等键构造契约（`Key` / `SettleRedeemKey`，业务号边界校验防键空间污染） |
| `pkg/order` | 支付订单号生成（前缀 + 时间戳 + 密码学安全随机序列） |
| `pkg/billing` | 结算抵扣计算契约（纯函数，无存储依赖）：默认顺序「满减 → 折扣 → 代金券 → 余额」、力度选择、余额不足校验 |
| `pkg/httpapi` | HTTP JSON API 公共件：统一响应/错误体、服务错误映射（`Mapper`）、分页解析 |
| `pkg/middleware` | HTTP 通用中间件：请求日志（方法/路径/状态码/耗时） |
| `internal/` | 库的内部实现，外部不可导入 |
| `contracttest/` | 契约测试 Go runner（消费工具库根 `tests/fixtures/`，fixtures 是契约唯一权威） |
| `docs/` | 使用指南与契约文档（[`user-guide/money.md`](docs/user-guide/money.md) 为金额转换契约；[`dev-guide/money.md`](docs/dev-guide/money.md) 为设计决策记录） |
| `examples/` | 使用示例（`go run ./examples/payment`） |
| `test/` | 集成测试辅助数据/环境（预留，按需扩展） |

## 使用

```go
import (
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/money"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/order"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/status"
)

// 金额：整数分 + ISO 4217 币种，避免浮点误差
total, err := money.New(1234, money.CNY).Add(money.New(100, money.CNY))
// total.Amount() == 1334

// 状态：取值可直接存入数据库、与渠道报文流转
s := status.PaymentStatusSucceeded
// s == "succeeded"

// 订单号：PAY20260803143000 + 10 位随机数字
orderNo, err := order.GenerateNow("PAY")
```

## 契约测试

共享测试向量（JSON fixtures，工具库根 `tests/fixtures/`）是契约**唯一权威**，各语言实现消费同一份断言一致。Go runner 位于 `contracttest/`：

```bash
go test ./contracttest
```

## 文档

| 文档 | 内容 |
|------|------|
| [ROADMAP.md](ROADMAP.md) | 任务状态与待定决策（G1–G6） |
| [STATUS.md](STATUS.md) | 交接状态与待办优先事项 |
| [CONTRIBUTING.md](CONTRIBUTING.md) | 包纪律、契约测试纪律、文档同步清单 |
| [AGENTS.md](AGENTS.md) | agent 工作索引与 harness 纪律 |
| `docs/dev-guide/` | 设计决策记录（改契约前先读） |
| `docs/user-guide/` | 使用契约 |

## 开发

```bash
make test    # 单元测试
make vet     # 静态检查
make fmt     # 格式化
make lint    # golangci-lint（需先安装）
```

## 许可

[Apache License 2.0](../../LICENSE)
