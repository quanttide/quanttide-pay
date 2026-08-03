package ledger

import "testing"

func TestAffectsBalance(t *testing.T) {
	cases := []struct {
		t    Type
		want bool
	}{
		{TypeRecharge, true},
		{TypeRefund, true},
		{TypeConsume, true},
		{TypeIssue, false},
		{TypeRedeem, false},
		{Type("unknown"), false},
	}
	for _, c := range cases {
		if got := AffectsBalance(c.t); got != c.want {
			t.Errorf("AffectsBalance(%s) = %v, want %v", c.t, got, c.want)
		}
	}
}

func TestSignedAmount(t *testing.T) {
	cases := []struct {
		t    Type
		want int64
	}{
		{TypeRecharge, 100},
		{TypeRefund, -100},
		{TypeConsume, -100},
		{TypeIssue, 0},
		{TypeRedeem, 0},
	}
	for _, c := range cases {
		if got := SignedAmount(c.t, 100); got != c.want {
			t.Errorf("SignedAmount(%s) = %d, want %d", c.t, got, c.want)
		}
	}
}

func TestIsValid(t *testing.T) {
	for _, typ := range []Type{TypeRecharge, TypeRefund, TypeConsume, TypeIssue, TypeRedeem} {
		if !IsValid(typ) {
			t.Errorf("IsValid(%s) = false, want true", typ)
		}
	}
	for _, typ := range []Type{"", "payment", "REVERSAL"} {
		if IsValid(typ) {
			t.Errorf("IsValid(%q) = true, want false", typ)
		}
	}
}
