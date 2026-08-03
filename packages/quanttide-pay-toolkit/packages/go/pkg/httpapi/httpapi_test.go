package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, http.StatusCreated, map[string]string{"id": "acc_1"})
	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"id":"acc_1"`) {
		t.Errorf("body = %s", rec.Body.String())
	}
}

func TestWriteJSON_EncodeError(t *testing.T) {
	// 不可序列化值走错误分支（仅记日志，不影响状态码）
	rec := httptest.NewRecorder()
	WriteJSON(rec, http.StatusOK, make(chan int))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusBadRequest, "bad request")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
	if rec.Body.String() != "{\"error\":\"bad request\"}\n" {
		t.Errorf("body = %s", rec.Body.String())
	}
}

func TestWriteServiceError(t *testing.T) {
	mapper := Mapper(func(err error) int {
		if err.Error() == "not found" {
			return http.StatusNotFound
		}
		return 0
	})

	// 映射命中
	rec := httptest.NewRecorder()
	WriteServiceError(rec, &errNotFound{}, mapper)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}

	// 未识别 → 500
	rec2 := httptest.NewRecorder()
	WriteServiceError(rec2, &errBoom{}, mapper)
	if rec2.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec2.Code)
	}

	// 无 mapper → 500
	rec3 := httptest.NewRecorder()
	WriteServiceError(rec3, &errBoom{}, nil)
	if rec3.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec3.Code)
	}
}

type errNotFound struct{}

func (e *errNotFound) Error() string { return "not found" }

type errBoom struct{}

func (e *errBoom) Error() string { return "boom" }

func TestParsePagination(t *testing.T) {
	cases := []struct {
		query      string
		wantLimit  int
		wantOffset int
	}{
		{"", 20, 0},
		{"limit=5", 5, 0},
		{"limit=1000", 100, 0},
		{"limit=0", 20, 0},
		{"limit=abc", 20, 0},
		{"offset=10", 20, 10},
		{"offset=-1", 20, 0},
		{"limit=50&offset=25", 50, 25},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/x?"+c.query, nil)
		limit, offset := ParsePagination(req)
		if limit != c.wantLimit || offset != c.wantOffset {
			t.Errorf("ParsePagination(%q) = %d,%d; want %d,%d", c.query, limit, offset, c.wantLimit, c.wantOffset)
		}
	}
}
