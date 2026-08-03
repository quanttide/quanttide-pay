package status

import "testing"

func TestPaymentStatusValues(t *testing.T) {
	values := []struct {
		status PaymentStatus
		want   string
	}{
		{PaymentStatusCreated, "created"},
		{PaymentStatusPaying, "paying"},
		{PaymentStatusSucceeded, "succeeded"},
		{PaymentStatusFailed, "failed"},
		{PaymentStatusClosed, "closed"},
		{PaymentStatusRefunding, "refunding"},
		{PaymentStatusRefunded, "refunded"},
	}
	for _, v := range values {
		if string(v.status) != v.want {
			t.Errorf("PaymentStatus(%s) = %q, 期望 %q", v.status, v.status, v.want)
		}
	}
}

func TestParseWechatTradeState(t *testing.T) {
	tests := []struct {
		code string
		want PaymentStatus
	}{
		{"SUCCESS", PaymentStatusSucceeded},
		{"REFUND", PaymentStatusRefunding},
		{"NOTPAY", PaymentStatusCreated},
		{"USERPAYING", PaymentStatusPaying},
		{"CLOSED", PaymentStatusClosed},
		{"PAYERROR", PaymentStatusFailed},
	}
	for _, tt := range tests {
		got, err := ParseWechatTradeState(tt.code)
		if err != nil {
			t.Errorf("ParseWechatTradeState(%q) 错误: %v", tt.code, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseWechatTradeState(%q) = %q, 期望 %q", tt.code, got, tt.want)
		}
	}
}

func TestParseWechatTradeState_Unknown(t *testing.T) {
	if _, err := ParseWechatTradeState("SOME_NEW_STATE"); err == nil {
		t.Error("未知微信 TradeState 应返回错误，而不是静默降级")
	}
}

func TestParseAlipayTradeStatus(t *testing.T) {
	tests := []struct {
		code string
		want PaymentStatus
	}{
		{"TRADE_SUCCESS", PaymentStatusSucceeded},
		{"TRADE_FINISHED", PaymentStatusSucceeded},
		{"WAIT_BUYER_PAY", PaymentStatusCreated},
		{"TRADE_CLOSED", PaymentStatusClosed},
	}
	for _, tt := range tests {
		got, err := ParseAlipayTradeStatus(tt.code)
		if err != nil {
			t.Errorf("ParseAlipayTradeStatus(%q) 错误: %v", tt.code, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseAlipayTradeStatus(%q) = %q, 期望 %q", tt.code, got, tt.want)
		}
	}
}

func TestParseAlipayTradeStatus_Unknown(t *testing.T) {
	if _, err := ParseAlipayTradeStatus("SOME_UNKNOWN_STATUS"); err == nil {
		t.Error("未知支付宝 trade_status 应返回错误，而不是静默降级")
	}
}

func TestIsValidPaymentStatus(t *testing.T) {
	if !IsValidPaymentStatus(PaymentStatusSucceeded) {
		t.Error("PaymentStatusSucceeded 应合法")
	}
	for _, s := range []PaymentStatus{"", "SUCCESS", "unknown"} {
		if IsValidPaymentStatus(s) {
			t.Errorf("PaymentStatus(%q) 应不合法", s)
		}
	}
}
