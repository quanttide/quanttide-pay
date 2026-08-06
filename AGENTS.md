# AGENTS（quanttide-pay）

面向在仓库内工作的编码 agent 的指令。**动手前先读「关键文件」一节**；涉及工具库契约时严格遵守「harness 工作纪律」。

## 仓库是什么

量潮支付工程编排仓库：`apps/qtcloud-pay`（支付云服务，子模块）、`data/*`（知识子模块）、`packages/quanttide-pay-toolkit`（共享工具集，子模块，独立仓库 `quanttide/quanttide-pay-toolkit`）。

## 关键文件（按优先级阅读）

| 文件 | 作用 | 何时必读 |
|------|------|----------|
| `README.md` | 仓库结构、工具库概述 | 每次工作前 |
| `CONTRIBUTING.md` | 提交规范、子模块纪律、文档同步要求 | 每次提交前 |
| `packages/quanttide-pay-toolkit/ROADMAP.md` | 工具库多语言路线与下一步（G3 收尾、G5、G6 Dart/Rust） | 涉及工具库规划 |
| `packages/quanttide-pay-toolkit/STATUS.md` | 工具库交接状态、**待办优先事项**（README 不符、pyproject 清理、Go 侧待办） | 接手工具库任务前 |
| `packages/quanttide-pay-toolkit/packages/go/ROADMAP.md` | Go 侧任务状态（已完成/搁置/待定决策） | 修改 Go 包前 |
| `packages/quanttide-pay-toolkit/packages/go/STATUS.md` | Go 侧交接状态与已知问题（`Transaction.Contract()` 不符、billing 收尾、middleware 局限） | 修改 Go 包前 |
| `packages/quanttide-pay-toolkit/packages/go/CONTRIBUTING.md` | Go 包贡献规范 | 修改 Go 包时 |
| `packages/quanttide-pay-toolkit/packages/go/AGENTS.md` | Go 包内特殊文件索引与契约纪律 | 修改 Go 包时 |
| `packages/quanttide-pay-toolkit/packages/go/docs/dev-guide/*.md` | 设计决策记录（D1–D5：money/ledger/status/idempotency） | **新增/变更契约前必读**（尤其 `status.md` D4：状态不散落业务域） |
| `packages/quanttide-pay-toolkit/packages/go/docs/user-guide/*.md` | 使用契约（money/ledger JSON 形状等） | 涉及金额/账本序列化 |
| `packages/quanttide-pay-toolkit/tests/README.md` | 契约测试架构说明 | 改 fixtures/runner 前 |
| `packages/quanttide-pay-toolkit/tests/fixtures/*.json` | **契约唯一权威**（money/status/ledger） | 涉及契约语义时 |

## harness 工作纪律（最高优先级）

harness = 契约测试骨架：`tests/fixtures/`（JSON 契约向量） + `contracttest`（Go runner）。以下纪律不可妥协：

1. **fixtures 先行**：新增/变更契约，先改 `tests/fixtures/*.json` 与 runner 用例，再改实现；实现与 fixtures 不一致 = 失败
2. **禁止端侧发明语义**：不得在任一端（Go 实现、provider、future Dart/Rust）自行新增状态/渠道码/字段含义——契约演进由 fixtures 驱动两端同步
3. **未知码必须报错**：渠道原始码（TradeState 等）解析遇未知值返回错误，**不用 UNKNOWN 兜底**（新渠道状态暴露而不是被掩盖）
4. **金额整数分**：全链路 int64 分，禁止浮点；JSON 形状 `{"amount": 分, "currency": "CNY"}`，非零必带币种
5. **状态小写字符串**：直接存库、与渠道报文流转；`IsValid*` 存库前防御校验
6. **主干边界**：纯逻辑（值对象/枚举/纯函数）进工具库；实体、gorm、事务编排、渠道适配留在 `apps/qtcloud-pay`
7. **新增契约四同步**：包文档注释（含 JSON 形状）、单元测试、fixtures+runner（若跨端）、`dev-guide/` 决策记录；另同步 README 布局表与 ROADMAP 状态

## 提交纪律

- 子模块（`apps/qtcloud-pay`、`data/*`、`packages/quanttide-pay-toolkit`）内容先提交到各自仓库，主仓库只记录指针
- Conventional Commits，中文描述（如 `feat(toolkit): ...`）；文档与代码同批提交
- 提交前核对：相关 `STATUS.md` 待办是否已消化、文档是否同步、`go test ./... && go vet ./...` 是否通过
