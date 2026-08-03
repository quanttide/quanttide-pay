package status

import "testing"

func TestCouponStatusValues(t *testing.T) {
	values := []struct {
		status CouponStatus
		want   string
	}{
		{CouponStatusIssued, "issued"},
		{CouponStatusUsed, "used"},
		{CouponStatusExpired, "expired"},
	}
	for _, v := range values {
		if string(v.status) != v.want {
			t.Errorf("coupon status = %q, want %q", v.status, v.want)
		}
		if !IsValidCouponStatus(v.status) {
			t.Errorf("IsValidCouponStatus(%q) = false", v.status)
		}
	}
	if IsValidCouponStatus("unknown") {
		t.Error("IsValidCouponStatus(unknown) = true, want false")
	}
}
