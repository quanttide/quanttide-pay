# quanttide-pay-toolkit ROADMAP（工具库级）

工具库目标：**领域主干从 provider 抽出到工具库，以契约测试（共享测试向量）统一各语言实现，从而统一多端**。

主干边界：会算钱、会定状态、会推导结果的纯逻辑（值对象/枚举/纯函数）进工具库；实体、存储（gorm）、事务编排、渠道适配留在应用壳（provider）。

> 语言实现级进度以 [`packages/go/ROADMAP.md`](packages/go/ROADMAP.md) 为准；本文档描述工具库整体（多语言 + 契约测试）的路线。

## 布局

| 路径 | 说明 | 状态 |
|------|------|------|
| `packages/go/` | Go 实现（8 个 pkg + `contracttest` runner + 文档） | ✅ 活跃，测试全绿 |
| `packages/dart/` | Dart 实现（计划中） | ⏸ 未启（多语言对齐暂缓，见 G6） |
| `packages/rust/` | Rust 实现（计划中） | ⏸ 未启（多语言对齐暂缓，见 G6） |
| `tests/fixtures/` | 共享契约测试向量（JSON，契约唯一权威） | ✅ 3 份（money/status/ledger），billing 待补 |
| `pyproject.toml` | 遗留：Python 包配置留位（多语言方向已定为 Dart/Rust，待清理） | ⏸ 留位，无 `src/` 实现 |
| `README.md` | 工具库概览（含各语言使用示例） | ⚠️ 与现状不符：描述 Python 模块，实际未实现 |

## 当前阶段（2026-08-03）

| 任务 | 落点 | 状态 |
|------|------|------|
| 主干抽取 Go 侧：ledger.Balance（G1）、契约测试骨架（G2）、idempotency（G4）、billing.Calculate 主体（G3）、状态契约扩展（S1）、httpapi/middleware（S2） | `packages/go/pkg/*` | ✅ 已完成，详见 [go ROADMAP](packages/go/ROADMAP.md) |
| 契约测试骨架：fixtures 唯一权威 + Go runner | `tests/fixtures/` + `packages/go/contracttest` | ✅ 已完成（3 份 fixtures 全过） |

## 下一步

| # | 优先级 | 任务 | 落点 | 状态 |
|---|--------|------|------|------|
| G3 收尾 | P0 | billing fixtures（抵扣顺序/力度/余额不足/非法金额）+ `dev-guide/billing.md` 决策记录 | `tests/fixtures/billing.json` + Go runner 用例 + 文档 | 待办 |
| G5 | P1 | `pkg/order` 完整订单模型（`Order` 记录契约、`SettleDetail`、生命周期规则）——`OrderStatus` 已提前落 `pkg/status` | `packages/go/pkg/order` 扩展 + fixtures | 搁置（与 G3 收尾同批） |
| G6 | P2 | Dart/Rust 实现对齐：`packages/dart`、`packages/rust` 消费同一 fixtures（先实现 money/status/order 号，再扩 billing/ledger）；补各自 runner | `packages/dart/` + `packages/rust/`（test runner 留位） | 暂缓（fixtures/runner 结构已留位，需要时再启） |

## 不做

- ~~多语言强制对齐~~：G6 暂缓（方向为 Dart/Rust，不做 Python 对齐），不设各语言实现时间线；契约演进由 fixtures 驱动，各端自行跟进

## 核心约定

- fixtures 是契约**唯一权威**：新增/变更契约先改 fixtures，各语言实现消费同一份断言一致，禁止各端自行加状态/语义
- 金额一律整数分（int64），禁止浮点；状态取值英文小写字符串，可直接存库
- 纯逻辑进工具库，实体/存储/编排留在应用壳
