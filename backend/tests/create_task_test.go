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

func TestCreateTaskIntegration(t *testing.T) {
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

	t.Run("create task with valid payload returns 201", func(t *testing.T) {
		token := registerAndLogin("createtask1@example.com", "secret123")
		projectID := createProject(token, "Test Project", "Description")

		body := map[string]interface{}{
			"title":       "Design homepage",
			"description": "Create mockups for the new homepage",
			"priority":    "high",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["id"] == nil {
			t.Error("Expected id in response")
		}
		if resp["title"] != "Design homepage" {
			t.Errorf("Expected title Design homepage, got %s", resp["title"])
		}
		if resp["status"] != "todo" {
			t.Errorf("Expected status todo, got %s", resp["status"])
		}
	})

	t.Run("create task without optional fields returns 201", func(t *testing.T) {
		token := registerAndLogin("createtask2@example.com", "secret123")
		projectID := createProject(token, "Test Project", "Description")

		body := map[string]string{
			"title": "Simple task",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["priority"] != "medium" {
			t.Errorf("Expected priority medium, got %s", resp["priority"])
		}
	})

	t.Run("create task with missing title returns 400", func(t *testing.T) {
		token := registerAndLogin("createtask3@example.com", "secret123")
		projectID := createProject(token, "Test Project", "Description")

		body := map[string]string{
			"description": "Some description",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("create task with title less than 3 chars returns 400", func(t *testing.T) {
		token := registerAndLogin("createtask4@example.com", "secret123")
		projectID := createProject(token, "Test Project", "Description")

		body := map[string]string{
			"title": "ab",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("create task on another users project returns 403", func(t *testing.T) {
		token1 := registerAndLogin("createtask5a@example.com", "secret123")
		token2 := registerAndLogin("createtask5b@example.com", "secret123")

		projectID := createProject(token1, "User1 Project", "Description")

		body := map[string]string{
			"title": "Hacked task",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token2)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", rr.Code)
		}
	})

	t.Run("create task on non-existent project returns 404", func(t *testing.T) {
		token := registerAndLogin("createtask6@example.com", "secret123")

		body := map[string]string{
			"title": "Task",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/00000000-0000-0000-0000-000000000000/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("create task without auth returns 401", func(t *testing.T) {
		body := map[string]string{
			"title": "Task",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/some-id/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
