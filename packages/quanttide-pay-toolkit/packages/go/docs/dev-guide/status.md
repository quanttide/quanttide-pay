# Status 设计决策记录

`pkg/status` 的设计决策与演进记录（按时间序）。本文件即该包的正式决策记录。

## D1 统一状态契约（2026-08-03）

**背景**：渠道层支付/退款状态表示不统一：

- 微信支付：透传原生 `TradeState`（`SUCCESS`/`REFUND`/`NOTPAY`/`USERPAYING`/`CLOSED`/`PAYERROR`…）
- 支付宝：本地 map 映射为大写枚举（`PENDING`/`SUCCESS`/`CLOSED`/`UNKNOWN`），未知码兜底 `UNKNOWN`
- 两套码并存、语义未固化；`UNKNOWN` 兜底会掩盖新状态出现的问题，错误语义可能静默入账

**决策**：支付单、退款单状态以 `pkg/status` 枚举为准（小写英文、可直接存库、与渠道报文流转一致）：

- `PaymentStatus`：created / paying / succeeded / failed / closed / refunding / refunded
- `RefundStatus`：created / processing / succeeded / failed

**理由**：渠道码与存储值分离——存库/流转统一用小写值，渠道差异收敛在适配器边界；支付网关、账务等后续应用直接复用同一状态契约。

**影响**：各应用不再各自解释渠道状态，状态值直接落库。

## D2 渠道码严格映射（2026-08-03）

**背景**：渠道原始码 → 统一状态的翻译此前在 provider 适配器本地维护（微信透传原生码、支付宝本地 map），没有权威位置。

**决策**：`pkg/status` 新增 `ParseWechatTradeState` / `ParseAlipayTradeStatus`——渠道原始码 → 统一状态；provider `channel/adapters.go` 删除本地映射表，改用工具库映射。

**理由**：翻译收敛到单一权威位置；新渠道码出现时只需改一处与对应测试。

**影响**：映射关系固化为表格（见下），作为解析函数与测试的双重依据。

### 微信 `TradeState` → `PaymentStatus`

| 渠道码 | 统一状态 |
|--------|----------|
| `SUCCESS` | succeeded |
| `REFUND` | refunding |
| `NOTPAY` | created |
| `USERPAYING` | paying |
| `CLOSED` | closed |
| `PAYERROR` | failed |
| 其他 | 错误 |

### 支付宝 `trade_status` → `PaymentStatus`

| 渠道码 | 统一状态 |
|--------|----------|
| `TRADE_SUCCESS` | succeeded |
| `TRADE_FINISHED` | succeeded |
| `WAIT_BUYER_PAY` | created |
| `TRADE_CLOSED` | closed |
| 其他 | 错误 |

### 微信退款 `status` → `RefundStatus`（补充）

决策补充：微信退款响应 `status` 字段的映射。微信退款终态 `CLOSED`/`ABNORMAL` 均归一为 `failed`（`RefundStatus` 无 closed 枚举）。

| 渠道码 | 统一状态 |
|--------|----------|
| `SUCCESS` | succeeded |
| `PROCESSING` | processing |
| `CLOSED` | failed |
| `ABNORMAL` | failed |
| 其他 | 错误 |

## D3 未知码显式报错（2026-08-03）

**背景**：支付宝适配器未知码兜底 `UNKNOWN`，新渠道状态出现时被静默掩盖，错误语义可能静默入账。

**决策**：映射失败返回错误，不用 `UNKNOWN` 兜底——与金额边界拒绝未知币种同一原则（见 money 决策）：新渠道状态出现时显式暴露，而非静默降级。

**理由**：静默降级把"渠道新增状态"变成不可见事件；显式报错迫使调用方先回答"这个新码到底意味着什么"再放行。

**影响**：`IsValidPaymentStatus` / `IsValidRefundStatus` 供存库前防御校验。

## D4 状态契约独立为共享模块，不散落业务域（2026-08-03）

**背景**：状态契约的归属存疑——单独模块，还是放在各业务域（order / transaction / billing / channel）？

**决策**：作为跨应用共享契约放工具库 `pkg/status`（与 `pkg/money`、`pkg/ledger` 同级），不放入任何业务域，也不放应用内 `internal/core`。

**理由**：

- 状态是聚合（支付单/退款单）的**值对象**与生命周期词汇，不是某个聚合的私有状态：`PaymentStatus` 被 channel（渠道翻译）、order（流转）、transaction/billing（对 `succeeded` 等事件反应）共用；`RefundStatus` 不属于 order——放入业务域会把同一契约劈成两半
- 各域各自定义 = 每处边界多一层翻译，语义漂移必然复发（本决策要解决的正是两处映射表漂移）
- 依赖方向：状态若放 order 域，channel 为取枚举需反向依赖 order，适配器依赖整个订单域，方向错误
- Go `internal/` 规则：放应用内 `internal/core` 则支付网关、账务等后续应用永远无法引用，与跨应用复用目标冲突（与 `pkg/money` 提炼至工具库同一原则）
- 状态机合法流转（如 succeeded 后才能 refund）与枚举同处一地，不变量有单一落点

**影响**：provider 经 go.mod `replace` 本地依赖；"枚举 + 解析 + 校验"三合一，"渠道状态变化"这个变化点集中为一个测试失败/编译错误，而不是散落的静默逻辑。

## 放弃的方案

| 方案 | 放弃原因 |
|------|----------|
| 各业务域自行定义状态（order/transaction/billing 各一份） | 同一契约被拆分，每处边界需要翻译层，语义漂移复发（D4） |
| 状态放 order 域、其他域引用 | `RefundStatus` 不属于 order；channel 反向依赖 order 域，依赖方向错误（D4） |
| 放应用内 `internal/core` | Go `internal/` 只允许同模块导入，网关、账务等后续应用无法复用（D4） |
| 适配器本地映射 + `UNKNOWN` 兜底 | 两套码并存、未知码静默降级，新渠道状态被掩盖（D3） |

## 实现细节备忘

- **包内按模块分文件**：`payment.go`（支付状态：枚举 + 微信/支付宝解析 + 校验）、`refund.go`（退款状态：枚举 + 微信退款解析 + 校验），对应测试同分 `payment_test.go`/`refund_test.go`
- **provider `channel/adapters.go` 已删除本地映射表**（`alipayTradeStatus`、`UNKNOWN` 兜底、微信透传），改用 `pkg/status` 解析函数（已实施）
- **`channel/model.go` 状态字段已类型化**：`OrderStatus.Status` → `status.PaymentStatus`、`RefundResponse.Status` → `status.RefundStatus`（已实施）
- **存库前防御校验**：`IsValidPaymentStatus` / `IsValidRefundStatus`（已实施）
- **新渠道码处理流程**：解析函数返回错误 → 测试先红 → 显式决定新码语义后再放行，禁止追加兜底分支
