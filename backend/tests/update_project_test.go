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

func TestUpdateProjectIntegration(t *testing.T) {
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

	t.Run("update project name returns 200", func(t *testing.T) {
		token := registerAndLogin("updproj1@example.com", "secret123")
		projectID := createProject(token, "Original Name", "Original description")

		body := map[string]string{
			"name": "Updated Name",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/"+projectID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["name"] != "Updated Name" {
			t.Errorf("Expected name Updated Name, got %s", resp["name"])
		}
	})

	t.Run("update project description returns 200", func(t *testing.T) {
		token := registerAndLogin("updproj2@example.com", "secret123")
		projectID := createProject(token, "Project", "Original description")

		body := map[string]string{
			"description": "Updated description",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/"+projectID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["description"] != "Updated description" {
			t.Errorf("Expected description Updated description, got %s", resp["description"])
		}
	})

	t.Run("update both name and description returns 200", func(t *testing.T) {
		token := registerAndLogin("updproj3@example.com", "secret123")
		projectID := createProject(token, "Original Name", "Original description")

		body := map[string]string{
			"name":        "New Name",
			"description": "New description",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/"+projectID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["name"] != "New Name" || resp["description"] != "New description" {
			t.Errorf("Expected updated name and description")
		}
	})

	t.Run("update name less than 3 chars returns 400", func(t *testing.T) {
		token := registerAndLogin("updproj4@example.com", "secret123")
		projectID := createProject(token, "Project", "Description")

		body := map[string]string{
			"name": "ab",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/"+projectID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("update another users project returns 403", func(t *testing.T) {
		token1 := registerAndLogin("updproj5a@example.com", "secret123")
		token2 := registerAndLogin("updproj5b@example.com", "secret123")

		projectID := createProject(token1, "User1 Project", "Description")

		body := map[string]string{
			"name": "Hacked Name",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/"+projectID, bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token2)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", rr.Code)
		}
	})

	t.Run("update non-existent project returns 404", func(t *testing.T) {
		token := registerAndLogin("updproj6@example.com", "secret123")

		body := map[string]string{
			"name": "New Name",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/00000000-0000-0000-0000-000000000000", bytes.NewReader(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("update project without auth returns 401", func(t *testing.T) {
		body := map[string]string{
			"name": "New Name",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PATCH", "/projects/some-id", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
