package order

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateFormat(t *testing.T) {
	now := time.Date(2026, 8, 3, 14, 30, 0, 0, time.UTC)
	orderNo, err := Generate("PAY", now)
	if err != nil {
		t.Fatalf("Generate() 返回错误: %v", err)
	}
	if !strings.HasPrefix(orderNo, "PAY20260803143000") {
		t.Fatalf("Generate() = %q, 缺少前缀与时间戳", orderNo)
	}
	if len(orderNo) != 27 {
		t.Fatalf("Generate() 长度 = %d, 期望 27", len(orderNo))
	}
}

func TestGenerateUniqueness(t *testing.T) {
	now := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	seen := make(map[string]bool, 100)
	for i := 0; i < 100; i++ {
		orderNo, err := Generate("PAY", now)
		if err != nil {
			t.Fatalf("Generate() 返回错误: %v", err)
		}
		if seen[orderNo] {
			t.Fatalf("订单号重复: %s", orderNo)
		}
		seen[orderNo] = true
	}
}

func TestGeneratePrefix(t *testing.T) {
	now := time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC)
	orderNo, err := Generate("REF", now)
	if err != nil {
		t.Fatalf("Generate() 返回错误: %v", err)
	}
	if !strings.HasPrefix(orderNo, "REF") {
		t.Fatalf("Generate() = %q, 期望以 REF 开头", orderNo)
	}
}
