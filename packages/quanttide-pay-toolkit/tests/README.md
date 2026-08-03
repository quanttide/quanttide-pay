# 契约测试

共享测试向量（fixtures）是各语言实现必须满足的**契约唯一权威**：同一输入 → 同一输出 / 同一拒绝行为。各语言消费同一 fixtures 并断言一致，即完成多端对齐。

## fixtures

| 文件 | 覆盖契约 |
|------|----------|
| `fixtures/money.json` | 金额 JSON 边界（严格整数分、币种校验、非零必带币种） |
| `fixtures/status.json` | 渠道码 → 统一状态映射（微信/支付宝/退款），未知码必须报错 |
| `fixtures/ledger.json` | 交易类型余额语义（AffectsBalance/SignedAmount）与余额推导（Balance） |

## runner

- Go：`packages/go/contracttest`——`go test ./contracttest`（从 Go 模块根运行）
- 新增语言（如 `packages/dart`、`packages/rust`）：实现同一契约，消费同一 fixtures 断言输出一致
