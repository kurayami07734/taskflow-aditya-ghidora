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

func TestDeleteTaskIntegration(t *testing.T) {
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

	t.Run("delete own task returns 204", func(t *testing.T) {
		token := registerAndLogin("deltask1@example.com", "secret123")
		projectID := createProject(token, "Project", "Description")
		taskID := createTask(token, projectID, "Task to delete")

		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", rr.Code)
		}
	})

	t.Run("delete another users task returns 403", func(t *testing.T) {
		token1 := registerAndLogin("deltask2a@example.com", "secret123")
		token2 := registerAndLogin("deltask2b@example.com", "secret123")

		projectID := createProject(token1, "Project", "Description")
		taskID := createTask(token1, projectID, "Task")

		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
		req.Header.Set("Authorization", "Bearer "+token2)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", rr.Code)
		}
	})

	t.Run("delete non-existent task returns 404", func(t *testing.T) {
		token := registerAndLogin("deltask3@example.com", "secret123")

		req, _ := http.NewRequest("DELETE", "/tasks/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("delete task without auth returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/tasks/some-id", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
