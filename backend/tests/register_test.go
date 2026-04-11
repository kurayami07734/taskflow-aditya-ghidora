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

func TestRegisterIntegration(t *testing.T) {
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

	t.Run("valid registration returns 201", func(t *testing.T) {
		body := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", rr.Code, rr.Body.String())
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

	t.Run("missing name returns 400", func(t *testing.T) {
		body := map[string]string{
			"email":    "jane@example.com",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("missing email returns 400", func(t *testing.T) {
		body := map[string]string{
			"name":     "Jane Doe",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("invalid email format returns 400", func(t *testing.T) {
		body := map[string]string{
			"name":     "Jane Doe",
			"email":    "not-an-email",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("short password returns 400", func(t *testing.T) {
		body := map[string]string{
			"name":     "Jane Doe",
			"email":    "jane@example.com",
			"password": "short",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("duplicate email returns 400", func(t *testing.T) {
		body := map[string]string{
			"name":     "Duplicate User",
			"email":    "john@example.com",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}

		expectedMsg := "User already registered!"
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedMsg)) {
			t.Errorf("Expected '%s' in body, got: %s", expectedMsg, rr.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}
