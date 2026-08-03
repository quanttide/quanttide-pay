# quanttide-pay
量潮支付工程

## 概述

量潮支付工程（quanttide-pay）是量潮知识管理体系中的**支付工程**领域，涵盖支付流程、账务处理、对账结算、风控合规等核心能力。

## 领域边界

- **支付流程**：收银台、支付网关、路由策略
- **账务处理**：账户体系、会计分录、余额管理
- **对账结算**：交易对账、资金结算、差错处理
- **风控合规**：交易风控、反欺诈、监管合规

## 仓库结构

| 路径 | 类型 | 说明 |
|------|------|------|
| `apps/qtcloud-pay` | 子模块 | 支付云服务平台（独立仓库，见 [.gitmodules](.gitmodules)） |
| `data/context` | 子模块 | 支付工程语境 |
| `data/journal` | 子模块 | 支付工程日志 |
| `data/library` | 子模块 | 支付工程图书馆 |
| `data/intention` | 子模块 | 支付工程意图 |
| `data/insight` | 子模块 | 支付工程洞察 |
| `data/roadmap` | 子模块 | 支付工程路线图 |
| `data/report` | 子模块 | 支付工程报告 |
| `packages/quanttide-pay-toolkit` | 本仓库目录 | 支付工程共享工具集：Go 基础库（`packages/go/`）+ 契约测试 fixtures（`tests/`） |

子模块操作：`git submodule update --init --recursive`；子模块内部改动须先在子模块仓库内提交，再回主仓库更新指针（见 [CONTRIBUTING.md](CONTRIBUTING.md)）。

## 共享工具集

`packages/quanttide-pay-toolkit` 承载领域主干（会算钱、会定状态、会推导结果的纯逻辑）：

- **Go 实现**（`packages/go/`）：money / ledger / status / idempotency / order / billing / httpapi / middleware
- **契约测试**：`tests/fixtures/`（JSON）是契约唯一权威，各语言实现消费同一份断言一致；当前 Go runner 已就位，**多语言方向为 Dart / Rust**（暂缓，详见工具库 [ROADMAP](packages/quanttide-pay-toolkit/ROADMAP.md)）
- **状态文档**：工具库 [STATUS.md](packages/quanttide-pay-toolkit/STATUS.md)（交接状态与待办）

## 参与贡献

- 工作纪律与关键文件索引见 [AGENTS.md](AGENTS.md)
- 提交流程、子模块与文档纪律见 [CONTRIBUTING.md](CONTRIBUTING.md)

## 许可

[CC BY 4.0](LICENSE)
