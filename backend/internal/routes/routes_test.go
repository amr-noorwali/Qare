package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"qare/backend/internal/handlers"
	"qare/backend/internal/repositories"
	"testing"
	"time"
)

func TestCoreFlow(t *testing.T) {
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		t.Skip("MYSQL_TEST_DSN is not set")
	}
	db, err := repositories.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	server := New(handlers.New(db), "http://localhost:3000")
	request := func(method, path, body, token string, expected int) map[string]any {
		t.Helper()
		r := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		server.ServeHTTP(w, r)
		if w.Code != expected {
			t.Fatalf("%s %s: got %d, want %d: %s", method, path, w.Code, expected, w.Body.String())
		}
		var result any
		if w.Body.Len() > 0 && json.Unmarshal(w.Body.Bytes(), &result) != nil {
			t.Fatalf("invalid JSON: %s", w.Body.String())
		}
		if m, ok := result.(map[string]any); ok {
			return m
		}
		return nil
	}
	request("GET", "/api/books?q=الخيميائي", "", "", 200)
	email := fmt.Sprintf("reader-%d@example.com", time.Now().UnixNano())
	registered := request("POST", "/api/auth/register", fmt.Sprintf(`{"name":"قارئ جديد","email":%q,"password":"secret123"}`, email), "", 201)
	token := registered["token"].(string)
	defer db.Exec("DELETE FROM users WHERE email=?", email)
	request("POST", "/api/auth/login", fmt.Sprintf(`{"email":%q,"password":"secret123"}`, email), "", 200)
	request("PUT", "/api/books/1/review", `{"rating":5,"body":"كتاب رائع ومميز"}`, "", 401)
	request("PUT", "/api/books/999/review", `{"rating":5,"body":"كتاب رائع ومميز"}`, token, 404)
	request("PUT", "/api/books/1/review", `{"rating":0,"body":"قصير"}`, token, 400)
	request("DELETE", "/api/books/1/review", "", token, 404)
	request("PUT", "/api/books/1/review", `{"rating":5,"body":"كتاب رائع ومميز"}`, token, 200)
	request("PUT", "/api/books/1/review", `{"rating":4,"body":"قراءة ممتعة"}`, token, 200)
	request("PUT", "/api/library/999", `{"status":"reading"}`, token, 404)
	request("PUT", "/api/library/1", `{"status":"unknown"}`, token, 400)
	request("PUT", "/api/library/1", `{"status":"reading"}`, token, 200)
	result := request("GET", "/api/library", "", token, 200)
	_ = result
	request("DELETE", "/api/books/1/review", "", token, 204)
	request("DELETE", "/api/library/1", "", token, 204)
	request("POST", "/api/auth/logout", "", token, 204)
	request("GET", "/api/auth/me", "", token, 401)
}
