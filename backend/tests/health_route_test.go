package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kurayami07734/taskflow-aditya-ghidora/src/routers"
)

func TestHealthRoute(t *testing.T) {
	r := routers.CreateBaseRouter()
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}
}
