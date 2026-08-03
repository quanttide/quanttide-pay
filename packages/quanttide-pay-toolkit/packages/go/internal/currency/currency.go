// Package currency 提供 ISO 4217 币种代码的校验实现，
// 供库内公共包复用，外部不可导入。
package currency

import "regexp"

var pattern = regexp.MustCompile(`^[A-Z]{3}$`)

// Valid 报告币种代码是否为合法的 ISO 4217 三字母大写形式（如 CNY、USD）。
func Valid(code string) bool {
	return pattern.MatchString(code)
}
