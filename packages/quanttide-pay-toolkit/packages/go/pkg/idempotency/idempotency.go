// Package idempotency 提供幂等键构造契约。
//
// 键空间规则：{业务}:{业务号}——前缀按业务类型隔离键空间，同一业务号在不同业务中互不冲突；
// 业务号为空或含 ":" 时返回错误（防键空间污染：业务号自带冒号会破坏隔离）。
// 本包只负责键的构造与校验，幂等写入语义（冲突回滚视为成功）由各服务实现。
package idempotency

import (
	"fmt"
	"strings"
)

// 业务前缀常量（键空间隔离）。
const (
	Recharge = "recharge" // 充值：打款凭证号
	Refund   = "refund"   // 退款：退款凭证号
	Issue    = "issue"    // 发券：批次号
	Settle   = "settle"   // 结算：商户订单号
)

// 发券子命名空间：优惠券与代金券批次号各自自增、可能相同，须按券类型隔离。
const (
	IssueCoupon  = Issue + ":coupon"  // 发券（优惠券）：批次号
	IssueVoucher = Issue + ":voucher" // 发券（代金券）：批次号
)

// Key 构造 {biz}:{bizNo}；业务号为空或含 ":" 时返回错误，防键空间污染。
func Key(biz, bizNo string) (string, error) {
	if biz == "" {
		return "", fmt.Errorf("idempotency: biz required")
	}
	if bizNo == "" || strings.Contains(bizNo, ":") {
		return "", fmt.Errorf("idempotency: invalid biz no %q", bizNo)
	}
	return biz + ":" + bizNo, nil
}

// SettleRedeemKey 构造结算核销复合键 settle:{orderID}:redeem:{kind}:{refID}。
// 一次结算按抵扣项（券/余额）分别写入核销交易，复合键保证每项独立幂等。
func SettleRedeemKey(orderID, kind string, refID int64) (string, error) {
	if orderID == "" || strings.Contains(orderID, ":") {
		return "", fmt.Errorf("idempotency: invalid order id %q", orderID)
	}
	if kind == "" {
		return "", fmt.Errorf("idempotency: kind required")
	}
	return fmt.Sprintf("%s:%s:redeem:%s:%d", Settle, orderID, kind, refID), nil
}
