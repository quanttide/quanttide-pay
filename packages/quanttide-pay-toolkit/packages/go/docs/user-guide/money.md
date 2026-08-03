# Money 金额值对象使用指南

本指南是金额表示的**固定契约**：`pkg/money` 的实现必须与本文所述行为一致；金额相关改动先改本文，再改实现。设计决策演进见[开发指南](dev-guide/money.md)。

## 金额表示约定（全链路整数分）

- 金额一律以货币最小单位（分）的 int64 存储、计算与**传输**，禁止浮点金额
- 内部（账本、余额、交易、券、订单）与 API 边界单位一致——不存在元↔分换算；**元仅存在于展示层**（`Display()` 或前端格式化）
- 金额值对象基于 go-money（Go 社区事实标准），完整 API 直接复用：构造、加减、比较、Split/Allocate 分账、Display 格式化

## JSON 转换契约

传输格式：`{"amount": <整数分>, "currency": "<ISO 4217 币种>"}`

序列化（`MarshalJSON`）：输出整数分，与 go-money 默认一致，如 `New(9999, "CNY")` → `{"amount":9999,"currency":"CNY"}`。

反序列化（`UnmarshalJSON`）严格校验，**实现必须遵循**：

- **amount 必须为 JSON 整数**：拒绝小数（`99.99`）、字符串（`"9999"`）、指数记法（`1e3`）、前导加号（`+100`）
- **币种必须有效**：ISO 4217 代码（大小写不敏感，规范化存储为大写）；未知币种拒绝
- **非零金额必须携带币种**：如 `{"amount":100}` 拒绝
- **零金额允许空币种**：`{"amount":0,"currency":""}` 合法（零值 Money）
- **非 JSON 输入**：拒绝

**为什么覆盖 go-money 默认反序列化**：默认实现把 amount 按 float64 解析后 `int64()` 截断——非整数会静默舍入入账（如 `99.99` → 99 分），与"零误差"契约冲突。序列化无此问题，保持默认。

## 使用示例

```go
// 构造（最小单位：分）
m := money.New(9999, money.CNY) // 99.99 元

// 运算（币种不一致返回错误）
sum, err := m.Add(money.New(1, money.CNY))
diff, err := m.Subtract(money.New(1, money.CNY))

// 比较
ok, err := m.Equals(money.New(9999, money.CNY))

// 分账（余数 round-robin 分配，不丢分）
parts, err := m.Split(3)

// JSON 边界
b, _ := json.Marshal(m) // {"amount":9999,"currency":"CNY"}
var got money.Money
_ = json.Unmarshal(b, &got)

// 展示（元）：99.99 元
m.Display()

// 取分（nil 视为零金额）
cents := money.CentsOf(m) // 9999
```

## 约束

- **禁止使用 go-money 的 `NewFromFloat`**：浮点构造入口，可能引入舍入误差；本包刻意不转发
- 金额计算与存储一律整数运算，不经过浮点舍入
- 渠道层（第三方支付）对接时，金额一律以分传输，勿在边界做 `int(x*100)` 类浮点换算
