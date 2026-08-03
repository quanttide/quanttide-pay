# quanttide-pay-toolkit Go 基础库

量潮支付工程工具箱的 Go 基础库，零外部依赖，提供支付领域通用基础能力。遵循 Go 社区标准布局：`pkg/` 为可被外部导入的公共包，`internal/` 为库内部实现，`examples/` 为使用示例。

## 布局

| 路径 | 说明 |
|------|------|
| `pkg/money` | 金额值对象（分 + ISO 4217 币种），支持加减与币种校验 |
| `pkg/status` | `PaymentStatus` / `RefundStatus` 支付单、退款单状态 |
| `pkg/order` | 支付订单号生成（前缀 + 时间戳 + 密码学安全随机序列） |
| `internal/` | 库的内部实现，外部不可导入 |
| `examples/` | 使用示例（`go run ./examples/payment`） |
| `test/` | 集成测试辅助数据/环境（按需扩展） |

## 使用

```go
import (
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/money"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/order"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/status"
)

// 金额：以分为单位，避免浮点误差
total, err := money.MustMoney(1234, "CNY").Add(money.MustMoney(100, "CNY"))
// total.Cents == 1334

// 状态：取值可直接存入数据库、与渠道报文流转
s := status.PaymentStatusSucceeded
// s == "succeeded"

// 订单号：PAY20260803143000 + 10 位随机数字
orderNo, err := order.GenerateNow("PAY")
```

## 开发

```bash
make test    # 单元测试
make vet     # 静态检查
make fmt     # 格式化
make lint    # golangci-lint（需先安装）
```

## 许可

[Apache License 2.0](../../LICENSE)
