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

func TestProjectsIntegration(t *testing.T) {
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

	registerAndLogin := func(email, password string) string {
		regBody := map[string]string{
			"name":     "Test User",
			"email":    email,
			"password": password,
		}
		regBodyBytes, _ := json.Marshal(regBody)
		regReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(regBodyBytes))
		regReq.Header.Set("Content-Type", "application/json")
		regRR := httptest.NewRecorder()
		router.ServeHTTP(regRR, regReq)

		loginBody := map[string]string{
			"email":    email,
			"password": password,
		}
		loginBodyBytes, _ := json.Marshal(loginBody)
		loginReq, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(loginBodyBytes))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRR := httptest.NewRecorder()
		router.ServeHTTP(loginRR, loginReq)

		var resp map[string]interface{}
		json.Unmarshal(loginRR.Body.Bytes(), &resp)
		return resp["token"].(string)
	}

	t.Run("authenticated request returns 200", func(t *testing.T) {
		token := registerAndLogin("projects1@example.com", "secret123")

		req, _ := http.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["projects"] == nil {
			t.Error("Expected projects in response")
		}
		if resp["pagination"] == nil {
			t.Error("Expected pagination in response")
		}
	})

	t.Run("authenticated request returns empty projects list", func(t *testing.T) {
		token := registerAndLogin("projects2@example.com", "secret123")

		req, _ := http.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		projects := resp["projects"].([]interface{})
		if len(projects) != 0 {
			t.Errorf("Expected empty projects list, got %d", len(projects))
		}
	})

	t.Run("authenticated request with pagination returns correct data", func(t *testing.T) {
		token := registerAndLogin("projects4@example.com", "secret123")

		for i := 0; i < 5; i++ {
			body := map[string]string{
				"name":        "Project " + string(rune('1'+i)),
				"description": "Description",
			}
			bodyBytes, _ := json.Marshal(body)
			req, _ := http.NewRequest("POST", "/projects", bytes.NewReader(bodyBytes))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
		}

		req, _ := http.NewRequest("GET", "/projects?page=1&limit=3", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		projects := resp["projects"].([]interface{})
		if len(projects) != 3 {
			t.Errorf("Expected 3 projects (limit), got %d", len(projects))
		}
		pagination := resp["pagination"].(map[string]interface{})
		if pagination["total"].(float64) != 5 {
			t.Errorf("Expected total 5, got %v", pagination["total"])
		}
		if pagination["total_pages"].(float64) != 2 {
			t.Errorf("Expected total_pages 2, got %v", pagination["total_pages"])
		}
	})

	t.Run("missing auth header returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})

	t.Run("malformed auth header returns 401", func(t *testing.T) {
		token := registerAndLogin("projects3@example.com", "secret123")

		req, _ := http.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
