package billing

import (
	"errors"
	"testing"
)

func fullReduction(id int64, threshold, amount int64) CouponInput {
	return CouponInput{ID: id, Type: "full_reduction", Threshold: threshold, Amount: amount}
}

func discount(id int64, rate int) CouponInput {
	return CouponInput{ID: id, Type: "discount", Rate: rate}
}

func voucher(id int64, amount int64) VoucherInput {
	return VoucherInput{ID: id, Amount: amount}
}

func sum(plan []Deduction) int64 {
	var total int64
	for _, d := range plan {
		total += d.Amount
	}
	return total
}

func TestCalculate(t *testing.T) {
	cases := []struct {
		name      string
		amount    int64
		coupons   []CouponInput
		vouchers  []VoucherInput
		balance   int64
		wantSum   int64 // 抵扣总额（等于订单金额时全额覆盖）
		wantKinds map[string]int
	}{
		{
			name: "无券全额余额", amount: 10000, balance: 100000,
			wantSum: 10000, wantKinds: map[string]int{"balance": 1},
		},
		{
			name: "仅满减", amount: 10000,
			coupons: []CouponInput{fullReduction(1, 8000, 2000)}, balance: 100000,
			wantSum: 10000, wantKinds: map[string]int{"coupon": 1, "balance": 1},
		},
		{
			name: "满减门槛未满足", amount: 5000,
			coupons: []CouponInput{fullReduction(1, 8000, 2000)}, balance: 100000,
			wantSum: 5000, wantKinds: map[string]int{"balance": 1},
		},
		{
			name: "满减取力度最大", amount: 10000, balance: 100000,
			coupons: []CouponInput{
				fullReduction(1, 8000, 1000),
				fullReduction(2, 9000, 3000),
				fullReduction(3, 500, 50),
			},
			wantSum: 10000, wantKinds: map[string]int{"coupon": 1, "balance": 1},
		},
		{
			name: "仅折扣", amount: 10000,
			coupons: []CouponInput{discount(1, 90)}, balance: 100000,
			wantSum: 10000, wantKinds: map[string]int{"coupon": 1, "balance": 1},
		},
		{
			name: "折扣向下取整", amount: 9999,
			coupons: []CouponInput{discount(1, 90)}, balance: 100000,
			wantSum: 9999, wantKinds: map[string]int{"coupon": 1, "balance": 1},
		},
		{
			name: "满减加折扣", amount: 10000, balance: 100000,
			coupons: []CouponInput{fullReduction(1, 8000, 2000), discount(2, 90)},
			// 满减 2000 → 剩余 8000 → 折扣省 800（9 折）→ 余额 7200
			wantSum: 10000, wantKinds: map[string]int{"coupon": 2, "balance": 1},
		},
		{
			name: "代金券全额抵扣", amount: 10000,
			vouchers: []VoucherInput{voucher(1, 6000), voucher(2, 4000)}, balance: 100000,
			wantSum: 10000, wantKinds: map[string]int{"voucher": 2},
		},
		{
			name: "代金券部分抵扣", amount: 10000,
			vouchers: []VoucherInput{voucher(1, 6000)}, balance: 100000,
			wantSum: 10000, wantKinds: map[string]int{"voucher": 1, "balance": 1},
		},
		{
			name: "混合抵扣", amount: 10000, balance: 100000,
			coupons:  []CouponInput{fullReduction(1, 8000, 2000)},
			vouchers: []VoucherInput{voucher(2, 3000)},
			// 满减 2000 → 代金券 3000 → 余额 5000
			wantSum: 10000, wantKinds: map[string]int{"coupon": 1, "voucher": 1, "balance": 1},
		},
		{
			name: "券额度超出订单", amount: 5000, balance: 100000,
			coupons:  []CouponInput{fullReduction(1, 1000, 4000)},
			vouchers: []VoucherInput{voucher(2, 3000)},
			// 满减 4000 → 剩余 1000 → 代金券 1000（部分使用）→ 余额 0
			wantSum: 5000, wantKinds: map[string]int{"coupon": 1, "voucher": 1},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			plan, err := Calculate(c.amount, c.coupons, c.vouchers, c.balance)
			if err != nil {
				t.Fatalf("Calculate: %v", err)
			}
			if got := sum(plan); got != c.wantSum {
				t.Errorf("sum = %d, want %d (plan=%+v)", got, c.wantSum, plan)
			}
			kinds := map[string]int{}
			for _, d := range plan {
				kinds[d.Kind]++
			}
			for k, n := range c.wantKinds {
				if kinds[k] != n {
					t.Errorf("kind %s count = %d, want %d", k, kinds[k], n)
				}
			}
			// 抵扣总额不得超过订单金额
			if sum(plan) > c.amount {
				t.Errorf("sum %d exceeds amount %d", sum(plan), c.amount)
			}
		})
	}
}

func TestCalculate_Errors(t *testing.T) {
	// 金额非正
	if _, err := Calculate(0, nil, nil, 100); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("amount=0 err = %v, want ErrInvalidAmount", err)
	}
	if _, err := Calculate(-1, nil, nil, 100); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("amount<0 err = %v, want ErrInvalidAmount", err)
	}

	// 余额不足
	_, err := Calculate(10000, nil, nil, 9999)
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Errorf("err = %v, want ErrInsufficientBalance", err)
	}
	// 券后仍不足
	_, err = Calculate(10000, []CouponInput{fullReduction(1, 1000, 2000)}, nil, 100)
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Errorf("err = %v, want ErrInsufficientBalance", err)
	}
}

// TestCalculate_PickBestCoupon 多张同类券时选力度最大：满减选减额最大、折扣选 rate 最低。
func TestCalculate_PickBestCoupon(t *testing.T) {
	// 满减：满 8000 减 1000 / 满 9000 减 3000 → 选减额最大的 3000
	plan, err := Calculate(10000, []CouponInput{
		fullReduction(1, 8000, 1000),
		fullReduction(2, 9000, 3000),
	}, nil, 100000)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != 2 || plan[0].RefID != 2 || plan[0].Amount != 3000 {
		t.Errorf("plan = %+v, want 满减券 ID2 减 3000", plan)
	}

	// 折扣：rate 90（9 折，省 10%）与 rate 80（8 折，省 20%）→ 选 8 折（力度最大，省 2000）
	plan, err = Calculate(10000, []CouponInput{
		discount(1, 90),
		discount(2, 80),
	}, nil, 100000)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != 2 || plan[0].RefID != 2 || plan[0].Amount != 2000 {
		t.Errorf("plan = %+v, want 折扣券 ID2（rate 80）省 2000", plan)
	}
}

func TestCalculate_InvalidCoupons(t *testing.T) {
	// 非法折扣券（rate 越界）被忽略
	plan, err := Calculate(10000, []CouponInput{
		{ID: 1, Type: "discount", Rate: 0},
		{ID: 2, Type: "discount", Rate: 101},
		{ID: 3, Type: "unknown"},
	}, nil, 10000)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}
	if len(plan) != 1 || plan[0].Kind != KindBalance {
		t.Errorf("plan = %+v, want only balance", plan)
	}

	// 折扣为 100%（9 折语义下省 0%）→ 不生成抵扣项，全额余额支付
	plan, err = Calculate(100, []CouponInput{{ID: 1, Type: "discount", Rate: 100}}, nil, 100)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}
	if len(plan) != 1 || plan[0].Kind != KindBalance || plan[0].Amount != 100 {
		t.Errorf("plan = %+v, want balance 100", plan)
	}
}
