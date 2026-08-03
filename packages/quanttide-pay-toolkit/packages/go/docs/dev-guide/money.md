# Money 设计决策记录

`pkg/money` 的设计决策与演进记录（按时间序）。本文件即该包的正式决策记录，与 [user-guide/money.md](../user-guide/money.md)（契约）配套：契约固定实现行为，本文固定设计脉络。

## D1 全链路整数分（2026-08-03）

**背景**：v0.1.0 最初约定"内部 int64 分 + API 传输元（两位小数）"，由自研 `Cents` 类型承担边界转换，严格校验拒绝三位及以上小数。

**决策**：金额一律以货币最小单位（分）的 int64 存储、计算与**传输**——内部与 API 边界单位统一，不存在元↔分换算；元仅存在于展示层（`Display()`/前端）。

**理由**：

- 传输层与库默认（go-money 整数分）、行业惯例（Stripe、微信支付 `total_fee`）一致
- 对账 CSV（`amount_cents`）天然同单位，全链路无换算
- "三位小数静默舍入"整类问题没有存在空间（无小数）

**影响**：API 金额形状从 `99.99`（元数字）变为 `{"amount":9999,"currency":"CNY"}`（结构化，破坏性变更，允许）。

## D2 采用 go-money 为事实标准（2026-08-03）

**背景**：原计划自研完整的 Money 值对象（曾起草 `Money{分+币种}` + `internal/currency` 币种校验）。

**决策**：金额值对象基于 `github.com/Rhymond/go-money` v1.0.15（Go 社区事实标准：int64 最小单位 + ISO 4217 币种 + 完整运算/分账/格式化），不重复实现。

**理由**：自研草案被放弃——单币种（CNY）场景下币种感知的完整值对象设计是过度设计，且重复造轮子。

## D3 薄契约层而非完整封装（2026-08-03）

**背景**：采用 go-money 后，质疑"是否还有必要单独写一个 money 包"。

**决策**：保留**薄契约层**（`pkg/money`，真逻辑约 30 行），职责：

1. 类型别名与转发（`Money = gomoney.Money`、`New`、`GetCurrency`、常用币种常量、`CentsOf`）——零逻辑
2. 严格整数反序列化（见 D4）——唯一自定义逻辑
3. 约束文档化（禁止 `NewFromFloat`，刻意不转发）

**理由**：go-money 已覆盖全部金额能力；但"严格反序列化"是跨应用共享的契约，必须有一个权威位置，否则每个应用各自实现。删除项：元/小数位处理（scale 表、按币种小数位格式化）——随 D1 失去意义。

## D4 覆盖 go-money 默认 JSON 反序列化（2026-08-03）

**背景**：go-money 默认 `UnmarshalJSON` 把 amount 按 float64 解析后 `int64()` 截断——`{"amount":99.99}` 会静默变成 99 分入账。

**决策**：经 go-money 官方注入点（包级 `MarshalJSON`/`UnmarshalJSON` 变量）覆盖反序列化为严格整数校验：仅接受 JSON 整数（拒绝小数、字符串、指数、前导加号）；币种必须有效（`GetCurrency` 非 nil，大小写不敏感）；非零金额必须携带币种；零金额允许空币种。序列化保持 go-money 默认（本就是整数分输出）。

**理由**：静默截断与"零误差"契约冲突；注入点是 go-money 官方支持的定制机制。

**影响**：覆盖是包级全局状态（`init()` 中设置），工具库是唯一设置方，使用时需知晓。

## D5 提炼至工具库 pkg/money（2026-08-03）

**背景**：`pkg/money.Cents` 原在 qtcloud-pay provider 内，金额能力需要跨应用复用。

**决策**：提炼至 `quanttide-pay-toolkit` Go 模块，放 `pkg/`（公共）而非 `internal/`；provider 经 go.mod `replace` 本地依赖。

**理由**：Go 的 `internal/` 规则——`internal` 包只能被同模块导入，放在 `internal/` 则 provider（不同模块）永远无法引用，工具库失去意义。工具库 README 布局约定本就声明 `pkg/` 为公共包。

## 放弃的方案

| 方案 | 放弃原因 |
|------|----------|
| 自研 `Cents`（元传输 + 纯整数解析） | 传输层单位与内部不一致，换算类 bug 有存在空间；被 D1 取代 |
| 自研完整 Money 值对象（`Money{分+币种}` + `internal/currency`） | 重复造轮子；单币种场景过度设计 |
| 元 JSON（按币种小数位格式化，支持 JPY/BHD） | 与 go-money 默认、行业惯例不一致；D1 后无意义 |
| 完全不用封装（provider 直接用 go-money） | 严格反序列化契约会散落各应用重复实现，与工具库定位冲突（D3） |

## 实现细节备忘

- **`json.Marshal` 会 compact Marshaler 输出**：Go 标准库对 `MarshalJSON` 返回的字节做空白压缩——测试断言须用紧凑格式（`{"amount":9999,"currency":"CNY"}`，无空格）
- **go-money 未知币种静默降级**：`New(100, "XYZ")` 不报错，会生成默认格式的伪币种（`getDefault()`）——校验必须显式用 `GetCurrency(code) == nil` 判断
- **币种常量是无类型字符串常量**：`const CNY = "CNY"`（go-money constants.go），转发导出用 `const CNY = gomoney.CNY`
- **类型别名不转发包级函数**：`type Money = gomoney.Money` 只别名类型；`New` 等包级函数与常量需手动转发导出
