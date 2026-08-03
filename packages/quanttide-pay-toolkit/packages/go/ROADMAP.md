# quanttide-pay-toolkit Go ROADMAP

工具库目标：**领域主干从 provider 抽出到工具库，以契约测试（共享测试向量）统一各语言实现，从而统一多端**。

主干边界：会算钱、会定状态、会推导结果的纯逻辑（值对象/枚举/纯函数）进工具库；实体、存储（gorm）、事务编排、渠道适配留在应用壳（provider）。

## 当前阶段：G1 / G2 / G4 / G3 主体 / 状态与公共件扩展（已完成 2026-08-03）

| # | 优先级 | 任务 | 落点 | 状态 |
|---|--------|------|------|------|
| G1 | P0 | `ledger.Balance`：余额推导纯函数（余额 = Σ带符号交易）——消除「SQL CASE WHEN 聚合」与「SignedAmount 循环」两处实现的漂移风险 | `pkg/ledger`（`ledger.go` 或 `transaction.go`）+ 测试；`docs/dev-guide/ledger.md` 补 D3 决策记录；provider `transaction/gorm` 的 `SumByAccount` 加与 `ledger.Balance` 的等价性测试 | ✅ 已完成（含 `Transaction.Contract()` 契约视图） |
| G2 | P0 | 契约测试骨架：共享测试向量（JSON fixtures）+ Go runner——fixtures 是契约唯一权威，各语言消费同一份 | 工具库根 `tests/`（`fixtures/` + Go runner，pytest 驱动留位）；先覆盖已抽的 money / status / ledger | ✅ 已完成（`tests/fixtures/` + `contracttest`，三份 fixtures 全过） |
| G4 | P0 | `pkg/idempotency` 幂等键契约实施：键构造 `Key(biz, bizNo)`、业务前缀常量、边界校验（业务号含 `:` 拒绝） | 新 `pkg/idempotency` + 测试；provider account/coupon/voucher/order 四个 service 改引用（决策已定：data/report `idempotency-key-contract.md`） | ✅ 已完成（含 reconciliation 对账匹配键、conventions.md 引用、dev-guide D1） |
| G3 | P1 | 抽 `billing.Calculate` 计费计算主干：抵扣顺序（满减→折扣→代金券→余额）、力度选择（门槛内最大减额 / 最低折扣率）、代金券逐张 min、余额不足/非法金额拒绝；`Deduction`/`CouponInput`/`VoucherInput` 输入输出契约 | `pkg/billing`（纯函数，无存储依赖）+ 测试；provider `internal/billing` 改为委托引用 | ✅ 主体完成（`1782784`，provider `Service.Calculate` 已委托工具库）；**待收尾**：billing fixtures + `docs/dev-guide/billing.md`（见下） |
| S1 | P1 | 状态契约扩展：`CouponStatus`（issued/used/expired）、`OrderStatus`（created/settled）及 `IsValid*` 存库前校验——与支付/退款状态同包统一（D4：状态契约独立为共享模块） | `pkg/status`（`coupon.go`/`order.go`）+ 测试；provider coupon/order model 改引用 | ✅ 已完成（`c4d9753`） |
| S2 | P2 | HTTP 公共件：`httpapi` 统一响应/服务错误映射/分页 + `middleware` 请求日志（应用壳侧共用，非领域主干） | `pkg/httpapi` + `pkg/middleware` + 测试；provider 各 transport 改引用 | ✅ 已完成（`c2b6910`） |

## G3 收尾（可随 G5 同批实施）

- billing fixtures：覆盖抵扣顺序、力度选择（满减门槛/折扣率）、代金券逐张、余额不足/非法金额拒绝，纳入契约测试骨架（G2 扩展）
- `docs/dev-guide/billing.md` 决策记录：抵扣顺序与力度规则为系统级契约 v0.1（规则引擎后置），金额均为分（int64）

## 搁置：G5（与 G3 收尾同批实施）

| # | 优先级 | 任务 | 落点 | 状态 |
|---|--------|------|------|------|
| G5 | P1 | `pkg/order` 完整订单模型：`Order` 记录契约（无 gorm）＋ `SettleDetail` 结算明细契约（引用 `pkg/billing` 的 `Deduction`，形状同源）＋ 生命周期规则（可结算判定/状态迁移）；`OrderStatus` 已提前落入 `pkg/status`（D4：状态不散落业务域），账户号 `acc_` 并入 ID 规则待定 | `pkg/order` 扩展 + fixtures；`docs/dev-guide/order.md` 决策记录；provider `internal/order` 改引用（gorm 模型 = 契约 + 存储约束） | 搁置 |

## 不做：G6（多语言对齐取消）

- ~~G6 Python 库契约对齐~~：多语言对齐暂不做（方向已定为 Dart/Rust）；契约测试骨架（G2）保留通用 runner 结构，后续需要时再启

## 待定决策

- **`SettleDetail` 与 `Deduction` 的归属**：`pkg/billing` 与 `pkg/ledger` 的归属已随 G3 落定——`Kind`/`Deduction` 与计算逻辑同包抽入 `pkg/billing`；`SettleDetail`（订单结算快照）与 `Deduction` 同源，G5 实施时确认是复用 `billing.Deduction` 还是独立形状
- **订单状态机演进**：`OrderStatus`（created/settled）已落入 `pkg/status`，与 `PaymentStatus` 严格区分（结算生命周期 ≠ 支付生命周期，`created` 同名同值不合并）；v0.2.0 渠道接入后订单状态必然扩展（支付/退款状态回写），契约演进由 fixtures 驱动两端同步变更，不得各端自行加状态
- **账户契约**：不预建 `pkg/account`——`Account` 实体是应用壳（存储/行锁/幂等登记），主干只有推导规则（G1）与 ID 规则（G5）；若多端出现真实账户概念消费者再提炼无 gorm 的账户契约
