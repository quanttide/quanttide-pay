// Package money 提供金额值对象：以最小货币单位（分）的整数存储，避免浮点误差。
package money

import (
	"fmt"

	"github.com/quanttide/quanttide-pay-toolkit/packages/go/internal/currency"
)

// Money 是金额值对象。
type Money struct {
	// Cents 是以分为单位的金额，如 1.23 元记作 123。
	Cents int64
	// Currency 是 ISO 4217 币种代码，如 CNY。
	Currency string
}

// NewMoney 以分为单位构造金额。
func NewMoney(cents int64, cur string) (Money, error) {
	if cents < 0 {
		return Money{}, fmt.Errorf("money: 金额不能为负数: %d", cents)
	}
	if !currency.Valid(cur) {
		return Money{}, fmt.Errorf("money: 币种代码不合法: %q", cur)
	}
	return Money{Cents: cents, Currency: cur}, nil
}

// MustMoney 以分为单位构造金额，输入非法时 panic，仅用于常量等确定场景。
func MustMoney(cents int64, cur string) Money {
	m, err := NewMoney(cents, cur)
	if err != nil {
		panic(err)
	}
	return m
}

// Add 返回两金额之和，币种不一致时返回错误。
func (m Money) Add(other Money) (Money, error) {
	if err := m.checkCurrency(other); err != nil {
		return Money{}, err
	}
	return Money{Cents: m.Cents + other.Cents, Currency: m.Currency}, nil
}

// Sub 返回两金额之差，币种不一致或结果为负数时返回错误。
func (m Money) Sub(other Money) (Money, error) {
	if err := m.checkCurrency(other); err != nil {
		return Money{}, err
	}
	return NewMoney(m.Cents-other.Cents, m.Currency)
}

// IsZero 报告金额是否为零。
func (m Money) IsZero() bool {
	return m.Cents == 0
}

func (m Money) checkCurrency(other Money) error {
	if m.Currency != other.Currency {
		return fmt.Errorf("money: 币种不一致: %s 与 %s", m.Currency, other.Currency)
	}
	return nil
}
