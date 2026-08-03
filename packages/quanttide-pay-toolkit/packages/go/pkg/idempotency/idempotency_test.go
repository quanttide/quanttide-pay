package idempotency

import "testing"

func TestKey(t *testing.T) {
	tests := []struct {
		biz   string
		bizNo string
		want  string
	}{
		{Recharge, "voucher-001", "recharge:voucher-001"},
		{Refund, "voucher-002", "refund:voucher-002"},
		{IssueCoupon, "b1", "issue:coupon:b1"},
		{IssueVoucher, "b1", "issue:voucher:b1"}, // 与优惠券同批次号不冲突（子命名空间隔离）
		{Settle, "ORD001", "settle:ORD001"},
	}
	for _, tt := range tests {
		got, err := Key(tt.biz, tt.bizNo)
		if err != nil {
			t.Errorf("Key(%q, %q) 错误: %v", tt.biz, tt.bizNo, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Key(%q, %q) = %q, 期望 %q", tt.biz, tt.bizNo, got, tt.want)
		}
	}
}

func TestKey_Invalid(t *testing.T) {
	// 业务号为空（缺业务号）或含 ":"（键空间污染）必须拒绝。
	if _, err := Key(Recharge, ""); err == nil {
		t.Error("空业务号应返回错误")
	}
	if _, err := Key(Recharge, "v1:evil"); err == nil {
		t.Error("含冒号业务号应返回错误")
	}
	if _, err := Key("", "v1"); err == nil {
		t.Error("空业务前缀应返回错误")
	}
}

func TestSettleRedeemKey(t *testing.T) {
	got, err := SettleRedeemKey("ORD001", "coupon", 7)
	if err != nil {
		t.Fatal(err)
	}
	if want := "settle:ORD001:redeem:coupon:7"; got != want {
		t.Errorf("SettleRedeemKey = %q, 期望 %q", got, want)
	}
}

func TestSettleRedeemKey_Invalid(t *testing.T) {
	if _, err := SettleRedeemKey("", "coupon", 1); err == nil {
		t.Error("空订单号应返回错误")
	}
	if _, err := SettleRedeemKey("ORD:1", "coupon", 1); err == nil {
		t.Error("含冒号订单号应返回错误")
	}
	if _, err := SettleRedeemKey("ORD001", "", 1); err == nil {
		t.Error("空 kind 应返回错误")
	}
}
