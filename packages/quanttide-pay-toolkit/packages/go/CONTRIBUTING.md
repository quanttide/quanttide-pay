# 贡献指南（quanttide-pay-toolkit Go 基础库）

## 本包定位

Go 基础库承载**领域主干**：会算钱、会定状态、会推导结果的纯逻辑（值对象/枚举/纯函数）。实体、存储（gorm）、事务编排、渠道适配留在应用壳（`apps/qtcloud-pay`）。

## 包纪律

1. **纯逻辑**：包内不引入 gorm/数据库/http client 等存储与 IO 依赖；唯一运行依赖为 `go-money`（金额值对象）
2. **契约即文档**：每个包的文档注释写清契约语义与 JSON 形状；取值用英文小写字符串，可直接存库
3. **防御校验**：枚举/状态提供 `IsValid*` 供存库前校验；渠道原始码解析**未知码必须报错**，不用 UNKNOWN 兜底
4. **金额整数分**：全链路 int64 分，禁止浮点；JSON 形状 `{"amount": 分, "currency": "CNY"}`，非零必带币种

## 契约测试纪律（harness）

harness = `tests/fixtures/`（JSON 契约向量，契约唯一权威）+ `contracttest/`（Go runner）：

1. **fixtures 先行**：新增/变更契约，先改 `tests/fixtures/*.json` 与 runner 用例，再改实现；实现与 fixtures 不一致 = 失败
2. 现有 fixtures：`money.json` / `status.json` / `ledger.json`；**billing 待补**（ROADMAP G3 收尾）
3. 禁止各端自行新增状态/渠道码/字段含义——契约演进由 fixtures 驱动两端同步

## 文档同步（新增/变更契约时缺一不可）

| 项 | 要求 |
|----|------|
| 包文档注释 | 含契约语义与 JSON 形状 |
| 单元测试 | 与包同目录 `*_test.go` |
| fixtures + runner | 若契约跨端（金额/状态/账本/计费） |
| `docs/dev-guide/*.md` | 设计决策记录（新增契约前先读已有 D1–D5） |
| `README.md` | 布局表同步 |
| `ROADMAP.md` / `STATUS.md` | 任务状态与交接待办同步 |

## 开发命令

```bash
make test    # go test ./...（含 contracttest，从模块根运行）
make vet     # go vet ./...
make fmt     # gofmt -w .
make lint    # golangci-lint（需先安装）
```

契约测试单独运行：`go test ./contracttest`（消费 `../../tests/fixtures/`）。

## 提交规范

Conventional Commits，中文描述，作用域 `toolkit`：

```
feat(toolkit): pkg/status 新增优惠券/订单状态契约与存库前校验测试
docs(toolkit): ROADMAP 对齐当前进度
```

提交前核对：`go test ./... && go vet ./...` 通过、STATUS.md 待办已消化、文档已同步。
