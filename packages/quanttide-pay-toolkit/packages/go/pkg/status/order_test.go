package status

import "testing"

func TestOrderStatusValues(t *testing.T) {
	values := []struct {
		status OrderStatus
		want   string
	}{
		{OrderStatusCreated, "created"},
		{OrderStatusSettled, "settled"},
	}
	for _, v := range values {
		if string(v.status) != v.want {
			t.Errorf("order status = %q, want %q", v.status, v.want)
		}
		if !IsValidOrderStatus(v.status) {
			t.Errorf("IsValidOrderStatus(%q) = false", v.status)
		}
	}
	if IsValidOrderStatus("unknown") {
		t.Error("IsValidOrderStatus(unknown) = true, want false")
	}
}
