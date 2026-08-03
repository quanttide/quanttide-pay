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
