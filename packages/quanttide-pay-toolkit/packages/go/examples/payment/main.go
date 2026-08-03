// Command payment 演示 quanttide-pay-toolkit Go 基础库的用法。
package main

import (
	"fmt"
	"log"

	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/money"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/order"
	"github.com/quanttide/quanttide-pay-toolkit/packages/go/pkg/status"
)

func main() {
	// 金额：以分为单位，避免浮点误差
	total, err := money.MustMoney(1234, "CNY").Add(money.MustMoney(100, "CNY"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("金额合计: %d 分（%s）\n", total.Cents, total.Currency)

	// 状态：取值可直接存入数据库、与渠道报文流转
	fmt.Printf("支付成功状态: %s\n", status.PaymentStatusSucceeded)

	// 订单号：PAY20260803143000 + 10 位随机数字
	orderNo, err := order.GenerateNow("PAY")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("订单号: %s\n", orderNo)
}
