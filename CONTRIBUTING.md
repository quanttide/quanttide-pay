# 贡献指南（quanttide-pay）

## 仓库结构须知

本仓库是支付工程领域的**编排仓库**：大部分内容在子模块（`apps/`、`data/`）与普通目录（`packages/quanttide-pay-toolkit`）中。

- **子模块**（`apps/qtcloud-pay`、`data/*`）：各自独立仓库，在主仓库**只提交指针更新**
- **工具库**（`packages/quanttide-pay-toolkit`）：本仓库直接跟踪，Go 基础库 + 契约测试 fixtures

## 提交规范

遵循 Conventional Commits，描述用中文，参照历史：

```
feat(toolkit): pkg/status 新增优惠券/订单状态契约与存库前校验测试
docs(toolkit): ROADMAP 对齐当前进度
fix(provider): 渠道码解析未知码不再静默降级
```

- 类型：`feat` / `fix` / `refactor` / `docs` / `test` / `chore`
- 作用域：`toolkit`（工具库）、`provider`（支付云服务）、`docs` 等
- 一次提交一个逻辑变更；文档与代码同批提交

## 子模块纪律

1. 先进入子模块仓库提交并推送，再回主仓库 `git add <子模块路径>` 更新指针
2. 不把子模块内容改作他用；主仓库 diff 中子模块条目只应是指针变化
3. `git submodule update --init --recursive` 初始化后开始工作

## 工具库契约纪律（重点）

工具库以 **fixtures 为契约唯一权威**（`packages/quanttide-pay-toolkit/tests/fixtures/`）：

1. 新增/变更契约（金额、状态、账本、计费等）**先改 fixtures 与 runner**，再改实现；各语言实现消费同一份 fixtures
2. 金额一律整数分（int64），禁止浮点；JSON 形状 `{"amount": 分, "currency": "CNY"}`
3. 状态取值英文小写字符串，可直接存库；渠道原始码解析**未知码必须报错**，不用 UNKNOWN 兜底
4. 纯逻辑（值对象/枚举/纯函数）进工具库；实体、gorm 存储、事务编排、渠道适配留在应用壳（provider）
5. 多语言方向为 **Dart / Rust**（暂缓），不做 Python 对齐

完整纪律见工具库 [AGENTS.md](packages/quanttide-pay-toolkit/AGENTS.md) 与 [ROADMAP.md](packages/quanttide-pay-toolkit/ROADMAP.md)。

## 文档同步

改代码必须同步文档，缺一不可（按仓库惯例）：

| 变更 | 需同步 |
|------|--------|
| 工具库新包/新契约 | 包文档注释、测试、fixtures（若跨端）、`dev-guide/` 决策记录、go `README.md` 布局表、go `ROADMAP.md` 状态 |
| 交接状态变化 | 工具库 `STATUS.md` 与 go `STATUS.md`（待办、验证基线） |
| 多语言/路线调整 | 工具库 `ROADMAP.md`、根 `README.md` |

## 验证

- Go 基础库：`cd packages/quanttide-pay-toolkit/packages/go && make test && make vet`（或 `go test ./... && go vet ./...`）
- 契约测试：`go test ./contracttest`（从 Go 模块根运行，消费 `tests/fixtures/`）
- 推送前确认本地无未提交的文档/代码遗漏

## 不做的事

- 不在主仓库直接修改子模块内容（先提交到子模块仓库）
- 不跳过 fixtures 直接改契约语义；不静默新增状态/渠道码
- 不向工具库引入存储依赖（gorm 等）；工具库保持纯逻辑
