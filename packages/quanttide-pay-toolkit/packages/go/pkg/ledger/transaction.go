package ledger

import "time"

// Transaction 一笔账本交易的不可变记录契约（JSON 形状，无存储绑定）。
// 存储实现（gorm 等）由消费应用负责；本结构是跨应用流转的统一形状。
// Type 强类型化为账本交易类型，避免裸字符串与类型语义脱节。
type Transaction struct {
	ID             int64     `json:"id"`
	AccountID      string    `json:"account_id"`
	Type           Type      `json:"type"`
	Amount         int64     `json:"amount"`
	BalanceAfter   int64     `json:"balance_after"`
	OrderID        string    `json:"order_id,omitempty"`
	IdempotencyKey string    `json:"-"`
	Note           string    `json:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// IsValid 报告记录是否为合法账本交易（类型已知、账户非空），供存库前防御校验。
func (t *Transaction) IsValid() bool {
	return t.AccountID != "" && IsValid(t.Type)
}

// Balance 从交易推导余额：Σ 带符号金额（余额 = 交易投影的唯一推导规则）。
// 发券/核销等不参与余额求和的交易自动贡献 0。
func Balance(txs []Transaction) int64 {
	var sum int64
	for _, t := range txs {
		sum += SignedAmount(t.Type, t.Amount)
	}
	return sum
}
