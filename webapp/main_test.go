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
		{"root path should return 200",
			"GET", "/", http.StatusOK},
		{"cs root path should return 200",
			"GET", "/cs/", http.StatusOK},
		{"en root path should return 200",
			"GET", "/en/", http.StatusOK},
		{"home path should return 200",
			"GET", "/home/", http.StatusOK},
		{"cs home path should return 200",
			"GET", "/cs/home/", http.StatusOK},
		{"en home path should return 200",
			"GET", "/en/home/", http.StatusOK},

		// Errors
		{"root non-existent page should return 404",
			"GET", "/homer/", http.StatusNotFound},
		{"root non-existing page should return 404",
			"GET", "/lorem", http.StatusNotFound},
	}
	runTests(t, tests)
}

func TestCoursesServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"courses root path should return 200",
			"GET", "/courses/", http.StatusOK},
		{"courses cs root path should return 200",
			"GET", "/cs/courses/", http.StatusOK},
		{"courses en root path should return 200",
			"GET", "/en/courses/", http.StatusOK},
		{"courses search query should return 200",
			"GET", "/courses/?search=peska", http.StatusOK},
		{"courses search query with random string should return 200",
			"GET", "/courses/?search=asdfghjkl", http.StatusOK},
		{"courses search query with some course should return 200",
			"GET", "/courses/?search=NSWI120", http.StatusOK},
		{"courses search query with empty string should return 200",
			"GET", "/courses/?search=", http.StatusOK},
		{"courses pagination should return 200",
			"GET", "/courses/?page=2", http.StatusOK},
		{"courses load more courses should return 200",
			"GET", "/courses/?hitsPerPage=10", http.StatusOK},
		{"search courses should return 200",
			"GET", "/courses/search", http.StatusOK},
		{"search search query should return 200",
			"GET", "/courses/search?search=peska", http.StatusOK},
		{"search search query on search with random string should return 200",
			"GET", "/courses/search?search=asdfghjkl", http.StatusOK},
		{"search search query on search with some course should return 200",
			"GET", "/courses/search?search=NSWI120", http.StatusOK},
		{"search search query on search with empty string should return 200",
			"GET", "/courses/search?search=", http.StatusOK},
		{"search pagination on search should return 200",
			"GET", "/courses/search?page=2", http.StatusOK},
		{"search load more courses on search should return 200",
			"GET", "/courses/search?hitsPerPage=10", http.StatusOK},
		{"courses add course to blueprint should return 200",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=0", http.StatusOK},

		// Errors
		{"courses non-existent page should return 404",
			"GET", "/courses/lorem", http.StatusNotFound},
		{"courses non-existing page should return 404",
			"GET", "/courses/NSWI120", http.StatusNotFound},
		{"courses non-number page should return 400",
			"GET", "/courses/?page=asdf", http.StatusBadRequest},
		{"courses negative page should return 400",
			"GET", "/courses/?page=-1", http.StatusBadRequest},
		{"courses zero page should return 400",
			"GET", "/courses/?page=0", http.StatusBadRequest},
		{"courses non-number hitsPerPage should return 400",
			"GET", "/courses/?hitsPerPage=asdf", http.StatusBadRequest},
		{"courses negative hitsPerPage should return 400",
			"GET", "/courses/?hitsPerPage=-1", http.StatusBadRequest},
		{"courses zero hitsPerPage should return 400",
			"GET", "/courses/?hitsPerPage=0", http.StatusBadRequest},
		{"search non-existent page should return 404",
			"GET", "/courses/search/lorem", http.StatusNotFound},
		{"search non-existing page should return 404",
			"GET", "/courses/search/NSWI120", http.StatusNotFound},
		{"search non-number page should return 400",
			"GET", "/courses/search?page=asdf", http.StatusBadRequest},
		{"search negative page should return 400",
			"GET", "/courses/search?page=-1", http.StatusBadRequest},
		{"search zero page should return 400",
			"GET", "/courses/search?page=0", http.StatusBadRequest},
		{"search non-number hitsPerPage should return 400",
			"GET", "/courses/search?hitsPerPage=asdf", http.StatusBadRequest},
		{"search negative hitsPerPage should return 400",
			"GET", "/courses/search?hitsPerPage=-1", http.StatusBadRequest},
		{"search zero hitsPerPage should return 400",
			"GET", "/courses/search?hitsPerPage=0", http.StatusBadRequest},
		{"courses add same course to blueprint should return 409",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=0", http.StatusConflict},
		{"courses add course to blueprint without course should return 400",
			"POST", "/courses/blueprint?year=0&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with empty course should return 400",
			"POST", "/courses/blueprint?course=&year=0&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with non-existent course should return 400",
			"POST", "/courses/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint without year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint without semester should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=0", http.StatusBadRequest},
		{"courses add course to blueprint with invalid course should return 400",
			"POST", "/courses/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with invalid year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=asdf&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with negative year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=-1&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with non-existent year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=1&semester=0", http.StatusBadRequest},
		{"courses add course to blueprint with invalid semester should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=asdf", http.StatusBadRequest},
		{"courses add course to blueprint with negative semester should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=-1", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestCourseDetailServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"course detail page for NSWI120 should return 200",
			"GET", "/course/NSWI120", http.StatusOK},
		{"course detail page for NPRG024 should return 200",
			"GET", "/course/NPRG024", http.StatusOK},
		{"course detail page for NSWI166 should return 200",
			"GET", "/course/NSWI166", http.StatusOK},
		{"course detail page for NSWI120 with cs language should return 200",
			"GET", "/cs/course/NSWI120", http.StatusOK},
		{"course detail page for NSWI120 with en language should return 200",
			"GET", "/en/course/NSWI120", http.StatusOK},
		{"survey for NSWI120 should return 200",
			"GET", "/course/survey/NSWI120", http.StatusOK},
		{"next survey for NSWI120 should return 200",
			"GET", "/course/survey/next/NSWI120", http.StatusOK},
		{"rating category 1 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/1?rating=0", http.StatusOK},
		{"rating category 2 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/2?rating=0", http.StatusOK},
		{"rating category 3 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/3?rating=0", http.StatusOK},
		{"rating category 4 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/4?rating=0", http.StatusOK},
		{"rating category 5 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/5?rating=0", http.StatusOK},
		{"delete rating category 1 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/1", http.StatusOK},
		{"delete rating category 2 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/2", http.StatusOK},
		{"delete rating category 3 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/3", http.StatusOK},
		{"delete rating category 4 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/4", http.StatusOK},
		{"delete rating category 5 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/5", http.StatusOK},
		{"rating NSWI120 positively should return 200",
			"PUT", "/course/rating/NSWI120?rating=1", http.StatusOK},
		{"repeated rating NSWI120 positively should return 200",
			"PUT", "/course/rating/NSWI120?rating=1", http.StatusOK},
		{"rating NSWI120 negatively should return 200",
			"PUT", "/course/rating/NSWI120?rating=0", http.StatusOK},
		{"repeated rating NSWI120 negatively should return 200",
			"PUT", "/course/rating/NSWI120?rating=0", http.StatusOK},
		{"delete rating NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120", http.StatusOK},
		{"add course to blueprint should return 200",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=0", http.StatusOK},

		// Errors
		{"course detail page for non-existent course should return 404",
			"GET", "/course/XXNOTXX", http.StatusNotFound},
		{"survey for non-existent course should return 404",
			"GET", "/course/XXNOTXX/survey", http.StatusNotFound},
		{"next survey for non-existent course should return 404",
			"GET", "/course/XXNOTXX/survey?next=true", http.StatusNotFound},
		{"rating non-existent 0- category for NSWI120 should return 400",
			"PUT", "/course/rating/NSWI120/0?rating=0", http.StatusBadRequest},
		{"rating non-existent 6+ category for NSWI120 should return 400",
			"PUT", "/course/rating/NSWI120/6?rating=0", http.StatusBadRequest},
		{"rating category for NSWI120 without rating should return 400",
			"PUT", "/course/rating/NSWI120/1", http.StatusBadRequest},
		{"rating category for NSWI120 with empty rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=", http.StatusBadRequest},
		{"rating category for NSWI120 with invalid rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=lorem", http.StatusBadRequest},
		{"rating category for NSWI120 with negative rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=-1", http.StatusBadRequest},
		{"rating category for NSWI120 with big rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=999", http.StatusBadRequest},
		{"delete non-existent 0- rating category for NSWI120 should return 400",
			"DELETE", "/course/rating/NSWI120/0", http.StatusBadRequest},
		{"delete non-existent 6+ rating category for NSWI120 should return 400",
			"DELETE", "/course/rating/NSWI120/6", http.StatusBadRequest},
		{"delete rating category for non-existent course should return 400",
			"DELETE", "/course/rating/XXNOTXX/1", http.StatusBadRequest},
		{"rating NSWI120 without rating should return 400",
			"PUT", "/course/rating/NSWI120", http.StatusBadRequest},
		{"rating NSWI120 with empty rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=", http.StatusBadRequest},
		{"rating NSWI120 with invalid rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=lorem", http.StatusBadRequest},
		{"rating NSWI120 with negative rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=-1", http.StatusBadRequest},
		{"rating NSWI120 with big rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=2", http.StatusBadRequest},
		{"repeated delete rating NSWI120 should return 400",
			"DELETE", "/course/rating/NSWI120", http.StatusBadRequest},
		{"delete rating for non-existent course should return 400",
			"DELETE", "/course/rating/XXNOTXX", http.StatusBadRequest},
		{"add same course to blueprint should return 409",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=0", http.StatusConflict},
		{"add course to blueprint without course should return 400",
			"POST", "/course/blueprint?year=0&semester=0", http.StatusBadRequest},
		{"add course to blueprint with empty course should return 400",
			"POST", "/course/blueprint?course=&year=0&semester=0", http.StatusBadRequest},
		{"add course to blueprint with non-existent course should return 400",
			"POST", "/course/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		{"add course to blueprint without year should return 400",
			"POST", "/course/blueprint?course=NSWI120&semester=0", http.StatusBadRequest},
		{"add course to blueprint without semester should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=0", http.StatusBadRequest},
		{"add course to blueprint with invalid course should return 400",
			"POST", "/course/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		{"add course to blueprint with invalid year should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=asdf&semester=0", http.StatusBadRequest},
		{"add course to blueprint with negative year should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=-1&semester=0", http.StatusBadRequest},
		{"add course to blueprint with non-existent year should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=1&semester=0", http.StatusBadRequest},
		{"add course to blueprint with invalid semester should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=asdf", http.StatusBadRequest},
		{"add course to blueprint with negative semester should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=-1", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestBlueprintServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"blueprint page should return 200",
			"GET", "/blueprint/", http.StatusOK},
		{"blueprint page with cs language should return 200",
			"GET", "/cs/blueprint/", http.StatusOK},
		{"blueprint page with en language should return 200",
			"GET", "/en/blueprint/", http.StatusOK},
		{"blueprint add year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		{"blueprint add second year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		{"blueprint add third year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		{"blueprint add fourth year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		{"blueprint remove last year should return 200",
			"DELETE", "/blueprint/year?unassign=true", http.StatusOK},
		// add some courses to blueprint for move and remove tests
		{"blueprint add NSWI152 to unassigned from courses should return 200",
			"POST", "/courses/blueprint?course=NSWI152&year=0&semester=0", http.StatusOK},
		{"blueprint add NSWI202 to unassigned from course detail should return 200",
			"POST", "/course/blueprint?course=NSWI202&year=0&semester=0", http.StatusOK},
		{"blueprint add NPRG045 to unassigned from degreeplan should return 200",
			"POST", "/degreeplan/blueprint?course=NPRG045&year=0&semester=0", http.StatusOK},
		{"blueprint add NSWI152 to first year from courses should return 200",
			"POST", "/courses/blueprint?course=NSWI152&year=1&semester=1", http.StatusOK},
		{"blueprint add NPRG024 to first year from courses should return 200",
			"POST", "/courses/blueprint?course=NPRG024&year=1&semester=1", http.StatusOK},
		{"blueprint add NSWI035 to second year from course detail should return 200",
			"POST", "/course/blueprint?course=NSWI035&year=2&semester=1", http.StatusOK},
		{"blueprint add NPRG069 to third year from degreeplan should return 200",
			"POST", "/degreeplan/blueprint?course=NPRG069&year=3&semester=1", http.StatusOK},
		{"blueprint add NSWI153 to first year second semester from courses should return 200",
			"POST", "/courses/blueprint?course=NSWI153&year=1&semester=2", http.StatusOK},
		{"blueprint add NSWI142 to second year second semester from course detail should return 200",
			"POST", "/course/blueprint?course=NSWI142&year=2&semester=2", http.StatusOK},
		{"blueprint add NDBI046 to third year second semester from degreeplan should return 200",
			"POST", "/degreeplan/blueprint?course=NDBI046&year=3&semester=2", http.StatusOK},
		// =========================================
		{"blueprint unassign second year first semester should return 200",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=2&semester=1", http.StatusOK},
		{"blueprint delete second year second semester should return 200",
			"DELETE", "/blueprint/courses?type=semester-remove&year=2&semester=2", http.StatusOK},
		{"blueprint folding unassigned should return 200",
			"PATCH", "/blueprint/fold?year=0&semester=0&folded=true", http.StatusOK},
		{"blueprint folding first year second semester should return 200",
			"PATCH", "/blueprint/fold?year=1&semester=2&folded=false", http.StatusOK},
		{"blueprint folding second year first semester should return 200",
			"PATCH", "/blueprint/fold?year=2&semester=1&folded=true", http.StatusOK},

		// Errors
		{"blueprint page for non-existent page should return 404",
			"GET", "/blueprint/lorem", http.StatusNotFound},
		{"blueprint page for non-existing page should return 404",
			"GET", "/blueprint/ipsum", http.StatusNotFound},
		{"blueprint remove year without unassign should return 400",
			"DELETE", "/blueprint/year", http.StatusBadRequest},
		{"blueprint remove year with empty unassign should return 400",
			"DELETE", "/blueprint/year?unassign=", http.StatusBadRequest},
		{"blueprint remove year with invalid unassign should return 400",
			"DELETE", "/blueprint/year?unassign=lorem", http.StatusBadRequest},
		{"blueprint fold without year should return 400",
			"PATCH", "/blueprint/fold?semester=1&folded=true", http.StatusBadRequest},
		{"blueprint fold with empty year should return 400",
			"PATCH", "/blueprint/fold?year=&semester=1&folded=true", http.StatusBadRequest},
		{"blueprint fold with invalid year should return 400",
			"PATCH", "/blueprint/fold?year=lorem&semester=1&folded=true", http.StatusBadRequest},
		{"blueprint fold with negative year should return 400",
			"PATCH", "/blueprint/fold?year=-1&semester=1&folded=true", http.StatusBadRequest},
		{"blueprint fold without semester should return 400",
			"PATCH", "/blueprint/fold?year=0&folded=true", http.StatusBadRequest},
		{"blueprint fold with empty semester should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=&folded=true", http.StatusBadRequest},
		{"blueprint fold with invalid semester should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=lorem&folded=true", http.StatusBadRequest},
		{"blueprint fold with negative semester should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=-1&folded=true", http.StatusBadRequest},
		{"blueprint fold without folded should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=1", http.StatusBadRequest},
		{"blueprint fold with empty folded should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=1&folded=", http.StatusBadRequest},
		{"blueprint fold with invalid folded should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=1&folded=lorem", http.StatusBadRequest},
		{"blueprint move course with bad ID should return 400",
			"PATCH", "/blueprint/course/987654321?year=1&semester=1&position=-1", http.StatusBadRequest},
		{"blueprint move courses without type should return 400",
			"PATCH", "/blueprint/courses", http.StatusBadRequest},
		{"blueprint move courses with empty type should return 400",
			"PATCH", "/blueprint/courses?type=", http.StatusBadRequest},
		{"blueprint move courses with invalid type should return 400",
			"PATCH", "/blueprint/courses?type=lorem", http.StatusBadRequest},
		{"blueprint move courses without year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&semester=1", http.StatusBadRequest},
		{"blueprint move courses with empty year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=&semester=1", http.StatusBadRequest},
		{"blueprint move courses with invalid year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=lorem&semester=1", http.StatusBadRequest},
		{"blueprint move courses with negative year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=-1&semester=1", http.StatusBadRequest},
		{"blueprint move courses with non-existent year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=5&semester=1", http.StatusBadRequest},
		{"blueprint move courses without semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1", http.StatusBadRequest},
		{"blueprint move courses with empty semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=", http.StatusBadRequest},
		{"blueprint move courses with invalid semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=lorem", http.StatusBadRequest},
		{"blueprint move courses with negative semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=-1", http.StatusBadRequest},
		{"blueprint move courses with non-existent semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=5", http.StatusBadRequest},
		{"blueprint move courses without position should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1", http.StatusBadRequest},
		{"blueprint move courses with empty position should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1&position=", http.StatusBadRequest},
		{"blueprint move courses with invalid position should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1&position=lorem", http.StatusBadRequest},
		{"blueprint move courses without selected should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1&position=-1", http.StatusBadRequest},
		{"blueprint unassign empty second year first semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=2&semester=1", http.StatusBadRequest},
		{"blueprint delete empty second year second semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=2&semester=2", http.StatusBadRequest},
		{"blueprint remove course with bad ID should return 400",
			"DELETE", "/blueprint/course/987654321?year=1&semester=1&position=-1", http.StatusBadRequest},
		{"blueprint remove courses without type should return 400",
			"DELETE", "/blueprint/courses", http.StatusBadRequest},
		{"blueprint remove courses with empty type should return 400",
			"DELETE", "/blueprint/courses?type=", http.StatusBadRequest},
		{"blueprint remove courses with invalid type should return 400",
			"DELETE", "/blueprint/courses?type=lorem", http.StatusBadRequest},
		{"blueprint remove courses without year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&semester=1", http.StatusBadRequest},
		{"blueprint remove courses with empty year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=&semester=1", http.StatusBadRequest},
		{"blueprint remove courses with invalid year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=lorem&semester=1", http.StatusBadRequest},
		{"blueprint remove courses with negative year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=-1&semester=1", http.StatusBadRequest},
		{"blueprint remove courses with non-existent year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=5&semester=1", http.StatusBadRequest},
		{"blueprint remove courses without semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1", http.StatusBadRequest},
		{"blueprint remove courses with empty semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=", http.StatusBadRequest},
		{"blueprint remove courses with invalid semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=lorem", http.StatusBadRequest},
		{"blueprint remove courses with negative semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=-1", http.StatusBadRequest},
		{"blueprint remove courses with non-existent semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=5", http.StatusBadRequest},
		{"blueprint remove courses without selected should return 400",
			"DELETE", "/blueprint/courses?type=selected-remove", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestDegreePlanServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"degree plan page should return 200",
			"GET", "/degreeplan/", http.StatusOK},
		{"degree plan page with cs language should return 200",
			"GET", "/cs/degreeplan/", http.StatusOK},
		{"degree plan page with en language should return 200",
			"GET", "/en/degreeplan/", http.StatusOK},
		{"degree plan existing code NIDAW19B and year 2025 should return 200",
			"GET", "/degreeplan/NIDAW19B?search-dp-year=2025", http.StatusOK},
		{"degree plan existing code NISD23N and year 2023 should return 200",
			"GET", "/degreeplan/NISD23N?search-dp-year=2023", http.StatusOK},
		{"degree plan existing code NIPP19B and year 2021 should return 200",
			"GET", "/degreeplan/NIPP19B?search-dp-year=2021", http.StatusOK},
		{"degree plan save NIPP19B and year 2025 should return 200",
			"PATCH", "/degreeplan/NIPP19B?save-dp-year=2025", http.StatusOK},
		{"degree plan save NISD23N and year 2023 should return 200",
			"PATCH", "/degreeplan/NISD23N?save-dp-year=2023", http.StatusOK},
		{"degree plan save NIDAW19B and year 2021 should return 200",
			"PATCH", "/degreeplan/NIDAW19B?save-dp-year=2021", http.StatusOK},
		{"degree plan search by code NIDAW19B should return 200",
			"GET", "/degreeplan/search?search-dp-query=NIDAW19B", http.StatusOK},
		{"degree plan search by 'softwarove' should return 200",
			"GET", "/degreeplan/search?search-dp-query=softwarove", http.StatusOK},
		{"degree plan search by 'a' should return 200",
			"GET", "/degreeplan/search?search-dp-query=a", http.StatusOK},
		{"degree plan search by empty query should return 200",
			"GET", "/degreeplan/search?search-dp-query=", http.StatusOK},
		{"degree plan add course to blueprint should return 200",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=0", http.StatusOK},
		{"degree plan add courses to blueprint should return 200",
			"POST", "/degreeplan/blueprint?selected-courses=NSWI035&selected-courses=NPRG024&year=0&semester=0", http.StatusOK},

		// Errors
		{"degree plan page for non-existent plan should return 404",
			"GET", "/degreeplan/lorem?search-dp-year=2025", http.StatusNotFound},
		{"degree plan specific without year should return 400",
			"GET", "/degreeplan/NISD23N", http.StatusBadRequest},
		{"degree plan specific with empty year should return 400",
			"GET", "/degreeplan/NISD23N?search-dp-year=", http.StatusBadRequest},
		{"degree plan specific with invalid year should return 400",
			"GET", "/degreeplan/NISD23N?search-dp-year=lorem", http.StatusBadRequest},
		{"degree plan specific with negative year should return 404",
			"GET", "/degreeplan/NISD23N?search-dp-year=-1", http.StatusNotFound},
		{"degree plan specific with future year should return 404",
			"GET", "/degreeplan/NISD23N?search-dp-year=2050", http.StatusNotFound},
		{"degree plan specific with old year should return 404",
			"GET", "/degreeplan/NISD23N?search-dp-year=1999", http.StatusNotFound},
		{"degree plan save with invalid code should return 400",
			"PATCH", "/degreeplan/lorem?save-dp-year=2025", http.StatusBadRequest},
		{"degree plan save without year should return 400",
			"PATCH", "/degreeplan/NISD23N", http.StatusBadRequest},
		{"degree plan save with empty year should return 400",
			"PATCH", "/degreeplan/NISD23N?save-dp-year=", http.StatusBadRequest},
		{"degree plan save with invalid year should return 400",
			"PATCH", "/degreeplan/NISD23N?save-dp-year=lorem", http.StatusBadRequest},
		{"degree plan save with negative year should return 404",
			"PATCH", "/degreeplan/NISD23N?save-dp-year=-1", http.StatusNotFound},
		{"degree plan save with future year should return 404",
			"PATCH", "/degreeplan/NISD23N?save-dp-year=2050", http.StatusNotFound},
		{"degree plan save with old year should return 404",
			"PATCH", "/degreeplan/NISD23N?save-dp-year=1999", http.StatusNotFound},
		{"degree plan add same course to blueprint should return 409",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=0", http.StatusConflict},
		{"degree plan add course to blueprint without course should return 400",
			"POST", "/degreeplan/blueprint?year=0&semester=0", http.StatusBadRequest},
		{"degree plan add course to blueprint without year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&semester=0", http.StatusBadRequest},
		{"degree plan add course to blueprint without semester should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0", http.StatusBadRequest},
		{"degree plan add course to blueprint with invalid course should return 400",
			"POST", "/degreeplan/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		{"degree plan add course to blueprint with invalid year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=asdf&semester=0", http.StatusBadRequest},
		{"degree plan add course to blueprint with negative year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=-1&semester=0", http.StatusBadRequest},
		{"degree plan add course to blueprint with non-existent year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=1&semester=0", http.StatusBadRequest},
		{"degree plan add course to blueprint with invalid semester should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=asdf", http.StatusBadRequest},
		{"degree plan add course to blueprint with negative semester should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=-1", http.StatusBadRequest},
	}
	runTests(t, tests)
}

func TestPageServer(t *testing.T) {
	tests := []testCase{
		// Happy path
		{"quicksearching 'NSWI153' should return 200",
			"GET", "/page/quicksearch?q=NSWI153", http.StatusOK},
		{"quicksearching 'NSWI153' with cs should return 200",
			"GET", "/cs/page/quicksearch?q=NSWI153", http.StatusOK},
		{"quicksearching 'NSWI153' with en should return 200",
			"GET", "/en/page/quicksearch?q=NSWI153", http.StatusOK},
		{"quicksearching 'NPRG024' should return 200",
			"GET", "/page/quicksearch?q=NPRG024", http.StatusOK},
		{"quicksearching 'peska' should return 200",
			"GET", "/page/quicksearch?q=peska", http.StatusOK},
		{"quicksearching random text should return 200",
			"GET", "/page/quicksearch?q=asagafe45", http.StatusOK},
		{"quicksearching empty query should return 200",
			"GET", "/page/quicksearch?q=", http.StatusOK},
		{"quicksearching with missing query equals empty query should return 200",
			"GET", "/page/quicksearch", http.StatusOK},
	}

	runTests(t, tests)
}

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
