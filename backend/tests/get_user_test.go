package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kurayami07734/taskflow-aditya-ghidora/src/routers"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

func TestGetUserIntegration(t *testing.T) {
	testDB := SetupPostgres(t)
	defer testDB.Close()

	cfg := utils.Config{
		Port:      8080,
		JwtSecret: "test-jwt-secret-key-for-integration-tests",
		Db: struct {
			Port     int    `env:"DB_PORT" envDefault:"5432"`
			Name     string `env:"DB_NAME"`
			Host     string `env:"DB_HOST"`
			User     string `env:"DB_USER"`
			Password string `env:"DB_PASSWORD"`
			SslMode  string `env:"DB_SSL_MODE"`
		}{
			Port:     testDB.Config.Port,
			Name:     testDB.Config.Name,
			Host:     testDB.Config.Host,
			User:     testDB.Config.User,
			Password: testDB.Config.Password,
			SslMode:  "disable",
		},
	}

	router := routers.CreateBaseRouter(testDB.DB, cfg)

	t.Run("valid user ID returns user", func(t *testing.T) {
		regBody := map[string]string{
			"name":     "John Doe",
			"email":    "johndoe@test.com",
			"password": "password123",
		}
		regBodyBytes, _ := json.Marshal(regBody)
		regReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(regBodyBytes))
		regReq.Header.Set("Content-Type", "application/json")
		regRR := httptest.NewRecorder()
		router.ServeHTTP(regRR, regReq)

		var regResp map[string]interface{}
		json.Unmarshal(regRR.Body.Bytes(), &regResp)
		user := regResp["user"].(map[string]interface{})
		userID := user["id"].(string)

		req, _ := http.NewRequest("GET", "/users/"+userID, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["id"] == nil {
			t.Error("Expected user id in response")
		}
		if resp["name"] != "John Doe" {
			t.Errorf("Expected name 'John Doe', got %v", resp["name"])
		}
		if resp["email"] != "johndoe@test.com" {
			t.Errorf("Expected email 'johndoe@test.com', got %v", resp["email"])
		}
	})

	t.Run("invalid UUID returns 400", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/invalid-uuid", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("non-existent user returns 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/00000000-0000-0000-0000-000000000000", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})
}
