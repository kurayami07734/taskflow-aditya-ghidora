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

func TestSearchUsersIntegration(t *testing.T) {
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

	t.Run("search by name returns matching users", func(t *testing.T) {
		regBody := map[string]string{
			"name":     "John Searchable",
			"email":    "john.searchable@test.com",
			"password": "password123",
		}
		regBodyBytes, _ := json.Marshal(regBody)
		regReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(regBodyBytes))
		regReq.Header.Set("Content-Type", "application/json")
		regRR := httptest.NewRecorder()
		router.ServeHTTP(regRR, regReq)

		req, _ := http.NewRequest("GET", "/users/search?q=John", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var users []interface{}
		json.Unmarshal(rr.Body.Bytes(), &users)
		if len(users) == 0 {
			t.Error("Expected at least one user in results")
		}
	})

	t.Run("search by email returns matching users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/search?q=john.searchable@test.com", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var users []interface{}
		json.Unmarshal(rr.Body.Bytes(), &users)
		if len(users) == 0 {
			t.Error("Expected at least one user in results")
		}
	})

	t.Run("search with no matches returns empty array", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/search?q=nonexistentuser123", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var users []interface{}
		json.Unmarshal(rr.Body.Bytes(), &users)
		if len(users) != 0 {
			t.Errorf("Expected empty array, got %d users", len(users))
		}
	})

	t.Run("missing query parameter returns 400", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/search", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("search is case insensitive", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/search?q=john", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var users []interface{}
		json.Unmarshal(rr.Body.Bytes(), &users)
		if len(users) == 0 {
			t.Error("Expected at least one user in results")
		}
	})
}
