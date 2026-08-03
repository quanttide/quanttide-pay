// Package contracttest 契约测试 Go runner。
//
// 读取工具库根 tests/fixtures/ 下的共享测试向量，断言 Go 实现与契约一致：
// 同一输入 → 同一输出 / 同一拒绝行为。其他语言实现（如 packages/dart、packages/rust）消费
// 同一 fixtures，即完成多端对齐——fixtures 是契约的唯一权威。
package contracttest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/ledger"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/money"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/status"
)

// fixture 读取共享测试向量。fixtures 位于工具库根 tests/fixtures
// （本包位于 packages/go/contracttest，相对路径上溯三级）。
func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "tests", "fixtures", name))
	if err != nil {
		t.Fatalf("读取 fixture %s: %v", name, err)
	}
	return b
}

func TestMoneyContract(t *testing.T) {
	var fx struct {
		Valid []struct {
			JSON        string `json:"json"`
			ExpectCents int64  `json:"expect_cents"`
		} `json:"valid"`
		Invalid []struct {
			JSON   string `json:"json"`
			Reason string `json:"reason"`
		} `json:"invalid"`
	}
	if err := json.Unmarshal(fixture(t, "money.json"), &fx); err != nil {
		t.Fatal(err)
	}
	for _, c := range fx.Valid {
		var m money.Money
		if err := json.Unmarshal([]byte(c.JSON), &m); err != nil {
			t.Errorf("契约应接受 %s: %v", c.JSON, err)
			continue
		}
		if got := m.Amount(); got != c.ExpectCents {
			t.Errorf("%s = %d 分, 契约期望 %d", c.JSON, got, c.ExpectCents)
		}
	}
	for _, c := range fx.Invalid {
		var m money.Money
		if err := json.Unmarshal([]byte(c.JSON), &m); err == nil {
			t.Errorf("契约应拒绝 %s（%s）", c.JSON, c.Reason)
		}
	}
}

func TestStatusContract(t *testing.T) {
	var fx struct {
		WechatTradeState   map[string]string `json:"wechat_trade_state"`
		AlipayTradeStatus  map[string]string `json:"alipay_trade_status"`
		WechatRefundStatus map[string]string `json:"wechat_refund_status"`
		UnknownCodes       []string          `json:"unknown_codes"`
	}
	if err := json.Unmarshal(fixture(t, "status.json"), &fx); err != nil {
		t.Fatal(err)
	}
	for code, want := range fx.WechatTradeState {
		got, err := status.ParseWechatTradeState(code)
		if err != nil {
			t.Errorf("ParseWechatTradeState(%q) 错误: %v", code, err)
			continue
		}
		if string(got) != want {
			t.Errorf("ParseWechatTradeState(%q) = %q, 契约期望 %q", code, got, want)
		}
	}
	for code, want := range fx.AlipayTradeStatus {
		got, err := status.ParseAlipayTradeStatus(code)
		if err != nil {
			t.Errorf("ParseAlipayTradeStatus(%q) 错误: %v", code, err)
			continue
		}
		if string(got) != want {
			t.Errorf("ParseAlipayTradeStatus(%q) = %q, 契约期望 %q", code, got, want)
		}
	}
	for code, want := range fx.WechatRefundStatus {
		got, err := status.ParseWechatRefundStatus(code)
		if err != nil {
			t.Errorf("ParseWechatRefundStatus(%q) 错误: %v", code, err)
			continue
		}
		if string(got) != want {
			t.Errorf("ParseWechatRefundStatus(%q) = %q, 契约期望 %q", code, got, want)
		}
	}
	for _, code := range fx.UnknownCodes {
		if _, err := status.ParseWechatTradeState(code); err == nil {
			t.Errorf("未知渠道码 %q 微信解析应报错（契约：不静默降级）", code)
		}
		if _, err := status.ParseAlipayTradeStatus(code); err == nil {
			t.Errorf("未知渠道码 %q 支付宝解析应报错（契约：不静默降级）", code)
		}
		if _, err := status.ParseWechatRefundStatus(code); err == nil {
			t.Errorf("未知渠道码 %q 微信退款解析应报错（契约：不静默降级）", code)
		}
	}
}

func TestLedgerContract(t *testing.T) {
	var fx struct {
		TypeSemantics map[string]struct {
			AffectsBalance bool  `json:"affects_balance"`
			SignedOf100    int64 `json:"signed_of_100"`
		} `json:"type_semantics"`
		BalanceCases []struct {
			Name         string `json:"name"`
			Transactions []struct {
				Type   string `json:"type"`
				Amount int64  `json:"amount"`
			} `json:"transactions"`
			Balance int64 `json:"balance"`
		} `json:"balance_cases"`
	}
	if err := json.Unmarshal(fixture(t, "ledger.json"), &fx); err != nil {
		t.Fatal(err)
	}
	for typ, want := range fx.TypeSemantics {
		if got := ledger.AffectsBalance(ledger.Type(typ)); got != want.AffectsBalance {
			t.Errorf("AffectsBalance(%q) = %v, 契约期望 %v", typ, got, want.AffectsBalance)
		}
		if got := ledger.SignedAmount(ledger.Type(typ), 100); got != want.SignedOf100 {
			t.Errorf("SignedAmount(%q, 100) = %d, 契约期望 %d", typ, got, want.SignedOf100)
		}
	}
	for _, c := range fx.BalanceCases {
		txs := make([]ledger.Transaction, 0, len(c.Transactions))
		for _, tx := range c.Transactions {
			txs = append(txs, ledger.Transaction{Type: ledger.Type(tx.Type), Amount: tx.Amount})
		}
		if got := ledger.Balance(txs); got != c.Balance {
			t.Errorf("Balance(%s) = %d, 契约期望 %d", c.Name, got, c.Balance)
		}
	}
}
