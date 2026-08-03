# quanttide-pay-toolkit Go ROADMAP

工具库目标：**领域主干从 provider 抽出到工具库，以契约测试（共享测试向量）统一各语言实现，从而统一多端**。

主干边界：会算钱、会定状态、会推导结果的纯逻辑（值对象/枚举/纯函数）进工具库；实体、存储（gorm）、事务编排、渠道适配留在应用壳（provider）。

## 当前阶段：G1 / G2 / G4（进行）

| # | 优先级 | 任务 | 落点 | 状态 |
|---|--------|------|------|------|
| G1 | P0 | `ledger.Balance`：余额推导纯函数（余额 = Σ带符号交易）——消除「SQL CASE WHEN 聚合」与「SignedAmount 循环」两处实现的漂移风险 | `pkg/ledger`（`ledger.go` 或 `transaction.go`）+ 测试；`docs/dev-guide/ledger.md` 补 D3 决策记录；provider `transaction/gorm` 的 `SumByAccount` 加与 `ledger.Balance` 的等价性测试 | 未开始 |
| G2 | P0 | 契约测试骨架：共享测试向量（JSON fixtures）+ Go runner——fixtures 是契约唯一权威，各语言消费同一份 | 工具库根 `tests/`（`fixtures/` + Go runner，pytest 驱动留位）；先覆盖已抽的 money / status / ledger | 未开始 |
| G4 | P0 | `pkg/idempotency` 幂等键契约实施：键构造 `Key(biz, bizNo)`、业务前缀常量、边界校验（业务号含 `:` 拒绝） | 新 `pkg/idempotency` + 测试；provider account/coupon/voucher/order 四个 service 改引用（决策已定：data/report `idempotency-key-contract.md`） | 未开始 |

## 搁置：G3 / G5（等待层边界仲裁，同批实施）

| # | 优先级 | 任务 | 落点 | 状态 |
|---|--------|------|------|------|
| G3 | P1 | 抽 `billing.Calculate` 计费计算主干：抵扣顺序（满减→折扣→代金券→余额）与全部券算法（折扣率/满减门槛/代金券 min） | 新 `pkg/billing`（纯函数 + 输入输出契约）+ fixtures；`docs/dev-guide/billing.md` 决策记录；provider `internal/billing` 改为引用，删本地实现 | 搁置 |
| G5 | P1 | `pkg/order` 完整订单模型：ID 生成（订单号已有 + 账户号 `acc_` 并入）＋ `Order` 记录契约（无 gorm）＋ `OrderStatus` 状态机（created/settled、合法流转、`IsValid`）＋ `SettleDetail` 结算明细契约 ＋ 生命周期规则（可结算判定/状态迁移） | `pkg/order` 扩展 + fixtures；`docs/dev-guide/order.md` 决策记录；provider `internal/order` 改引用（gorm 模型 = 契约 + 存储约束） | 搁置，**与 G3 同批实施**（`SettleDetail` 引用 `pkg/billing` 的 `Deduction`，形状同源） |

## 不做：G6（多语言对齐取消）

- ~~G6 Python 库契约对齐~~：多语言对齐暂不做；契约测试骨架（G2）保留通用 runner 结构，后续需要时再启

## 待定决策

- **`pkg/billing` 与 `pkg/ledger` 的归属**：计费计算（G3）是否同时把 `Kind` 枚举（coupon/voucher/balance）与 `Deduction` 输出契约一并抽入——倾向同包（与记录形状/语义同演进）；`SettleDetail`（订单结算快照）与 `Deduction` 同源，两包共享同一形状，G3/G5 同批实施时一并确认
- **订单状态机演进**：v0.2.0 渠道接入后订单状态必然扩展（支付/退款状态回写），契约演进由 fixtures 驱动两端同步变更，不得各端自行加状态；与 `pkg/status` 的 `PaymentStatus` 严格区分（结算生命周期 ≠ 支付生命周期，`created` 同名同值不合并）
- **账户契约**：不预建 `pkg/account`——`Account` 实体是应用壳（存储/行锁/幂等登记），主干只有推导规则（G1）与 ID 规则（G5）；若多端出现真实账户概念消费者再提炼无 gorm 的账户契约
