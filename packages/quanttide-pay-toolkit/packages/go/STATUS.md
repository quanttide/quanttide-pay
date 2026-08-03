# quanttide-pay-toolkit Go 基础库 · 状态交接

> 交接日期：2026-08-03 ｜ 最近提交：`b5b2eba`（ROADMAP 对齐）｜ 分支：`main`（领先 `origin/main` 3 个提交）
> 仓库：`packages/quanttide-pay-toolkit/packages/go`（monorepo `quanttide-pay` 内）

## 1. 总体状态

- **模块**：`github.com/quanttide/quanttide-pay-toolkit/packages/go`，Go 1.26，唯一依赖 `go-money v1.0.15`（零传递依赖）
- **定位**：支付领域主干（会算钱、会定状态、会推导结果的纯逻辑）抽离到工具库；实体、gorm 存储、事务编排、渠道适配留在应用壳（`apps/qtcloud-pay`）
- **git 状态**：工具库全部改动已提交（最近 3 个：`c4d9753` 状态契约、`c2b6910` httpapi/middleware、`b5b2eba` ROADMAP）；**未推送**，`apps/qtcloud-pay` 子模块有大量改动未提交（见 §6）
- **验证基线**：`go test ./...` 与 `go vet ./...` 全绿（契约测试、8 个 pkg 测试均通过）；`make lint` 需先安装 golangci-lint，未执行

## 2. 包清单（`pkg/`）

| 包 | 职责 | 关键 API / 常量 | 状态 |
|----|------|-----------------|------|
| `money` | 金额值对象（go-money 薄契约层）：整数分 + ISO 4217，JSON 严格整数校验（拒绝小数/字符串/指数，非零必带币种） | `New` / `CentsOf` / `CNY` 等币种常量；注意勿用 `NewFromFloat` | ✅ 稳定，有 fixtures 契约 |
| `ledger` | 账本交易类型（recharge/refund/consume/issue/redeem）+ 余额推导 | `Type` / `AffectsBalance` / `SignedAmount` / `Balance` / `IsValid`；`Transaction` 记录契约（无 gorm），JSON 形状有测试锁定 | ✅ 稳定，有 fixtures 契约 |
| `status` | 状态契约（小写英文字符串，可直接存库） | `PaymentStatus`（7 值）/ `RefundStatus`（4 值）+ 微信/支付宝原始码解析（未知码报错不兜底）+ `IsValid*` 存库前校验；`CouponStatus`（issued/used/expired）、`OrderStatus`（created/settled） | ✅ 稳定，支付/退款有 fixtures 契约；coupon/order 为新增 |
| `idempotency` | 幂等键构造：`{biz}:{bizNo}`，业务号空/含 `:` 拒绝（防键空间污染） | `Key` / `SettleRedeemKey` / 前缀常量 `Recharge`/`Refund`/`Issue`（含 `IssueCoupon`/`IssueVoucher`）/`Settle` | ✅ 稳定 |
| `order` | 支付订单号生成：前缀 + `yyyyMMddHHmmss` + 10 位密码学安全随机数 | `Generate(prefix, now)` / `GenerateNow(prefix)` | ✅ 稳定（**仅此功能**，完整订单模型见 ROADMAP G5） |
| `billing` | 结算抵扣计算（纯函数，无存储依赖）：满减→折扣→代金券→余额 | `Calculate` / `Deduction` / `CouponInput` / `VoucherInput` / `Kind*` / `ErrInsufficientBalance` / `ErrInvalidAmount` | ✅ 主体完成，**待收尾**：无 fixtures、无 `dev-guide/billing.md`（见 §5） |
| `httpapi` | HTTP JSON API 公共件 | `WriteJSON` / `WriteError` / `WriteServiceError`（`Mapper`）/ `ParsePagination`（limit 默认 20、上限 100） | ✅ 稳定（新增） |
| `middleware` | HTTP 通用中间件 | `Logging`（方法/路径/状态码/耗时） | ✅ 稳定（新增） |

## 3. 契约测试（`contracttest` + 工具库根 `tests/fixtures/`）

