// Package ledger 提供账本交易类型契约：类型枚举、余额影响语义与校验。
// 取值使用英文小写字符串，可直接存入数据库并与账务/对账应用流转。
package ledger

// Type 账本交易类型。
type Type string

// 交易类型。
const (
	TypeRecharge Type = "recharge" // 充值（对公打款入账）
	TypeRefund   Type = "refund"   // 退款（多退登记：对公退款出账）
	TypeConsume  Type = "consume"  // 消费（余额支付部分）
	TypeIssue    Type = "issue"    // 发券（信息性记录，不影响余额）
	TypeRedeem   Type = "redeem"   // 核销（券抵扣部分，不影响余额）
)

// AffectsBalance 该类型是否影响余额（发券/核销不参与余额求和）。
func AffectsBalance(t Type) bool {
	return t == TypeRecharge || t == TypeRefund || t == TypeConsume
}

// SignedAmount 余额方向的带符号金额：充值 +，退款/消费 −，其余 0。
func SignedAmount(t Type, amount int64) int64 {
	switch t {
	case TypeRecharge:
		return amount
	case TypeRefund, TypeConsume:
		return -amount
	default:
		return 0
	}
}

// IsValid 报告交易类型是否为已知类型（存库前防御校验）。
func IsValid(t Type) bool {
	switch t {
	case TypeRecharge, TypeRefund, TypeConsume, TypeIssue, TypeRedeem:
		return true
	}
	return false
}
