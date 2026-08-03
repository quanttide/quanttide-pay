# Idempotency 设计决策记录

`pkg/idempotency` 的设计决策与演进记录（按时间序）。本文件即该包的正式决策记录。契约决策原文见 data/report `technical-decisions/idempotency-key-contract.md`。

## D1 幂等键构造契约化（2026-08-03）

**背景**：账本写入唯一入口（充值/退款/发券/结算）依赖幂等键 + `transaction.idempotency_key` 唯一约束。此前键由 provider 四个 service 手拼字符串（`"recharge:"+voucherNo`、`"issue:coupon:"+batchNo`、`"settle:"+orderID` 等），命名规则只存在于 conventions.md；业务号若含 `:`（如凭证号自带冒号）会破坏键空间隔离，无边界校验。

**决策**：`pkg/idempotency` 提供键构造契约：

- `Key(biz, bizNo) (string, error)`：构造 `{biz}:{bizNo}`；业务号为空或含 `:` 时返回错误（防键空间污染）
- 业务前缀常量（键空间隔离）：`recharge` / `refund` / `issue` / `settle`；发券按券类型分子命名空间（`issue:coupon` / `issue:voucher`——两券种批次号各自自增可能相同）
- `SettleRedeemKey(orderID, kind, refID)`：结算核销复合键 `settle:{orderID}:redeem:{kind}:{ref}`
- provider 五个构造点（account 充值/退款、coupon 发券、voucher 发券、order 结算消费/核销、reconciliation 对账匹配键）全部改引工具库；conventions.md 幂等键规则改为引用本包

**理由**：键空间规则代码化 + 边界校验防呆，业务号含 `:` 在入口即被拒绝；新业务接入沿用前缀常量与构造器，不再散落字符串字面量。

**影响**：键构造与「冲突回滚视为成功」语义解耦——本包只负责键的构造与校验，幂等写入语义留在各服务。
