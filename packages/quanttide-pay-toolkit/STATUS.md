# quanttide-pay-toolkit 状态交接（工具库级）

> 交接日期：2026-08-03 ｜ 最近提交：`bf29189`（本地 `main`，领先 `origin/main` 2 个提交）
> 仓库：`packages/quanttide-pay-toolkit`（monorepo `quanttide-pay` 内）

## 1. 总体状态

- **多语言工具集**，当前仅 **Go 实现**（`packages/go`）处于活跃开发；Python 侧只有 `pyproject.toml` 留位，**尚无 `src/` 实现与测试**
- **契约测试架构**：`tests/fixtures/`（JSON）是契约唯一权威，Go runner（`packages/go/contracttest`）消费并断言；多端对齐通过各语言消费同一份 fixtures 完成
- **git 状态**：工具库全部改动已提交；本地 `main` 领先 `origin/main` 2 个提交（`b5b2eba` ROADMAP 对齐、`bf29189` STATUS 交接）
- **验证基线**：`go test ./...` 与 `go vet ./...` 全绿（Go 模块内）；**Python 侧无测试可跑**（见 §4 待办 2）

## 2. 目录结构

| 路径 | 说明 |
|------|------|
| `README.md` | 工具库概览：Python 使用示例 + Go 库说明 ⚠️（见 §4 待办 1） |
| `pyproject.toml` | Python 包配置留位（0.0.1，pytest testpaths=`tests`） |
| `packages/go/` | Go 实现：`pkg/`（money/ledger/status/idempotency/order/billing/httpapi/middleware）、`contracttest/`、`docs/`（user-guide + dev-guide）、`examples/`；含自身 `ROADMAP.md` 与 `STATUS.md` |
| `packages/py/` | 计划中，未创建 |
| `tests/README.md` | 契约测试说明（fixtures 权威、runner 指引） |
| `tests/fixtures/` | `money.json` / `status.json` / `ledger.json`（billing 待补） |
| `LICENSE` | Apache 2.0 |

## 3. 模块状态

| 模块 | 实现 | fixtures | 测试 | 状态 |
|------|------|----------|------|------|
| 金额（money） | ✅ Go | ✅ | ✅ Go + 契约 | 稳定 |
| 状态（status：支付/退款/优惠券/订单） | ✅ Go | ⚠️ 仅支付/退款 | ✅ Go + 契约 | 稳定；coupon/order 状态为新增，fixtures 未覆盖 |
| 账本（ledger：类型/余额推导/交易记录） | ✅ Go | ✅ | ✅ Go + 契约 | 稳定 |
| 幂等（idempotency） | ✅ Go | ❌（本地单测覆盖） | ✅ Go | 稳定 |
| 订单号（order.Generate） | ✅ Go | ❌ | ✅ Go | 稳定；仅订单号生成，完整模型见 G5 |
| 计费（billing.Calculate） | ✅ Go | ❌ **待补** | ✅ Go | 主体完成，**待收尾**（fixtures + dev-guide） |
| HTTP 公共件（httpapi/middleware） | ✅ Go | — | ✅ Go | 稳定（新增） |
| Python 实现 | ❌ 未实现 | — | ❌ | 留位（G6 暂缓） |

## 4. 已知问题与待办（接手优先事项）

1. **根 README 与实际不符（首要）**：README 的「模块」「安装」「快速开始」描述 Python 模块（`Money` / `PaymentStatus` / `generate_order_no`），但仓库内无 `src/`、无 Python 测试——**实现 Python 包或修正 README 二选一**，当前文档会误导新成员
2. **pytest 配置与现状不符**：`pyproject.toml` 的 `testpaths = ["tests"]` 指向的目录只有 fixtures（无 `test_*.py`），运行 `uv run pytest` 会报 "no tests ran"（退出码 5）；Python 测试启用前先调整配置
3. **Go 侧待办**（详见 [`packages/go/STATUS.md`](packages/go/STATUS.md)）：
   - billing 补 fixtures 与 `dev-guide/billing.md` 决策记录（G3 收尾）
   - `Transaction.Contract()` 文档声称已完成、代码中不存在（补实现或改文档）
   - `middleware.Logging` 未透传 `Flush`/`Hijack`（接入流式响应前需扩展）
   - 本地 `main` 未推送
4. **历史重复提交**：本地历史含 `c4d9753` / `c2b6910`，`origin/main` 的 `60e5959` 为整合版本（同一内容 + README/子模块指针更新）——已推送以 `60e5959` 为准，历史中的重复提交不建议回写，避免与远端分叉
5. **子模块**：`apps/qtcloud-pay` 工作区仍有未提交改动（transport 改引用工具库等），需在子模块仓库内提交并回主仓库更新指针

## 5. 核心约定（改代码前必读）

- fixtures 是契约**唯一权威**：新增/变更契约先改 fixtures 与 runner，各语言同步；禁止各端自行加状态/语义
- 金额一律整数分（int64），JSON 形状 `{"amount": 分, "currency": "CNY"}`，禁止浮点
- 状态取值英文小写字符串，可直接存库；渠道原始码解析未知码必须报错，不用 UNKNOWN 兜底
- 纯逻辑（值对象/枚举/纯函数）进工具库；实体、gorm 存储、事务编排、渠道适配留在应用壳
- 新契约/新包须同步：包文档注释、测试、fixtures（若跨端）、dev-guide 决策记录、README 布局表、ROADMAP 状态
