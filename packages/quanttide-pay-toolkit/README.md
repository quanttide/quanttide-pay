# quanttide-pay-toolkit

量潮支付工程工具箱 — 量潮知识管理体系中的**支付工程**共享工具集。

## 概述

`quanttide-pay-toolkit` 为量潮支付工程（quanttide-pay）下的各应用（支付网关、账务处理、对账结算等）提供通用基础能力，避免在各应用中重复实现，并统一金额、状态等核心概念的语义。

## 模块

| 模块 | 说明 |
|------|------|
| `Money` | 金额值对象：`Decimal` 精确表示 + ISO 4217 币种，自动规范化为两位小数（分） |
| `PaymentStatus` / `RefundStatus` | 支付单、退款单状态枚举 |
| `generate_order_no` | 支付订单号生成（前缀 + 时间戳 + 随机序列） |

## 安装

```bash
uv add quanttide-pay-toolkit
# 或
pip install quanttide-pay-toolkit
```

## 快速开始

```python
from quanttide_pay_toolkit import Money, PaymentStatus, generate_order_no

# 金额：精确十进制，自动规范化为两位小数（分）
amount = Money("12.34") + Money.from_cents(100)
assert amount == Money("13.34")
assert amount.cents == 1334

# 状态：取值可直接存入数据库、与渠道报文流转
assert PaymentStatus.SUCCEEDED.value == "succeeded"

# 订单号：PAY20260803143000 + 10 位随机数字
order_no = generate_order_no(prefix="PAY")
```

## Go 基础库

`packages/go/` 提供了与 Python 侧同构的 Go 实现（金额、状态、订单号），遵循 Go 社区标准布局（`pkg/` + `internal/` + `examples/`），详见其 [README](packages/go/README.md)。

## 开发

```bash
uv sync
uv run pytest
```

## 许可

[Apache License 2.0](LICENSE)
