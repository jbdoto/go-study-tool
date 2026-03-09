package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jbdoto/go-study-tool/internal/challenges"
	"github.com/jbdoto/go-study-tool/internal/questions"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	sets := []challenges.ChallengeSet{{
		Slug:  "test-set",
		Title: "Test Set",
		Challenges: []challenges.CodeChallenge{{
			ID:       "t1",
			Title:    "Test",
			Scaffold: "package main",
			Solution: "package main\n\nfunc main() {}",
			Hints:    []string{"hint one", "hint two"},
		}},
	}}
	templateFS := os.DirFS("../../cmd/server")
	srv, err := New([]questions.Chapter{}, sets, templateFS)
	if err != nil {
		t.Fatal(err)
	}
	return srv
}

func postCheck(t *testing.T, srv http.Handler, code string) *http.Response {
	t.Helper()
	form := url.Values{"code": {code}}
	req := httptest.NewRequest(http.MethodPost, "/challenges/test-set/check", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	return w.Result()
}

func TestCheckSyntax_ValidCode(t *testing.T) {
	srv := newTestServer(t)
	resp := postCheck(t, srv.Routes(), "package main\n\nfunc main() {}\n")
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if !strings.Contains(string(body), "Syntax OK") {
		t.Errorf("expected 'Syntax OK' in body, got: %s", body)
	}
}

func TestCheckSyntax_GarbageInput(t *testing.T) {
	srv := newTestServer(t)
	cases := []struct {
		name string
		code string
	}{
		{"random text", "this is not go code at all!!!"},
		{"missing package", "func main() {}"},
		{"unclosed brace", "package main\n\nfunc main() {"},
		{"invalid token", "package main\n\n@@@"},
		{"empty string", ""},
		{"partial declaration", "package main\n\nfunc ("},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := postCheck(t, srv.Routes(), tc.code)
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
			if strings.Contains(string(body), "Syntax OK") {
				t.Errorf("expected syntax error for %q, but got 'Syntax OK'", tc.code)
			}
			if !strings.Contains(string(body), "check-error") {
				t.Errorf("expected check-error class in response for %q, got: %s", tc.code, body)
			}
		})
	}
}
