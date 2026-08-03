// Package billing 提供结算抵扣计算契约（纯计算，无存储依赖）。
//
// v0.1.0 默认抵扣顺序（规则引擎后置）：
//  1. 满减券：满足门槛（≤ 剩余应付）中力度最大的一张
//  2. 折扣券：按 rate 优惠（9 折 = rate 90 = 省 10%），向下取整
//  3. 代金券：逐张抵扣 min(面值, 剩余应付)
//  4. 余额：补足剩余
//
// 输入输出金额均为分（int64）；余额不足返回 ErrInsufficientBalance。
// 顺序与力度规则为系统级契约（方案契约 v0.1），由服务端默认执行。
package billing

import "errors"

// 抵扣类型。
const (
	KindCoupon  = "coupon"  // 优惠券抵扣
	KindVoucher = "voucher" // 代金券抵现
	KindBalance = "balance" // 余额支付
)

// 错误。
var (
	// ErrInsufficientBalance 余额不足，无法完成结算。
	ErrInsufficientBalance = errors.New("billing: insufficient balance")
	// ErrInvalidAmount 订单金额必须为正整数（分）。
	ErrInvalidAmount = errors.New("billing: invalid amount")
)

// CouponInput 参与结算计算的优惠券（中立输入，不依赖具体优惠券模块）。
type CouponInput struct {
	ID        int64
	Type      string // discount / full_reduction
	Rate      int    // 折扣券：整数百分比（9 折 = 90）
	Threshold int64  // 满减券：门槛（分）
	Amount    int64  // 满减券：减额（分）
}

// VoucherInput 参与结算计算的代金券。
type VoucherInput struct {
	ID     int64
	Amount int64 // 面值（分）
}

// Deduction 一项抵扣。
type Deduction struct {
	Kind   string `json:"kind"`             // coupon / voucher / balance
	RefID  int64  `json:"ref_id,omitempty"` // 券 ID（balance 时为 0）
	Amount int64  `json:"amount"`           // 抵扣额（分）
}

// Calculate 结算计算（纯函数，无 I/O）：给定订单金额与可用券/余额，输出逐项抵扣明细。
func Calculate(amount int64, coupons []CouponInput, vouchers []VoucherInput, balance int64) ([]Deduction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	remain := amount
	var plan []Deduction

	// 1. 满减券（力度最大的一张）
	if c := bestFullReduction(coupons, remain); c != nil {
		plan = append(plan, Deduction{Kind: KindCoupon, RefID: c.ID, Amount: c.Amount})
		remain -= c.Amount
	}

	// 2. 折扣券（折扣力度最大的一张）：省 (100−rate)%
	if c := bestDiscount(coupons); c != nil {
		discount := remain * int64(100-c.Rate) / 100
		if discount > 0 {
			plan = append(plan, Deduction{Kind: KindCoupon, RefID: c.ID, Amount: discount})
			remain -= discount
		}
	}

	// 3. 代金券逐张抵现
	for _, v := range vouchers {
		if remain == 0 {
			break
		}
		d := min(v.Amount, remain)
		plan = append(plan, Deduction{Kind: KindVoucher, RefID: v.ID, Amount: d})
		remain -= d
	}

	// 4. 余额补足
	if remain > balance {
		return nil, ErrInsufficientBalance
	}
	if remain > 0 {
		plan = append(plan, Deduction{Kind: KindBalance, Amount: remain})
	}
	return plan, nil
}

// bestFullReduction 满足门槛（≤ remain）中减额最大的一张满减券；无则返回 nil。
func bestFullReduction(coupons []CouponInput, remain int64) *CouponInput {
	var best *CouponInput
	for i := range coupons {
		c := &coupons[i]
		if c.Type == "full_reduction" && c.Threshold <= remain && c.Amount > 0 {
			if best == nil || c.Amount > best.Amount {
				best = c
			}
		}
	}
	return best
}

// bestDiscount 折扣力度最大（rate 最低）的一张折扣券；无则返回 nil。
func bestDiscount(coupons []CouponInput) *CouponInput {
	var best *CouponInput
	for i := range coupons {
		c := &coupons[i]
		if c.Type == "discount" && c.Rate > 0 && c.Rate <= 100 {
			if best == nil || c.Rate < best.Rate {
				best = c
			}
		}
	}
	return best
}
