// Package money 提供完整金额值对象（基于 go-money，Go 社区事实标准）：
// 以货币最小单位（分）的整数存储，杜绝浮点误差；币种遵循 ISO 4217。
//
// 本包是 go-money 的薄契约层：类型与 API 直接复用 go-money（别名/转发），
// 仅覆盖其 JSON 反序列化注入点做严格整数校验——go-money 默认实现把 amount
// 按 float64 解析后 int64() 截断，非整数会静默舍入入账，与本库契约冲突。
// 序列化保持 go-money 默认（本就是整数分输出）。
//
// JSON 边界契约：{"amount": <整数分>, "currency": "<ISO 4217 币种>"}——
// amount 必须为 JSON 整数（拒绝小数、字符串、指数记法）；未知币种拒绝；
// 非零金额必须携带币种。全链路（内部、传输、对账）单位统一为分。
//
// 注意：勿使用 go-money 的 NewFromFloat（浮点构造入口，本包刻意不转发）。
package money

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	gomoney "github.com/Rhymond/go-money"
)

// Money 金额值对象（go-money）：整数最小单位 + ISO 4217 币种。
type Money = gomoney.Money

// Currency ISO 4217 币种信息。
type Currency = gomoney.Currency

// 常用 ISO 4217 币种代码（完整列表见 go-money constants.go）。
const (
	CNY = gomoney.CNY
	USD = gomoney.USD
	EUR = gomoney.EUR
	JPY = gomoney.JPY
	HKD = gomoney.HKD
	GBP = gomoney.GBP
)

// New 以最小单位（分）构造金额。
func New(amount int64, code string) *Money {
	return gomoney.New(amount, code)
}

// GetCurrency 返回 ISO 4217 币种信息，未知币种返回 nil。
func GetCurrency(code string) *Currency {
	return gomoney.GetCurrency(code)
}

// CentsOf 返回金额的分值，nil 视为零金额。
func CentsOf(m *Money) int64 {
	if m == nil {
		return 0
	}
	return m.Amount()
}

func init() {
	gomoney.UnmarshalJSON = func(m *gomoney.Money, b []byte) error {
		var raw struct {
			Amount   json.RawMessage `json:"amount"`
			Currency string          `json:"currency"`
		}
		if err := json.Unmarshal(b, &raw); err != nil {
			return fmt.Errorf("money: invalid json: %w", err)
		}
		code := strings.ToUpper(strings.TrimSpace(raw.Currency))
		if code != "" && gomoney.GetCurrency(code) == nil {
			return fmt.Errorf("money: unsupported currency %q", raw.Currency)
		}
		cents := int64(0)
		if len(raw.Amount) > 0 && string(raw.Amount) != "null" {
			var err error
			cents, err = parseCents(raw.Amount)
			if err != nil {
				return err
			}
		}
		if cents != 0 && code == "" {
			return fmt.Errorf("money: currency required for non-zero amount")
		}
		*m = *gomoney.New(cents, code)
		return nil
	}
}

// parseCents 严格整数解析：仅接受 JSON 整数（可选负号），拒绝小数、字符串、指数与加号。
func parseCents(b []byte) (int64, error) {
	s := string(b)
	if s == "" || s[0] == '"' || s[0] == '\'' || s[0] == '+' {
		return 0, fmt.Errorf("money: amount must be an integer number: %s", b)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("money: amount must be an integer number: %s", b)
	}
	return v, nil
}
