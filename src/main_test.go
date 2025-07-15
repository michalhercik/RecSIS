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

const (
	testUserID = "testuser"
)

//================================================================================
// Test Entry Points
//================================================================================

func TestHomeServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"root path should return 200", "GET", "/", http.StatusOK},
		{"cs root path should return 200", "GET", "/cs/", http.StatusOK},
		{"en root path should return 200", "GET", "/en/", http.StatusOK},
		{"home path should return 200", "GET", "/home/", http.StatusOK},
		{"cs home path should return 200", "GET", "/cs/home/", http.StatusOK},
		{"en home path should return 200", "GET", "/en/home/", http.StatusOK},

		// Errors
		{"root non-existent page should return 404", "GET", "/homer/", http.StatusNotFound},
		{"root non-existing page should return 404", "GET", "/lorem", http.StatusNotFound},
	}
	runTests(t, tests)
}

func TestCoursesServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"courses root path should return 200", "GET", "/courses/", http.StatusOK},
		{"courses cs root path should return 200", "GET", "/cs/courses/", http.StatusOK},
		{"courses en root path should return 200", "GET", "/en/courses/", http.StatusOK},
		{"courses search query should return 200", "GET", "/courses/?search=peska", http.StatusOK},
		{"courses search query with random string should return 200", "GET", "/courses/?search=asdfghjkl", http.StatusOK},
		{"courses search query with some course should return 200", "GET", "/courses/?search=NSWI120", http.StatusOK},
		{"courses search query with empty string should return 200", "GET", "/courses/?search=", http.StatusOK},
		{"courses pagination should return 200", "GET", "/courses/?page=2", http.StatusOK},
		{"courses load more courses should return 200", "GET", "/courses/?hitsPerPage=10", http.StatusOK},
		{"search courses should return 200", "GET", "/courses/search", http.StatusOK},
		{"search search query should return 200", "GET", "/courses/search?search=peska", http.StatusOK},
		{"search search query on search with random string should return 200", "GET", "/courses/search?search=asdfghjkl", http.StatusOK},
		{"search search query on search with some course should return 200", "GET", "/courses/search?search=NSWI120", http.StatusOK},
		{"search search query on search with empty string should return 200", "GET", "/courses/search?search=", http.StatusOK},
		{"search pagination on search should return 200", "GET", "/courses/search?page=2", http.StatusOK},
		{"search load more courses on search should return 200", "GET", "/courses/search?hitsPerPage=10", http.StatusOK},
		{"courses add course to blueprint should return 200", "POST", "/courses/blueprint?course=NSWI120&year=0&semester=0", http.StatusOK},

		// Errors
		{"courses non-existent page should return 404", "GET", "/courses/lorem", http.StatusNotFound},
		{"courses non-existing page should return 404", "GET", "/courses/NSWI120", http.StatusNotFound},
		{"courses non-number page should return 400", "GET", "/courses/?page=asdf", http.StatusBadRequest},
		{"courses negative page should return 400", "GET", "/courses/?page=-1", http.StatusBadRequest},
		{"courses zero page should return 400", "GET", "/courses/?page=0", http.StatusBadRequest},
		{"courses non-number hitsPerPage should return 400", "GET", "/courses/?hitsPerPage=asdf", http.StatusBadRequest},
		{"courses negative hitsPerPage should return 400", "GET", "/courses/?hitsPerPage=-1", http.StatusBadRequest},
		{"courses zero hitsPerPage should return 400", "GET", "/courses/?hitsPerPage=0", http.StatusBadRequest},
		{"search non-existent page should return 404", "GET", "/courses/search/lorem", http.StatusNotFound},
		{"search non-existing page should return 404", "GET", "/courses/search/NSWI120", http.StatusNotFound},
		{"search non-number page should return 400", "GET", "/courses/search?page=asdf", http.StatusBadRequest},
		{"search negative page should return 400", "GET", "/courses/search?page=-1", http.StatusBadRequest},
		{"search zero page should return 400", "GET", "/courses/search?page=0", http.StatusBadRequest},
		{"search non-number hitsPerPage should return 400", "GET", "/courses/search?hitsPerPage=asdf", http.StatusBadRequest},
		{"search negative hitsPerPage should return 400", "GET", "/courses/search?hitsPerPage=-1", http.StatusBadRequest},
		{"search zero hitsPerPage should return 400", "GET", "/courses/search?hitsPerPage=0", http.StatusBadRequest},
		{"courses add same course to blueprint should return 409", "POST", "/courses/blueprint?course=NSWI120&year=0&semester=0", http.StatusConflict},
		{"courses add course to blueprint without course should return 400", "POST", "/courses/blueprint?year=0&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint without year should return 400", "POST", "/courses/blueprint?course=NSWI120&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint without semester should return 400", "POST", "/courses/blueprint?course=NSWI120&year=0", http.StatusBadRequest},
		{"courses add course to blueprint with invalid course should return 400", "POST", "/courses/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with invalid year should return 400", "POST", "/courses/blueprint?course=NSWI120&year=asdf&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with negative year should return 400", "POST", "/courses/blueprint?course=NSWI120&year=-1&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with invalid semester should return 400", "POST", "/courses/blueprint?course=NSWI120&year=0&semester=asdf", http.StatusBadRequest},
		{"courses add course to blueprint with negative semester should return 400", "POST", "/courses/blueprint?course=NSWI120&year=0&semester=-1", http.StatusBadRequest},
		{"courses add course to blueprint with invalid course code should return 400", "POST", "/courses/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestCourseDetailServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"course detail page for NSWI120 should return 200", "GET", "/course/NSWI120", http.StatusOK},
		{"course detail page for NPRG024 should return 200", "GET", "/course/NPRG024", http.StatusOK},
		{"course detail page for NSWI166 should return 200", "GET", "/course/NSWI166", http.StatusOK},
		{"course detail page for NSWI120 with cs language should return 200", "GET", "/cs/course/NSWI120", http.StatusOK},
		{"course detail page for NSWI120 with en language should return 200", "GET", "/en/course/NSWI120", http.StatusOK},
		{"survey for NSWI120 should return 200", "GET", "/course/survey/NSWI120", http.StatusOK},
		{"next survey for NSWI120 should return 200", "GET", "/course/survey/next/NSWI120", http.StatusOK},
		{"rating category 1 for NSWI120 should return 200", "PUT", "/course/rating/NSWI120/1?rating=0", http.StatusOK},
		{"rating category 2 for NSWI120 should return 200", "PUT", "/course/rating/NSWI120/2?rating=0", http.StatusOK},
		{"rating category 3 for NSWI120 should return 200", "PUT", "/course/rating/NSWI120/3?rating=0", http.StatusOK},
		{"rating category 4 for NSWI120 should return 200", "PUT", "/course/rating/NSWI120/4?rating=0", http.StatusOK},
		{"rating category 5 for NSWI120 should return 200", "PUT", "/course/rating/NSWI120/5?rating=0", http.StatusOK},
		{"delete rating category 1 for NSWI120 should return 200", "DELETE", "/course/rating/NSWI120/1", http.StatusOK},
		{"delete rating category 2 for NSWI120 should return 200", "DELETE", "/course/rating/NSWI120/2", http.StatusOK},
		{"delete rating category 3 for NSWI120 should return 200", "DELETE", "/course/rating/NSWI120/3", http.StatusOK},
		{"delete rating category 4 for NSWI120 should return 200", "DELETE", "/course/rating/NSWI120/4", http.StatusOK},
		{"delete rating category 5 for NSWI120 should return 200", "DELETE", "/course/rating/NSWI120/5", http.StatusOK},
		{"rating NSWI120 positively should return 200", "PUT", "/course/rating/NSWI120?rating=1", http.StatusOK},
		{"repeated rating NSWI120 positively should return 200", "PUT", "/course/rating/NSWI120?rating=1", http.StatusOK},
		{"rating NSWI120 negatively should return 200", "PUT", "/course/rating/NSWI120?rating=0", http.StatusOK},
		{"repeated rating NSWI120 negatively should return 200", "PUT", "/course/rating/NSWI120?rating=0", http.StatusOK},
		{"delete rating NSWI120 should return 200", "DELETE", "/course/rating/NSWI120", http.StatusOK},
		{"repeated delete rating NSWI120 should return 200", "DELETE", "/course/rating/NSWI120", http.StatusOK},
		// TODO: add to blueprint
		// TODO: page not found

		// Errors
		{"course detail page for non-existent course should return 404", "GET", "/course/XXNOTXX", http.StatusNotFound},
		{"survey for non-existent course should return 404", "GET", "/course/XXNOTXX/survey", http.StatusNotFound},
		{"next survey for non-existent course should return 404", "GET", "/course/XXNOTXX/survey?next=true", http.StatusNotFound},
		{"rating non-existent 0- category for NSWI120 should return 400", "PUT", "/course/rating/NSWI120/0?rating=0", http.StatusBadRequest},
		{"rating non-existent 6+ category for NSWI120 should return 400", "PUT", "/course/rating/NSWI120/6?rating=0", http.StatusBadRequest},
		{"rating category for NSWI120 without rating should return 400", "PUT", "/course/rating/NSWI120/1", http.StatusBadRequest},
		{"rating category for NSWI120 with empty rating should return 400", "PUT", "/course/rating/NSWI120/1?rating=", http.StatusBadRequest},
		{"rating category for NSWI120 with invalid rating should return 400", "PUT", "/course/rating/NSWI120/1?rating=lorem", http.StatusBadRequest},
		{"rating category for NSWI120 with negative rating should return 400", "PUT", "/course/rating/NSWI120/1?rating=-1", http.StatusBadRequest},
		{"rating category for NSWI120 with big rating should return 400", "PUT", "/course/rating/NSWI120/1?rating=999", http.StatusBadRequest},
		{"delete non-existent 0- rating category for NSWI120 should return 400", "DELETE", "/course/rating/NSWI120/0", http.StatusBadRequest},
		{"delete non-existent 6+ rating category for NSWI120 should return 400", "DELETE", "/course/rating/NSWI120/6", http.StatusBadRequest},
		{"delete rating category for non-existent course should return 400", "DELETE", "/course/rating/XXNOTXX/1", http.StatusNotFound},
		{"rating NSWI120 without rating should return 400", "PUT", "/course/rating/NSWI120", http.StatusBadRequest},
		{"rating NSWI120 with empty rating should return 400", "PUT", "/course/rating/NSWI120?rating=", http.StatusBadRequest},
		{"rating NSWI120 with invalid rating should return 400", "PUT", "/course/rating/NSWI120?rating=lorem", http.StatusBadRequest},
		{"rating NSWI120 with negative rating should return 400", "PUT", "/course/rating/NSWI120?rating=-1", http.StatusBadRequest},
		{"rating NSWI120 with big rating should return 400", "PUT", "/course/rating/NSWI120?rating=2", http.StatusBadRequest},
		{"delete rating for non-existent course should return 400", "DELETE", "/course/rating/XXNOTXX", http.StatusNotFound},
	}

	runTests(t, tests)
}
func TestBlueprintServer(t *testing.T) {
	tests := []testCase{}

	runTests(t, tests)
}

func TestDegreePlanServer(t *testing.T) {
	tests := []testCase{}
	runTests(t, tests)
}

// TODO: add page server tests -> quicksearch

//================================================================================
// Test Utilities
//================================================================================

type testCase struct {
	name   string
	method string
	url    string
	want   int
}

func runTests(t *testing.T, tests []testCase) {
	ts := setupTestServer(t)
	defer ts.Close()

	client := ts.Client()
	sessionCookie := setupTestUser(ts, t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.url, nil)
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

func setupTestServer(t *testing.T) *httptest.Server {
	conf := configFrom("./config.dev.toml")
	handler := setupHandler(conf)
	removeTestUserFromDB(t, conf)
	return httptest.NewTLSServer(handler)
}

func removeTestUserFromDB(t *testing.T, conf config) {
	db := setupDB(conf)
	defer db.Close()
	_, err := db.Exec("DELETE FROM users WHERE id = $1", testUserID)
	if err != nil {
		t.Logf("Warning: Failed to remove test user from DB: %v", err)
	}
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
	form.Set("user-id", testUserID)
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
