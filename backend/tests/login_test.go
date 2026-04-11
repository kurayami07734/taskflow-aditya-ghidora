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

func TestLoginIntegration(t *testing.T) {
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

	registerUser := func(email, password string) {
		body := map[string]string{
			"name":     "Test User",
			"email":    email,
			"password": password,
		}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}

	registerUser("login@example.com", "secret123")

	t.Run("valid login returns 200", func(t *testing.T) {
		body := map[string]string{
			"email":    "login@example.com",
			"password": "secret123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["token"] == nil {
			t.Error("Expected token in response")
		}
		if resp["user"] == nil {
			t.Error("Expected user in response")
		}
	})

	t.Run("wrong email returns 400", func(t *testing.T) {
		body := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "secret123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}

		expectedMsg := "Invalid credentials"
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedMsg)) {
			t.Errorf("Expected '%s' in body, got: %s", expectedMsg, rr.Body.String())
		}
	})

	t.Run("wrong password returns 400", func(t *testing.T) {
		body := map[string]string{
			"email":    "login@example.com",
			"password": "wrongpassword",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}

		expectedMsg := "Invalid credentials"
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedMsg)) {
			t.Errorf("Expected '%s' in body, got: %s", expectedMsg, rr.Body.String())
		}
	})

	t.Run("missing email returns 400", func(t *testing.T) {
		body := map[string]string{
			"password": "secret123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("missing password returns 400", func(t *testing.T) {
		body := map[string]string{
			"email": "login@example.com",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}
