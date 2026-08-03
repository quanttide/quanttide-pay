// Package order 提供支付订单号生成。
package order

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

// 随机序列位数
const randomDigits = 10

// 随机序列范围（10^10）
const randomBound = 10000000000

// Generate 生成支付订单号：{prefix}{yyyyMMddHHmmss}{10 位随机数字}，
// 例如 PAY202608031430001234567890。随机序列使用密码学安全随机数，避免可预测性。
func Generate(prefix string, now time.Time) (string, error) {
	var buf [4]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("order: 生成随机序列失败: %w", err)
	}
	n := uint64(binary.BigEndian.Uint32(buf[:]))
	return fmt.Sprintf("%s%s%010d", prefix, now.Format("20060102150405"), n%randomBound), nil
}

// GenerateNow 使用当前时间生成支付订单号。
func GenerateNow(prefix string) (string, error) {
	return Generate(prefix, time.Now())
}
