package status

import "fmt"

// RefundStatus 是退款单状态。
type RefundStatus string

const (
	RefundStatusCreated    RefundStatus = "created"    // 已创建，待处理
	RefundStatusProcessing RefundStatus = "processing" // 退款处理中
	RefundStatusSucceeded  RefundStatus = "succeeded"  // 退款成功
	RefundStatusFailed     RefundStatus = "failed"     // 退款失败
)

// wechatRefundStatus 微信退款 status 原始码到统一退款状态的映射。
// 微信退款终态 CLOSED/ABNORMAL 归一为 failed（RefundStatus 无 closed 枚举）。
var wechatRefundStatus = map[string]RefundStatus{
	"SUCCESS":    RefundStatusSucceeded,
	"PROCESSING": RefundStatusProcessing,
	"CLOSED":     RefundStatusFailed,
	"ABNORMAL":   RefundStatusFailed,
}

// ParseWechatRefundStatus 将微信退款 status 原始码解析为统一退款状态。
// 未知码返回错误，不用兜底。
func ParseWechatRefundStatus(code string) (RefundStatus, error) {
	if s, ok := wechatRefundStatus[code]; ok {
		return s, nil
	}
	return "", fmt.Errorf("status: unknown wechat refund status %q", code)
}

// IsValidRefundStatus 报告 s 是否为合法退款状态，供存库前防御校验。
func IsValidRefundStatus(s RefundStatus) bool {
	switch s {
	case RefundStatusCreated, RefundStatusProcessing, RefundStatusSucceeded, RefundStatusFailed:
		return true
	}
	return false
}
