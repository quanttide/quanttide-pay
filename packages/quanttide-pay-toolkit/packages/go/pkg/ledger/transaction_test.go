package ledger

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTransactionJSONShape(t *testing.T) {
	createdAt := time.Date(2026, 8, 3, 10, 0, 0, 0, time.UTC)
	tx := Transaction{
		ID: 1, AccountID: "acc_1", Type: TypeRecharge, Amount: 10000,
		BalanceAfter: 10000, IdempotencyKey: "recharge:v1", Note: "对公打款", CreatedAt: createdAt,
	}
	b, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	// IdempotencyKey 不进入 JSON（json:"-"）；金额为整数分。
	want := `{"id":1,"account_id":"acc_1","type":"recharge","amount":10000,"balance_after":10000,"note":"对公打款","created_at":"2026-08-03T10:00:00Z"}`
	if string(b) != want {
		t.Errorf("json = %s, 期望 %s", b, want)
	}
	var got Transaction
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got.Type != TypeRecharge || got.Amount != 10000 {
		t.Errorf("round-trip = %+v", got)
	}
}

func TestTransactionIsValid(t *testing.T) {
	valid := Transaction{AccountID: "acc_1", Type: TypeConsume}
	if !valid.IsValid() {
		t.Error("合法记录应通过校验")
	}
	for _, tx := range []Transaction{
		{AccountID: "", Type: TypeConsume}, // 账户为空
		{AccountID: "acc_1", Type: "unknown"},
	} {
		if tx.IsValid() {
			t.Errorf("%+v 应不合法", tx)
		}
	}
}

func TestBalance(t *testing.T) {
	// 余额 = Σ带符号金额：充值 +，退款/消费 −，发券/核销不参与。
	// 10000 − 3000 − 2000 = 5000（发券/核销贡献 0）。
	txs := []Transaction{
		{Type: TypeRecharge, Amount: 10000},
		{Type: TypeConsume, Amount: 3000},
		{Type: TypeRefund, Amount: 2000},
		{Type: TypeIssue, Amount: 5000},
		{Type: TypeRedeem, Amount: 1000},
	}
	if got := Balance(txs); got != 5000 {
		t.Errorf("Balance = %d, want 5000", got)
	}
	if got := Balance(nil); got != 0 {
		t.Errorf("Balance(nil) = %d, want 0", got)
	}
}
