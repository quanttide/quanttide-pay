// Package status 提供支付单、退款单状态契约。
// 取值使用英文小写字符串，便于直接存入数据库并与渠道报文流转。
// 包内按模块划分文件：payment.go 为支付状态，refund.go 为退款状态；
// 每模块含枚举、渠道原始码解析、合法性校验三部分。
package status

import "fmt"

// PaymentStatus 是支付单状态。
type PaymentStatus string

const (
	PaymentStatusCreated   PaymentStatus = "created"   // 已创建，待支付
	PaymentStatusPaying    PaymentStatus = "paying"    // 支付处理中
	PaymentStatusSucceeded PaymentStatus = "succeeded" // 支付成功
	PaymentStatusFailed    PaymentStatus = "failed"    // 支付失败
	PaymentStatusClosed    PaymentStatus = "closed"    // 已关闭（超时未支付或主动取消）
	PaymentStatusRefunding PaymentStatus = "refunding" // 退款处理中
	PaymentStatusRefunded  PaymentStatus = "refunded"  // 已全额退款
)

// wechatTradeState 微信支付 TradeState 原始码到统一支付状态的映射。
// 未知码不在此表内，由 ParseWechatTradeState 显式报错（新渠道状态出现时暴露，不静默降级）。
var wechatTradeState = map[string]PaymentStatus{
	"SUCCESS":    PaymentStatusSucceeded,
	"REFUND":     PaymentStatusRefunding,
	"NOTPAY":     PaymentStatusCreated,
	"USERPAYING": PaymentStatusPaying,
	"CLOSED":     PaymentStatusClosed,
	"PAYERROR":   PaymentStatusFailed,
}

// alipayTradeStatus 支付宝 trade_status 原始码到统一支付状态的映射。
var alipayTradeStatus = map[string]PaymentStatus{
	"TRADE_SUCCESS":  PaymentStatusSucceeded,
	"TRADE_FINISHED": PaymentStatusSucceeded,
	"WAIT_BUYER_PAY": PaymentStatusCreated,
	"TRADE_CLOSED":   PaymentStatusClosed,
}

// ParseWechatTradeState 将微信支付 TradeState 原始码解析为统一支付状态。
// 未知码返回错误，不用 UNKNOWN 兜底。
func ParseWechatTradeState(code string) (PaymentStatus, error) {
	if s, ok := wechatTradeState[code]; ok {
		return s, nil
	}
	return "", fmt.Errorf("status: unknown wechat trade state %q", code)
}

// ParseAlipayTradeStatus 将支付宝 trade_status 原始码解析为统一支付状态。
// 未知码返回错误，不用 UNKNOWN 兜底。
func ParseAlipayTradeStatus(code string) (PaymentStatus, error) {
	if s, ok := alipayTradeStatus[code]; ok {
		return s, nil
	}
	return "", fmt.Errorf("status: unknown alipay trade status %q", code)
}

// IsValidPaymentStatus 报告 s 是否为合法支付状态，供存库前防御校验。
func IsValidPaymentStatus(s PaymentStatus) bool {
	switch s {
	case PaymentStatusCreated, PaymentStatusPaying, PaymentStatusSucceeded,
		PaymentStatusFailed, PaymentStatusClosed, PaymentStatusRefunding, PaymentStatusRefunded:
		return true
	}
	return false
}