- **架构**：fixtures（JSON）是契约唯一权威，各语言实现消费同一份即完成多端对齐；`contracttest` 是 Go runner，断言 Go 实现与契约一致
- **fixtures 现状**（3 份，均全过）：
  - `money.json`：金额 JSON 边界（严格整数分、币种校验、非零必带币种）
  - `status.json`：微信/支付宝/退款渠道码 → 统一状态映射，未知码必须报错
  - `ledger.json`：交易类型余额语义 + 余额推导
- **新增契约的流程约定**：新抽主干（如 billing）须同时补 fixtures 与 runner 用例，保证跨端同步

## 4. 文档

| 文档 | 内容 |
|------|------|
| `README.md` | 布局表、使用示例、开发命令（`make test/vet/fmt/lint`） |
| `ROADMAP.md` | 主干抽取路线图与待定决策（进度见 §5） |
| `docs/user-guide/money.md`、`ledger.md` | 使用契约 |
| `docs/dev-guide/money.md`、`ledger.md`、`status.md`、`idempotency.md` | 设计决策记录（D1–D5、放弃方案、实现备忘）——**新增契约请先读**，尤其 `status.md` D4（状态不散落业务域） |

## 5. ROADMAP 进度摘要

- **已完成**：G1（`ledger.Balance`）、G2（契约测试骨架）、G4（`pkg/idempotency`）、G3 主体（`pkg/billing.Calculate`，provider 已委托引用）、S1（coupon/order 状态入 `pkg/status`）、S2（httpapi/middleware）
- **G3 收尾（未做）**：billing fixtures + `docs/dev-guide/billing.md` 决策记录
- **搁置**：G5（`pkg/order` 完整模型：`Order` 记录契约、`SettleDetail`（与 `billing.Deduction` 形状同源待确认）、生命周期规则）——`OrderStatus` 已提前落 `pkg/status`，范围收窄
- **不做**：G6（多语言对齐暂缓，fixtures runner 结构留位）
- **待定决策**：`SettleDetail` 与 `Deduction` 归属；订单状态机演进（由 fixtures 驱动两端同步变更，不得各端自行加状态）；不预建 `pkg/account`

## 6. 已知问题与待办（接手优先事项）

1. **推送**：本地 `main` 领先 `origin/main` 3 个提交（`c4d9753`、`c2b6910`、`b5b2eba`），需 review 后推送
2. **子模块**：`apps/qtcloud-pay` 工作区有大量未提交改动（transport/model 改引用工具库 httpapi/money/status/idempotency/billing、删除本地 middleware/logging 与内部测试等），需在子模块仓库内提交并回主仓库更新指针
3. **G3 收尾**：billing 补 fixtures（抵扣顺序/力度/余额不足/非法金额）+ `dev-guide/billing.md` 决策记录
4. **`order` 包边界**：当前仅订单号生成；G5 实施前勿在 provider 侧堆积订单模型，避免与契约演进冲突
5. **`middleware.Logging` 局限**：`statusRecorder` 仅记录 WriteHeader 显式写入的状态码，隐式 200（未调用 WriteHeader）记录为默认 200，行为正确但 `Flush`/`Hijack` 未透传——如接入 WebSocket/流式响应需先扩展
6. **文档与代码不符**：ROADMAP G1 声称 `Transaction.Contract()` 契约视图已完成，但代码中不存在（实际为 `TestTransactionJSONShape`）——补实现或改文档二选一，并回查 ROADMAP 其余条目

## 7. 核心约定（改代码前必读）

- 金额一律整数分（int64），全链路（内部/传输/对账）禁止浮点；JSON 形状 `{"amount": 分, "currency": "CNY"}`
- 状态取值英文小写字符串，可直接存库；渠道原始码解析未知码**必须报错**，不用 UNKNOWN 兜底
- 领域主干只放纯逻辑（值对象/枚举/纯函数），实体与存储留在应用壳；契约变更走 fixtures 驱动，不静默改语义
- 新包/新契约须同步：包文档注释（含 JSON 形状）、测试、fixtures（若跨端）、dev-guide 决策记录、README 布局表、ROADMAP 状态
