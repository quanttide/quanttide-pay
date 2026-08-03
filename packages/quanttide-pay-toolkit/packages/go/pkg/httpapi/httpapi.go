// Package httpapi 提供 HTTP JSON API 公共件：统一响应、服务错误映射与分页。
//
// 服务端各 transport 模块共用：WriteJSON/WriteError 统一响应格式
// （错误体固定为 {"error": "..."}）；WriteServiceError 按模块注册的
// Mapper 映射业务错误为状态码，未识别错误记日志并返回 500。
package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// WriteJSON 以 JSON 格式写入响应。
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("httpapi: write json response: %v", err)
	}
}

// WriteError 写入错误响应：{"error": msg}。
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// Mapper 错误 → HTTP 状态码映射；返回 0 表示未识别（写 500）。
type Mapper func(err error) int

// WriteServiceError 按 mapper 映射服务错误并写响应；未识别错误记日志并返回 500。
func WriteServiceError(w http.ResponseWriter, err error, mapper Mapper) {
	status := http.StatusInternalServerError
	if mapper != nil {
		if s := mapper(err); s != 0 {
			status = s
		}
	}
	if status == http.StatusInternalServerError {
		log.Printf("httpapi: unhandled error: %v", err)
	}
	WriteError(w, status, http.StatusText(status))
}

// ParsePagination 解析 limit/offset 查询参数：limit 默认 20、最大 100，offset 默认 0。
func ParsePagination(r *http.Request) (limit, offset int) {
	limit, offset = 20, 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > 100 {
		limit = 100
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}
