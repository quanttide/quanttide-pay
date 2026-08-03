// Command payment 演示 quanttide-pay-toolkit Go 基础库的用法。
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/money"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/order"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/status"
)

func main() {
	// 金额：完整值对象（整数分 + ISO 4217 币种），避免浮点误差
	total, err := money.New(1234, money.CNY).Add(money.New(100, money.CNY))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("金额合计: %s\n", total.Display())

	// JSON 边界：整数分 + 币种，严格校验
	b, err := json.Marshal(money.New(9999, money.CNY))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("序列化: %s\n", b)

	var amount money.Money
	if err := json.Unmarshal([]byte(`{"amount": 9999, "currency": "CNY"}`), &amount); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("反序列化: %d 分\n", amount.Amount())

	// 状态：取值可直接存入数据库、与渠道报文流转
	fmt.Printf("支付成功状态: %s\n", status.PaymentStatusSucceeded)

	// 订单号：PAY20260803143000 + 10 位随机数字
	orderNo, err := order.GenerateNow("PAY")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("订单号: %s\n", orderNo)
}
