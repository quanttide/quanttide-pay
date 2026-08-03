# Ledger 账本交易类型契约

对应实现：[`pkg/ledger`](../pkg/ledger/ledger.go)。

本指南是账本交易类型的**固定契约**：`pkg/ledger` 的实现必须与本文所述行为一致。

## 交易类型

取值使用英文小写字符串，可直接存入数据库并与账务/对账应用流转：

| 类型 | 值 | 语义 | 影响余额 |
|------|-----|------|----------|
| `TypeRecharge` | `recharge` | 充值（对公打款入账） | ✅ |
| `TypeRefund` | `refund` | 退款（多退登记：对公退款出账） | ✅ |
| `TypeConsume` | `consume` | 消费（余额支付部分） | ✅ |
| `TypeIssue` | `issue` | 发券（信息性记录） | ❌ |
| `TypeRedeem` | `redeem` | 核销（券抵扣部分） | ❌ |

## 语义函数

- `AffectsBalance(t)`：是否影响余额——`recharge`/`refund`/`consume` 为 true，其余 false。**余额求和与对账必须以此为准**，不可自行枚举字符串
- `SignedAmount(t, amount)`：余额方向的带符号金额——充值 `+amount`，退款/消费 `−amount`，其余 `0`
- `IsValid(t)`：是否为已知类型（存库前防御校验）

## 约束

- 新增交易类型必须先更新本契约（含上表与各语义函数），再改实现
- 余额相关的求和/对账逻辑一律经 `AffectsBalance`/`SignedAmount`，禁止散落 `== "recharge"` 字符串比较
