package status

import "testing"

func TestRefundStatusValues(t *testing.T) {
	values := []struct {
		status RefundStatus
		want   string
	}{
		{RefundStatusCreated, "created"},
		{RefundStatusProcessing, "processing"},
		{RefundStatusSucceeded, "succeeded"},
		{RefundStatusFailed, "failed"},
	}
	for _, v := range values {
		if string(v.status) != v.want {
			t.Errorf("RefundStatus(%s) = %q, 期望 %q", v.status, v.status, v.want)
		}
	}
}

func TestParseWechatRefundStatus(t *testing.T) {
	tests := []struct {
		code string
		want RefundStatus
	}{
		{"SUCCESS", RefundStatusSucceeded},
		{"PROCESSING", RefundStatusProcessing},
		{"CLOSED", RefundStatusFailed},
		{"ABNORMAL", RefundStatusFailed},
	}
	for _, tt := range tests {
		got, err := ParseWechatRefundStatus(tt.code)
		if err != nil {
			t.Errorf("ParseWechatRefundStatus(%q) 错误: %v", tt.code, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseWechatRefundStatus(%q) = %q, 期望 %q", tt.code, got, tt.want)
		}
	}
}

func TestParseWechatRefundStatus_Unknown(t *testing.T) {
	if _, err := ParseWechatRefundStatus("SOME_NEW_STATUS"); err == nil {
		t.Error("未知微信退款 status 应返回错误，而不是兜底")
	}
}

func TestIsValidRefundStatus(t *testing.T) {
	if !IsValidRefundStatus(RefundStatusSucceeded) {
		t.Error("RefundStatusSucceeded 应合法")
	}
	for _, s := range []RefundStatus{"", "SUCCESS", "closed"} {
		if IsValidRefundStatus(s) {
			t.Errorf("RefundStatus(%q) 应不合法", s)
		}
	}
}
