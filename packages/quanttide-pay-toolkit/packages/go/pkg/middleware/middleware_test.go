package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_PassesThrough(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	ts := httptest.NewServer(Logging(inner))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want 201", resp.StatusCode)
	}
}

func TestStatusRecorder_ExplicitAndDefault(t *testing.T) {
	// 显式 WriteHeader：记录该状态码
	rec1 := httptest.NewRecorder()
	r1 := &statusRecorder{ResponseWriter: rec1}
	r1.WriteHeader(http.StatusNotFound)
	if r1.status != http.StatusNotFound {
		t.Errorf("status = %d, want 404", r1.status)
	}

	// 未调用 WriteHeader：默认 200（status 初始化为 OK，Write 的隐式 WriteHeader 不经过包装）
	rec2 := httptest.NewRecorder()
	r2 := &statusRecorder{ResponseWriter: rec2, status: http.StatusOK}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler.ServeHTTP(r2, httptest.NewRequest(http.MethodGet, "/", nil))
	if r2.status != http.StatusOK {
		t.Errorf("status = %d, want 200", r2.status)
	}
}
