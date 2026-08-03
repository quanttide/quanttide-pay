package status

// CouponStatus 是优惠券状态。
type CouponStatus string

const (
	CouponStatusIssued  CouponStatus = "issued"  // 已发放，可用
	CouponStatusUsed    CouponStatus = "used"    // 已使用（核销于结算）
	CouponStatusExpired CouponStatus = "expired" // 已过期（查询时惰性流转）
)

// IsValidCouponStatus 报告券状态是否为已知状态（存库前防御校验）。
func IsValidCouponStatus(s CouponStatus) bool {
	switch s {
	case CouponStatusIssued, CouponStatusUsed, CouponStatusExpired:
		return true
	}
	return false
}
