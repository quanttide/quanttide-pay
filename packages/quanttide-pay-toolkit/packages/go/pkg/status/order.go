package status

// OrderStatus 是订单状态。
type OrderStatus string

const (
	OrderStatusCreated OrderStatus = "created" // 已创建，待结算
	OrderStatusSettled OrderStatus = "settled" // 已结算
)

// IsValidOrderStatus 报告订单状态是否为已知状态（存库前防御校验）。
func IsValidOrderStatus(s OrderStatus) bool {
	switch s {
	case OrderStatusCreated, OrderStatusSettled:
		return true
	}
	return false
}
