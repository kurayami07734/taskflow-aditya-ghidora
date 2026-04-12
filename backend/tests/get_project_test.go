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

func TestGetProjectIntegration(t *testing.T) {
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

	createProject := func(token, name, description string) string {
		body := map[string]string{
			"name":        name,
			"description": description,
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		return resp["id"].(string)
	}

	createTask := func(token, projectID, title string) {
		body := map[string]interface{}{"title": title}
		bodyBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}

	t.Run("get existing project returns 200", func(t *testing.T) {
		token := registerAndLogin("getproj1@example.com", "secret123")
		projectID := createProject(token, "Test Project", "Test description")

		req, _ := http.NewRequest("GET", "/projects/"+projectID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["id"] == nil {
			t.Error("Expected id in response")
		}
		if resp["name"] != "Test Project" {
			t.Errorf("Expected name Test Project, got %s", resp["name"])
		}
		if resp["pagination"] == nil {
			t.Error("Expected pagination in response")
		}
	})

	t.Run("get project with pagination returns correct data", func(t *testing.T) {
		token := registerAndLogin("getproj5@example.com", "secret123")
		projectID := createProject(token, "Test Project", "Test description")

		for i := 0; i < 5; i++ {
			createTask(token, projectID, "Task "+string(rune('1'+i)))
		}

		req, _ := http.NewRequest("GET", "/projects/"+projectID+"?page=1&limit=3", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		tasks := resp["tasks"].([]interface{})
		if len(tasks) != 3 {
			t.Errorf("Expected 3 tasks (limit), got %d", len(tasks))
		}
		pagination := resp["pagination"].(map[string]interface{})
		if pagination["total"].(float64) != 5 {
			t.Errorf("Expected total 5, got %v", pagination["total"])
		}
		if pagination["total_pages"].(float64) != 2 {
			t.Errorf("Expected total_pages 2, got %v", pagination["total_pages"])
		}
	})

	t.Run("get non-existent project returns 404", func(t *testing.T) {
		token := registerAndLogin("getproj2@example.com", "secret123")

		req, _ := http.NewRequest("GET", "/projects/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("get another users project returns 403", func(t *testing.T) {
		token1 := registerAndLogin("getproj3a@example.com", "secret123")
		token2 := registerAndLogin("getproj3b@example.com", "secret123")

		projectID := createProject(token1, "User1 Project", "Description")

		req, _ := http.NewRequest("GET", "/projects/"+projectID, nil)
		req.Header.Set("Authorization", "Bearer "+token2)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", rr.Code)
		}
	})

	t.Run("get project with invalid id returns 400", func(t *testing.T) {
		token := registerAndLogin("getproj4@example.com", "secret123")

		req, _ := http.NewRequest("GET", "/projects/invalid-uuid", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("get project without auth returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/some-id", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
