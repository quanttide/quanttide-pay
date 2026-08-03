package money

import (
	"encoding/json"
	"testing"
)

func TestNewAndAccessors(t *testing.T) {
	m := New(1234, "CNY")
	if m.Amount() != 1234 {
		t.Errorf("Amount() = %d, want 1234", m.Amount())
	}
	if m.Currency().Code != "CNY" {
		t.Errorf("Currency() = %s, want CNY", m.Currency().Code)
	}
	if !m.IsPositive() || m.IsNegative() || m.IsZero() {
		t.Error("1234 分应为正金额")
	}
	if !New(0, "CNY").IsZero() {
		t.Error("0 分应视为零金额")
	}
	if m.Currency().Fraction != 2 {
		t.Errorf("CNY Fraction = %d, want 2", m.Currency().Fraction)
	}
}

func TestArithmetic(t *testing.T) {
	sum, err := New(110, "CNY").Add(New(220, "CNY"))
	if err != nil || sum.Amount() != 330 {
		t.Fatalf("Add = %d, %v; want 330", sum.Amount(), err)
	}
	diff, err := New(330, "CNY").Subtract(New(110, "CNY"))
	if err != nil || diff.Amount() != 220 {
		t.Fatalf("Subtract = %d, %v; want 220", diff.Amount(), err)
	}
	if _, err := New(100, "CNY").Add(New(100, "USD")); err == nil {
		t.Fatal("不同币种相加未返回错误")
	}
	if v, err := New(9999, "CNY").Compare(New(10000, "CNY")); err != nil || v != -1 {
		t.Errorf("Compare = %d, %v; want -1", v, err)
	}
	parts, err := New(100, "CNY").Split(3)
	if err != nil {
		t.Fatal(err)
	}
	if parts[0].Amount() != 34 || parts[1].Amount() != 33 || parts[2].Amount() != 33 {
		t.Errorf("Split(3) = %d,%d,%d; want 34,33,33", parts[0].Amount(), parts[1].Amount(), parts[2].Amount())
	}
}

func TestCentsOf(t *testing.T) {
	if CentsOf(nil) != 0 {
		t.Error("CentsOf(nil) 应为 0")
	}
	if CentsOf(New(500, "CNY")) != 500 {
		t.Error("CentsOf(500分) 应为 500")
	}
}

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		in   *Money
		want string
	}{
		{"零金额", New(0, "CNY"), `{"amount":0,"currency":"CNY"}`},
		{"整数分", New(9999, "CNY"), `{"amount":9999,"currency":"CNY"}`},
		{"负数", New(-150, "CNY"), `{"amount":-150,"currency":"CNY"}`},
		{"零位小数币种", New(100, "JPY"), `{"amount":100,"currency":"JPY"}`},
		{"零值 Money", &Money{}, `{"amount":0,"currency":""}`},
	}
	for _, c := range cases {
		b, err := json.Marshal(c.in)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if string(b) != c.want {
			t.Errorf("%s: Marshal = %s, want %s", c.name, b, c.want)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int64
		ok   bool
	}{
		{"整数分", `{"amount":9999,"currency":"CNY"}`, 9999, true},
		{"零金额带币种", `{"amount":0,"currency":"CNY"}`, 0, true},
		{"零金额无币种", `{"amount":0,"currency":""}`, 0, true},
		{"负数", `{"amount":-150,"currency":"CNY"}`, -150, true},
		{"小写币种", `{"amount":100,"currency":"cny"}`, 100, true},
		{"零位小数币种", `{"amount":100,"currency":"JPY"}`, 100, true},
		{"小数拒绝", `{"amount":99.99,"currency":"CNY"}`, 0, false},
		{"字符串拒绝", `{"amount":"9999","currency":"CNY"}`, 0, false},
		{"指数记法拒绝", `{"amount":1e3,"currency":"CNY"}`, 0, false},
		{"未知币种拒绝", `{"amount":100,"currency":"XYZ"}`, 0, false},
		{"非零金额缺币种拒绝", `{"amount":100}`, 0, false},
		{"非 JSON 拒绝", `not-json`, 0, false},
	}
	for _, c := range cases {
		var m Money
		err := json.Unmarshal([]byte(c.in), &m)
		if c.ok {
			if err != nil {
				t.Errorf("%s: Unmarshal 错误: %v", c.name, err)
				continue
			}
			if m.Amount() != c.want {
				t.Errorf("%s: Amount = %d, want %d", c.name, m.Amount(), c.want)
			}
		} else if err == nil {
			t.Errorf("%s: 应拒绝，但解析成功: %+v", c.name, m)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, cents := range []int64{0, 1, 9999, 1000000, 123456789, -150, -123456789} {
		m := New(cents, "CNY")
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		var got Money
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("Unmarshal(%s): %v", b, err)
		}
		if got.Amount() != cents || got.Currency().Code != "CNY" {
			t.Errorf("round trip %d → %s → %d (%s)", cents, b, got.Amount(), got.Currency().Code)
		}
	}
}
