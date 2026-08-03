# Ledger 设计决策记录

`pkg/ledger` 的设计决策与演进记录（按时间序）。本文件即该包的正式决策记录，与 [user-guide/ledger.md](../user-guide/ledger.md)（契约）配套。

## D1 账本交易类型契约化（2026-08-03）

**背景**：qtcloud-pay provider 的 `internal/transaction/model.go` 定义了交易类型常量（recharge/refund/consume/issue/redeem）与余额影响判定（`AffectsBalance`），但这些是账务处理、对账结算等应用的共享语义，锁在单个应用内部。

**决策**：提炼为工具库 `pkg/ledger`——类型枚举 + `AffectsBalance`/`SignedAmount`/`IsValid` 三个语义函数；provider 的 `transaction` 模型常量改为引用工具库（保留本地别名 `TypeRecharge` 等，内部 API 不变）。

**理由**：

- 与 `pkg/money` 同一范式：共享领域契约进工具库，应用内只留薄适配
- 命名用 `ledger` 而非 `transaction`：避免与 provider 的 `internal/transaction` 模块撞名，且"账本类型"比"交易"更准确地表达余额语义

**影响**：余额求和与对账逻辑统一走 `AffectsBalance`/`SignedAmount`，禁止散落字符串比较。

## D2 交易记录契约（2026-08-03）

**背景**：`internal/transaction` 的 `Transaction` 记录形状（JSON 字段、`Type` 类型）是账务/对账等应用读账本的统一形状，仍锁在 provider 应用内；`Type` 字段为裸 `string`，与 `ledger.Type` 契约脱节，跨应用流转时语义可能再次分叉。

**决策**：`pkg/ledger` 新增 `Transaction` **记录契约**（无 gorm tag，`Type` 强类型化为 `ledger.Type`，`IdempotencyKey` 不进 JSON）；provider 的 gorm 模型保留存储约束（唯一索引/长度），`Type` 字段改为 `ledger.Type`，本地类型常量改为类型化别名。

**理由**：

- 记录形状与类型语义同演进、强绑定（`Type` 字段直接引用 `ledger.Type`），放同一包零依赖；拆出独立 `pkg/transaction` 会产生账本相关包间依赖与命名摩擦，收益为负
- gorm 约束（`uniqueIndex`/`size`）是应用存储实现，工具库零传递依赖不进 gorm
- 与 `pkg/status` 同范式：契约进工具库，应用内保留带存储/框架约束的实现

**影响**：跨应用读账本以 `ledger.Transaction` 契约为准；gorm 模型与契约的显式转换在真实跨应用消费者出现时补充（当前 provider JSON 输出已与契约同形状）。

## D3 余额推导纯函数（2026-08-03）

**背景**：余额 = Σ带符号交易是全系统的不变式（对账、账单、一致性校验都依赖），但推导规则存在两处实现：`ledger.SignedAmount`（Go 循环）与 provider `transaction/gorm` 的 `SumByAccount`（SQL `CASE WHEN` 聚合）——两者必须永远等价，各自维护有漂移风险（新增影响余额的类型时漏改一处即错账）。

**决策**：`pkg/ledger` 新增 `Balance(txs []Transaction) int64`——Σ `SignedAmount`，是余额推导的**唯一权威规则**；provider 的 SQL 聚合保留作性能优化，但用等价性测试锁定「SQL 结果 == `ledger.Balance` 结果」（经 `Transaction.Contract()` 契约视图转换）。

**理由**：与「契约进工具库、实现留应用」同范式——推导规则跨应用共享且演进一致（多端统一对账/账单），必须单一权威；SQL 是同一语义的等价实现，等价性测试把漂移变成编译期可见的测试失败。

**影响**：新增影响余额的交易类型时，`AffectsBalance`/`SignedAmount`/`Balance` 一处修改，等价性测试强制 SQL 同步。

## 未纳入

- 交易金额、幂等键等不属于类型契约，分别由 `pkg/money`、`pkg/idempotency`（另见 data/report 决策）负责
