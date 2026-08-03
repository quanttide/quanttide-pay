// Package status 提供支付单、退款单状态枚举。
// 取值使用英文小写字符串，便于直接存入数据库并与渠道报文流转。
package status

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

// RefundStatus 是退款单状态。
type RefundStatus string

const (
	RefundStatusCreated    RefundStatus = "created"    // 已创建，待处理
	RefundStatusProcessing RefundStatus = "processing" // 退款处理中
	RefundStatusSucceeded  RefundStatus = "succeeded"  // 退款成功
	RefundStatusFailed     RefundStatus = "failed"     // 退款失败
)
