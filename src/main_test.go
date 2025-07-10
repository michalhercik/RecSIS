package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//================================================================================
// Test Entry Points
//================================================================================

func TestHomeServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"Root", "/", http.StatusOK},
		{"CS Root", "/cs/", http.StatusOK},
		{"EN Root", "/en/", http.StatusOK},
		{"Home", "/home/", http.StatusOK},
		{"CS Home", "/cs/home/", http.StatusOK},
		{"EN Home", "/en/home/", http.StatusOK},

		// Errors
		{"Not Found", "/homer/", http.StatusNotFound},
		{"Unknown Path", "/lorem", http.StatusNotFound},
	}
	runTests(t, tests)
}

func TestCoursesServer(t *testing.T)      { runTests(t, []testCase{ /* TODO */ }) }
func TestCourseDetailServer(t *testing.T) { runTests(t, []testCase{ /* TODO */ }) }
func TestBlueprintServer(t *testing.T)    { runTests(t, []testCase{ /* TODO */ }) }
func TestDegreePlanServer(t *testing.T)   { runTests(t, []testCase{ /* TODO */ }) }

//================================================================================
// Test Utilities
//================================================================================

type testCase struct {
	name string
	url  string
	want int
}

func runTests(t *testing.T, tests []testCase) {
	ts := setupTestServer()
	defer ts.Close()

	client := ts.Client()
	sessionCookie := setupTestUser(ts, t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+test.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.AddCookie(sessionCookie)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, test.want, resp.StatusCode)
		})
	}
}

func setupTestServer() *httptest.Server {
	conf := configFrom("./config.dev.toml")
	handler := setupHandler(conf)
	return httptest.NewTLSServer(handler)
}

func setupTestUser(ts *httptest.Server, t *testing.T) *http.Cookie {
	insecureClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	form := url.Values{}
	form.Set("user-id", "testuser")
	form.Set("service", ts.URL+"/cas/login")

	req, err := http.NewRequest("POST", "https://localhost:8001/cas/login", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal("Failed to create CAS login request:", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := insecureClient.Do(req)
	if err != nil {
		t.Fatal("CAS login failed:", err)
	}
	defer res.Body.Close()

	redirectURL := res.Header.Get("Location")
	if redirectURL == "" {
		t.Fatal("No redirect URL from CAS login")
	}

	req2, _ := http.NewRequest("GET", redirectURL, nil)
	res2, err := insecureClient.Do(req2)
	if err != nil {
		t.Fatal("Redirect to app login failed:", err)
	}
	defer res2.Body.Close()

	for _, c := range res2.Cookies() {
		if c.Name == "recsis_session_key" {
			return c
		}
	}
	t.Fatal("No session cookie returned from login")
	return nil
}
