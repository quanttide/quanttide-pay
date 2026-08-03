# AGENTS（quanttide-pay-toolkit Go 基础库）

面向在 Go 包内工作的编码 agent 的指令。**动手前先读「关键文件」一节**；契约相关改动严格遵守「harness 工作纪律」。上级索引见仓库根 [AGENTS.md](../../../AGENTS.md)。

## 本包是什么

工具库的 Go 实现（`pkg/`：money/ledger/status/idempotency/order/billing/httpapi/middleware）+ 契约测试 runner（`contracttest/`）+ 设计文档（`docs/`）。fixtures 在工具库根 `tests/fixtures/`（本包相对路径上溯三级）。

## 关键文件（按优先级阅读）

| 文件 | 作用 | 何时必读 |
|------|------|----------|
| `README.md` | 布局表、使用示例、开发命令 | 每次工作前 |
| `ROADMAP.md` | 任务状态（G1–G6：已完成/搁置/待定决策） | 每次规划前 |
| `STATUS.md` | 交接状态、**待办优先事项**（billing fixtures、`Transaction.Contract()` 不符、middleware 局限、未推送） | 接手任务前 |
| `CONTRIBUTING.md` | 包纪律、契约测试纪律、文档同步清单 | 修改代码前 |
| `docs/dev-guide/*.md` | 设计决策记录（money D1–D5 / ledger D3 / status D4 / idempotency D1） | **新增/变更契约前必读** |
| `docs/user-guide/*.md` | 使用契约（money/ledger JSON 形状） | 涉及金额/账本序列化 |
| `contracttest/contract_test.go` | Go runner：断言实现与 fixtures 一致 | 改 runner 前 |
| `../../tests/README.md` | 契约测试架构说明 | 改 fixtures/runner 前 |
| `../../tests/fixtures/*.json` | **契约唯一权威**（money/status/ledger） | 涉及契约语义时 |
| `Makefile` | `make test/vet/fmt/lint` | 验证时 |

## harness 工作纪律（最高优先级，不可妥协）

harness = `tests/fixtures/`（JSON 契约向量） + `contracttest`（Go runner）。

1. **fixtures 先行**：新增/变更契约，先改 `tests/fixtures/*.json` 与 runner 用例，再改实现；实现与 fixtures 不一致 = 失败
2. **禁止端侧发明语义**：不得在本包或任一端自行新增状态/渠道码/字段含义——契约演进由 fixtures 驱动各端同步；D4：状态不散落业务域，放 `pkg/status` 不放业务包
3. **未知码必须报错**：渠道原始码（TradeState/trade_status/退款 status）解析遇未知值返回错误，**不用 UNKNOWN 兜底**
4. **金额整数分**：全链路 int64 分，禁止浮点；JSON 形状 `{"amount": 分, "currency": "CNY"}`，非零必带币种；勿用 go-money `NewFromFloat`
5. **状态小写字符串**：直接存库、与渠道报文流转；`IsValid*` 存库前防御校验
6. **主干边界**：纯逻辑（值对象/枚举/纯函数）进 `pkg/`；实体、gorm、事务编排、渠道适配留在 `apps/qtcloud-pay`；本包无存储依赖
7. **新增契约四同步**：包文档注释（含 JSON 形状）、单元测试、fixtures+runner（若跨端）、`docs/dev-guide/` 决策记录；另同步 README 布局表、ROADMAP 状态、STATUS 待办

## 已知状态（动手前核对 STATUS.md）

- `Transaction.Contract()`：ROADMAP 声称已完成但代码不存在（补实现或改文档）
- billing：主体完成，缺 fixtures 与 `dev-guide/billing.md`（G3 收尾）
- `middleware.Logging`：未透传 `Flush`/`Hijack`
- 本地 `main` 领先 `origin/main`，待推送

## 验证

```bash
go test ./... && go vet ./...   # 从模块根运行（含 contracttest）
```

提交前核对：测试/vet 通过、文档已同步、STATUS 待办已消化。
