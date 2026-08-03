package money

import (
	"strings"
	"testing"
)

func TestNewMoney(t *testing.T) {
	m, err := NewMoney(1234, "CNY")
	if err != nil {
		t.Fatalf("NewMoney() 返回错误: %v", err)
	}
	if m.Cents != 1234 || m.Currency != "CNY" {
		t.Fatalf("NewMoney() = %+v, 期望 {1234 CNY}", m)
	}
}

func TestNewMoneyNegativeCents(t *testing.T) {
	if _, err := NewMoney(-1, "CNY"); err == nil || !strings.Contains(err.Error(), "不能为负数") {
		t.Fatalf("NewMoney(-1) 错误 = %v, 期望包含「不能为负数」", err)
	}
}

func TestNewMoneyInvalidCurrency(t *testing.T) {
	if _, err := NewMoney(100, "cny"); err == nil || !strings.Contains(err.Error(), "币种代码不合法") {
		t.Fatalf("NewMoney(100, \"cny\") 错误 = %v, 期望包含「币种代码不合法」", err)
	}
}

func TestMustMoney(t *testing.T) {
	m := MustMoney(100, "CNY")
	if m.Cents != 100 {
		t.Fatalf("MustMoney(100) = %+v", m)
	}
}

func TestMoneyAdd(t *testing.T) {
	sum, err := MustMoney(110, "CNY").Add(MustMoney(220, "CNY"))
	if err != nil {
		t.Fatalf("Add() 返回错误: %v", err)
	}
	if sum.Cents != 330 {
		t.Fatalf("Add() = %d, 期望 330", sum.Cents)
	}
}

func TestMoneySub(t *testing.T) {
	diff, err := MustMoney(330, "CNY").Sub(MustMoney(110, "CNY"))
	if err != nil {
		t.Fatalf("Sub() 返回错误: %v", err)
	}
	if diff.Cents != 220 {
		t.Fatalf("Sub() = %d, 期望 220", diff.Cents)
	}
}

func TestMoneyAddDifferentCurrency(t *testing.T) {
	if _, err := MustMoney(100, "CNY").Add(MustMoney(100, "USD")); err == nil {
		t.Fatal("不同币种相加未返回错误")
	}
}

func TestMoneySubNegative(t *testing.T) {
	if _, err := MustMoney(10, "CNY").Sub(MustMoney(20, "CNY")); err == nil {
		t.Fatal("结果为负数的减法未返回错误")
	}
}

func TestMoneyIsZero(t *testing.T) {
	if !MustMoney(0, "CNY").IsZero() {
		t.Fatal("Money{0} 应视为零")
	}
	if MustMoney(1, "CNY").IsZero() {
		t.Fatal("Money{1} 不应视为零")
	}
}
