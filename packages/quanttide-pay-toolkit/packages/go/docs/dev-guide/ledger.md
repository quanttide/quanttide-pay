# Ledger 设计决策记录

`pkg/ledger` 的设计决策与演进记录（按时间序）。本文件即该包的正式决策记录，与 [user-guide/ledger.md](../user-guide/ledger.md)（契约）配套。

## D1 账本交易类型契约化（2026-08-03）

**背景**：qtcloud-pay provider 的 `internal/transaction/model.go` 定义了交易类型常量（recharge/refund/consume/issue/redeem）与余额影响判定（`AffectsBalance`），但这些是账务处理、对账结算等应用的共享语义，锁在单个应用内部。

**决策**：提炼为工具库 `pkg/ledger`——类型枚举 + `AffectsBalance`/`SignedAmount`/`IsValid` 三个语义函数；provider 的 `transaction` 模型常量改为引用工具库（保留本地别名 `TypeRecharge` 等，内部 API 不变）。

**理由**：

- 与 `pkg/money` 同一范式：共享领域契约进工具库，应用内只留薄适配
- 命名用 `ledger` 而非 `transaction`：避免与 provider 的 `internal/transaction` 模块撞名，且"账本类型"比"交易"更准确地表达余额语义

**影响**：余额求和与对账逻辑统一走 `AffectsBalance`/`SignedAmount`，禁止散落字符串比较。

## 未纳入

- 交易金额、幂等键等不属于类型契约，分别由 `pkg/money`、`pkg/idempotency`（另见 data/report 决策）负责
