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

func TestUpdateTaskIntegration(t *testing.T) {
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

	createTask := func(token, projectID, title string) string {
		body := map[string]string{
			"title": title,
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/projects/"+projectID+"/tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		return resp["id"].(string)
	}

	t.Run("update task title returns 200", func(t *testing.T) {
		token := registerAndLogin("updatetask1@example.com", "secret123")
		projectID := createProject(token, "Project", "Description")
		taskID := createTask(token, projectID, "Original title")

		body := map[string]string{
			"title": "Updated title",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/"+taskID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["title"] != "Updated title" {
			t.Errorf("Expected title Updated title, got %s", resp["title"])
		}
	})

	t.Run("update task status returns 200", func(t *testing.T) {
		token := registerAndLogin("updatetask2@example.com", "secret123")
		projectID := createProject(token, "Project", "Description")
		taskID := createTask(token, projectID, "Task")

		body := map[string]string{
			"status": "done",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/"+taskID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["status"] != "done" {
			t.Errorf("Expected status done, got %s", resp["status"])
		}
	})

	t.Run("update task priority returns 200", func(t *testing.T) {
		token := registerAndLogin("updatetask3@example.com", "secret123")
		projectID := createProject(token, "Project", "Description")
		taskID := createTask(token, projectID, "Task")

		body := map[string]string{
			"priority": "high",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/"+taskID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["priority"] != "high" {
			t.Errorf("Expected priority high, got %s", resp["priority"])
		}
	})

	t.Run("update task with invalid status returns 400", func(t *testing.T) {
		token := registerAndLogin("updatetask4@example.com", "secret123")
		projectID := createProject(token, "Project", "Description")
		taskID := createTask(token, projectID, "Task")

		body := map[string]string{
			"status": "invalid",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/"+taskID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("update task on another users project returns 403", func(t *testing.T) {
		token1 := registerAndLogin("updatetask5a@example.com", "secret123")
		token2 := registerAndLogin("updatetask5b@example.com", "secret123")

		projectID := createProject(token1, "Project", "Description")
		taskID := createTask(token1, projectID, "Task")

		body := map[string]string{
			"title": "Hacked title",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/"+taskID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token2)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", rr.Code)
		}
	})

	t.Run("update non-existent task returns 404", func(t *testing.T) {
		token := registerAndLogin("updatetask6@example.com", "secret123")

		body := map[string]string{
			"title": "New title",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/00000000-0000-0000-0000-000000000000", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("update task without auth returns 401", func(t *testing.T) {
		body := map[string]string{
			"title": "New title",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/tasks/some-id", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
