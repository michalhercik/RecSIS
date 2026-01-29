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
	tests := []testRunner{
		// Happy path
		testCase{"root path should return 200",
			"GET", "/", http.StatusOK},
		testCase{"cs root path should return 200",
			"GET", "/cs/", http.StatusOK},
		testCase{"en root path should return 200",
			"GET", "/en/", http.StatusOK},
		testCase{"home path should return 200",
			"GET", "/home/", http.StatusOK},
		testCase{"cs home path should return 200",
			"GET", "/cs/home/", http.StatusOK},
		testCase{"en home path should return 200",
			"GET", "/en/home/", http.StatusOK},

		// Errors
		testCase{"root non-existent page should return 404",
			"GET", "/homer/", http.StatusNotFound},
		testCase{"root non-existing page should return 404",
			"GET", "/lorem", http.StatusNotFound},
	}
	runTests(t, tests)
}

func TestCoursesServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"courses root path should return 200",
			"GET", "/courses/", http.StatusOK},
		testCase{"courses cs root path should return 200",
			"GET", "/cs/courses/", http.StatusOK},
		testCase{"courses en root path should return 200",
			"GET", "/en/courses/", http.StatusOK},
		testCase{"courses search query should return 200",
			"GET", "/courses/?search=peska", http.StatusOK},
		testCase{"courses search query with random string should return 200",
			"GET", "/courses/?search=asdfghjkl", http.StatusOK},
		testCase{"courses search query with some course should return 200",
			"GET", "/courses/?search=NSWI120", http.StatusOK},
		testCase{"courses search query with empty string should return 200",
			"GET", "/courses/?search=", http.StatusOK},
		testCase{"courses pagination should return 200",
			"GET", "/courses/?page=2", http.StatusOK},
		testCase{"courses load more courses should return 200",
			"GET", "/courses/?hitsPerPage=10", http.StatusOK},
		testCase{"search courses should return 200",
			"GET", "/courses/search", http.StatusOK},
		testCase{"search search query should return 200",
			"GET", "/courses/search?search=peska", http.StatusOK},
		testCase{"search search query on search with random string should return 200",
			"GET", "/courses/search?search=asdfghjkl", http.StatusOK},
		testCase{"search search query on search with some course should return 200",
			"GET", "/courses/search?search=NSWI120", http.StatusOK},
		testCase{"search search query on search with empty string should return 200",
			"GET", "/courses/search?search=", http.StatusOK},
		testCase{"search pagination on search should return 200",
			"GET", "/courses/search?page=2", http.StatusOK},
		testCase{"search load more courses on search should return 200",
			"GET", "/courses/search?hitsPerPage=10", http.StatusOK},
		testCase{"courses add course to blueprint should return 200",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=0", http.StatusOK},

		// Errors
		testCase{"courses non-existent page should return 404",
			"GET", "/courses/lorem", http.StatusNotFound},
		testCase{"courses non-existing page should return 404",
			"GET", "/courses/NSWI120", http.StatusNotFound},
		testCase{"courses non-number page should return 400",
			"GET", "/courses/?page=asdf", http.StatusBadRequest},
		testCase{"courses negative page should return 400",
			"GET", "/courses/?page=-1", http.StatusBadRequest},
		testCase{"courses zero page should return 400",
			"GET", "/courses/?page=0", http.StatusBadRequest},
		testCase{"courses non-number hitsPerPage should return 400",
			"GET", "/courses/?hitsPerPage=asdf", http.StatusBadRequest},
		testCase{"courses negative hitsPerPage should return 400",
			"GET", "/courses/?hitsPerPage=-1", http.StatusBadRequest},
		testCase{"courses zero hitsPerPage should return 400",
			"GET", "/courses/?hitsPerPage=0", http.StatusBadRequest},
		testCase{"search non-existent page should return 404",
			"GET", "/courses/search/lorem", http.StatusNotFound},
		testCase{"search non-existing page should return 404",
			"GET", "/courses/search/NSWI120", http.StatusNotFound},
		testCase{"search non-number page should return 400",
			"GET", "/courses/search?page=asdf", http.StatusBadRequest},
		testCase{"search negative page should return 400",
			"GET", "/courses/search?page=-1", http.StatusBadRequest},
		testCase{"search zero page should return 400",
			"GET", "/courses/search?page=0", http.StatusBadRequest},
		testCase{"search non-number hitsPerPage should return 400",
			"GET", "/courses/search?hitsPerPage=asdf", http.StatusBadRequest},
		testCase{"search negative hitsPerPage should return 400",
			"GET", "/courses/search?hitsPerPage=-1", http.StatusBadRequest},
		testCase{"search zero hitsPerPage should return 400",
			"GET", "/courses/search?hitsPerPage=0", http.StatusBadRequest},
		testCase{"courses add same course to blueprint should return 409",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=0", http.StatusConflict},
		testCase{"courses add course to blueprint without course should return 400",
			"POST", "/courses/blueprint?year=0&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with empty course should return 400",
			"POST", "/courses/blueprint?course=&year=0&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with non-existent course should return 400",
			"POST", "/courses/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint without year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint without semester should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with invalid course should return 400",
			"POST", "/courses/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with invalid year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=asdf&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with negative year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=-1&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with non-existent year should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=1&semester=0", http.StatusBadRequest},
		testCase{"courses add course to blueprint with invalid semester should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=asdf", http.StatusBadRequest},
		testCase{"courses add course to blueprint with negative semester should return 400",
			"POST", "/courses/blueprint?course=NSWI120&year=0&semester=-1", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestCourseDetailServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"course detail page for NSWI120 should return 200",
			"GET", "/course/NSWI120", http.StatusOK},
		testCase{"course detail page for NPRG024 should return 200",
			"GET", "/course/NPRG024", http.StatusOK},
		testCase{"course detail page for NSWI166 should return 200",
			"GET", "/course/NSWI166", http.StatusOK},
		testCase{"course detail page for NSWI120 with cs language should return 200",
			"GET", "/cs/course/NSWI120", http.StatusOK},
		testCase{"course detail page for NSWI120 with en language should return 200",
			"GET", "/en/course/NSWI120", http.StatusOK},
		testCase{"survey for NSWI120 should return 200",
			"GET", "/course/survey/NSWI120", http.StatusOK},
		testCase{"next survey for NSWI120 should return 200",
			"GET", "/course/survey/next/NSWI120", http.StatusOK},
		testCase{"rating category 1 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/1?rating=0", http.StatusOK},
		testCase{"rating category 2 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/2?rating=0", http.StatusOK},
		testCase{"rating category 3 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/3?rating=0", http.StatusOK},
		testCase{"rating category 4 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/4?rating=0", http.StatusOK},
		testCase{"rating category 5 for NSWI120 should return 200",
			"PUT", "/course/rating/NSWI120/5?rating=0", http.StatusOK},
		testCase{"delete rating category 1 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/1", http.StatusOK},
		testCase{"delete rating category 2 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/2", http.StatusOK},
		testCase{"delete rating category 3 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/3", http.StatusOK},
		testCase{"delete rating category 4 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/4", http.StatusOK},
		testCase{"delete rating category 5 for NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120/5", http.StatusOK},
		testCase{"rating NSWI120 positively should return 200",
			"PUT", "/course/rating/NSWI120?rating=1", http.StatusOK},
		testCase{"repeated rating NSWI120 positively should return 200",
			"PUT", "/course/rating/NSWI120?rating=1", http.StatusOK},
		testCase{"rating NSWI120 negatively should return 200",
			"PUT", "/course/rating/NSWI120?rating=0", http.StatusOK},
		testCase{"repeated rating NSWI120 negatively should return 200",
			"PUT", "/course/rating/NSWI120?rating=0", http.StatusOK},
		testCase{"delete rating NSWI120 should return 200",
			"DELETE", "/course/rating/NSWI120", http.StatusOK},
		testCase{"add course to blueprint should return 200",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=0", http.StatusOK},

		// Errors
		testCase{"course detail page for non-existent course should return 404",
			"GET", "/course/XXNOTXX", http.StatusNotFound},
		testCase{"survey for non-existent course should return 404",
			"GET", "/course/XXNOTXX/survey", http.StatusNotFound},
		testCase{"next survey for non-existent course should return 404",
			"GET", "/course/XXNOTXX/survey?next=true", http.StatusNotFound},
		testCase{"rating non-existent 0- category for NSWI120 should return 400",
			"PUT", "/course/rating/NSWI120/0?rating=0", http.StatusBadRequest},
		testCase{"rating non-existent 6+ category for NSWI120 should return 400",
			"PUT", "/course/rating/NSWI120/6?rating=0", http.StatusBadRequest},
		testCase{"rating category for NSWI120 without rating should return 400",
			"PUT", "/course/rating/NSWI120/1", http.StatusBadRequest},
		testCase{"rating category for NSWI120 with empty rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=", http.StatusBadRequest},
		testCase{"rating category for NSWI120 with invalid rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=lorem", http.StatusBadRequest},
		testCase{"rating category for NSWI120 with negative rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=-1", http.StatusBadRequest},
		testCase{"rating category for NSWI120 with big rating should return 400",
			"PUT", "/course/rating/NSWI120/1?rating=999", http.StatusBadRequest},
		testCase{"delete non-existent 0- rating category for NSWI120 should return 400",
			"DELETE", "/course/rating/NSWI120/0", http.StatusBadRequest},
		testCase{"delete non-existent 6+ rating category for NSWI120 should return 400",
			"DELETE", "/course/rating/NSWI120/6", http.StatusBadRequest},
		testCase{"delete rating category for non-existent course should return 400",
			"DELETE", "/course/rating/XXNOTXX/1", http.StatusBadRequest},
		testCase{"rating NSWI120 without rating should return 400",
			"PUT", "/course/rating/NSWI120", http.StatusBadRequest},
		testCase{"rating NSWI120 with empty rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=", http.StatusBadRequest},
		testCase{"rating NSWI120 with invalid rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=lorem", http.StatusBadRequest},
		testCase{"rating NSWI120 with negative rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=-1", http.StatusBadRequest},
		testCase{"rating NSWI120 with big rating should return 400",
			"PUT", "/course/rating/NSWI120?rating=2", http.StatusBadRequest},
		testCase{"repeated delete rating NSWI120 should return 400",
			"DELETE", "/course/rating/NSWI120", http.StatusBadRequest},
		testCase{"delete rating for non-existent course should return 400",
			"DELETE", "/course/rating/XXNOTXX", http.StatusBadRequest},
		testCase{"add same course to blueprint should return 409",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=0", http.StatusConflict},
		testCase{"add course to blueprint without course should return 400",
			"POST", "/course/blueprint?year=0&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint with empty course should return 400",
			"POST", "/course/blueprint?course=&year=0&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint with non-existent course should return 400",
			"POST", "/course/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint without year should return 400",
			"POST", "/course/blueprint?course=NSWI120&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint without semester should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=0", http.StatusBadRequest},
		testCase{"add course to blueprint with invalid course should return 400",
			"POST", "/course/blueprint?course=XXNOTXX&year=0&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint with invalid year should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=asdf&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint with negative year should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=-1&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint with non-existent year should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=1&semester=0", http.StatusBadRequest},
		testCase{"add course to blueprint with invalid semester should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=asdf", http.StatusBadRequest},
		testCase{"add course to blueprint with negative semester should return 400",
			"POST", "/course/blueprint?course=NSWI120&year=0&semester=-1", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestBlueprintServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"blueprint page should return 200",
			"GET", "/blueprint/", http.StatusOK},
		testCase{"blueprint page with cs language should return 200",
			"GET", "/cs/blueprint/", http.StatusOK},
		testCase{"blueprint page with en language should return 200",
			"GET", "/en/blueprint/", http.StatusOK},
		testCase{"blueprint add year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		testCase{"blueprint add second year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		testCase{"blueprint add third year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		testCase{"blueprint add fourth year should return 200",
			"POST", "/blueprint/year", http.StatusOK},
		testCase{"blueprint remove last year should return 200",
			"DELETE", "/blueprint/year?unassign=true", http.StatusOK},
		// add some courses to blueprint for move and remove tests
		testCase{"blueprint add NSWI152 to unassigned from courses should return 200",
			"POST", "/courses/blueprint?course=NSWI152&year=0&semester=0", http.StatusOK},
		testCase{"blueprint add NSWI202 to unassigned from course detail should return 200",
			"POST", "/course/blueprint?course=NSWI202&year=0&semester=0", http.StatusOK},
		testCaseWithReferer{"blueprint add NPRG045 to unassigned from degreeplan should return 200",
			"POST", "/degreeplan/blueprint?course=NPRG045&year=0&semester=0", "/degreeplan/NIPVS19B", http.StatusOK},
		testCase{"blueprint add NSWI152 to first year from courses should return 200",
			"POST", "/courses/blueprint?course=NSWI152&year=1&semester=1", http.StatusOK},
		testCase{"blueprint add NPRG024 to first year from courses should return 200",
			"POST", "/courses/blueprint?course=NPRG024&year=1&semester=1", http.StatusOK},
		testCase{"blueprint add NSWI035 to second year from course detail should return 200",
			"POST", "/course/blueprint?course=NSWI035&year=2&semester=1", http.StatusOK},
		testCaseWithReferer{"blueprint add NPRG069 to third year from degreeplan should return 200",
			"POST", "/degreeplan/blueprint?course=NPRG069&year=3&semester=1", "/degreeplan/NIPVS19B", http.StatusOK},
		testCase{"blueprint add NSWI153 to first year second semester from courses should return 200",
			"POST", "/courses/blueprint?course=NSWI153&year=1&semester=2", http.StatusOK},
		testCase{"blueprint add NSWI142 to second year second semester from course detail should return 200",
			"POST", "/course/blueprint?course=NSWI142&year=2&semester=2", http.StatusOK},
		testCaseWithReferer{"blueprint add NDBI046 to third year second semester from degreeplan should return 200",
			"POST", "/degreeplan/blueprint?course=NDBI046&year=3&semester=2", "/degreeplan/NIPVS19B", http.StatusOK},
		// =========================================
		testCase{"blueprint unassign second year first semester should return 200",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=2&semester=1", http.StatusOK},
		testCase{"blueprint delete second year second semester should return 200",
			"DELETE", "/blueprint/courses?type=semester-remove&year=2&semester=2", http.StatusOK},
		testCase{"blueprint folding unassigned should return 200",
			"PATCH", "/blueprint/fold?year=0&semester=0&folded=true", http.StatusOK},
		testCase{"blueprint folding first year second semester should return 200",
			"PATCH", "/blueprint/fold?year=1&semester=2&folded=false", http.StatusOK},
		testCase{"blueprint folding second year first semester should return 200",
			"PATCH", "/blueprint/fold?year=2&semester=1&folded=true", http.StatusOK},

		// Errors
		testCase{"blueprint page for non-existent page should return 404",
			"GET", "/blueprint/lorem", http.StatusNotFound},
		testCase{"blueprint page for non-existing page should return 404",
			"GET", "/blueprint/ipsum", http.StatusNotFound},
		testCase{"blueprint remove year without unassign should return 400",
			"DELETE", "/blueprint/year", http.StatusBadRequest},
		testCase{"blueprint remove year with empty unassign should return 400",
			"DELETE", "/blueprint/year?unassign=", http.StatusBadRequest},
		testCase{"blueprint remove year with invalid unassign should return 400",
			"DELETE", "/blueprint/year?unassign=lorem", http.StatusBadRequest},
		testCase{"blueprint fold without year should return 400",
			"PATCH", "/blueprint/fold?semester=1&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold with empty year should return 400",
			"PATCH", "/blueprint/fold?year=&semester=1&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold with invalid year should return 400",
			"PATCH", "/blueprint/fold?year=lorem&semester=1&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold with negative year should return 400",
			"PATCH", "/blueprint/fold?year=-1&semester=1&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold without semester should return 400",
			"PATCH", "/blueprint/fold?year=0&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold with empty semester should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold with invalid semester should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=lorem&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold with negative semester should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=-1&folded=true", http.StatusBadRequest},
		testCase{"blueprint fold without folded should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=1", http.StatusBadRequest},
		testCase{"blueprint fold with empty folded should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=1&folded=", http.StatusBadRequest},
		testCase{"blueprint fold with invalid folded should return 400",
			"PATCH", "/blueprint/fold?year=1&semester=1&folded=lorem", http.StatusBadRequest},
		testCase{"blueprint move course with bad ID should return 400",
			"PATCH", "/blueprint/course/987654321?year=1&semester=1&position=-1", http.StatusBadRequest},
		testCase{"blueprint move courses without type should return 400",
			"PATCH", "/blueprint/courses", http.StatusBadRequest},
		testCase{"blueprint move courses with empty type should return 400",
			"PATCH", "/blueprint/courses?type=", http.StatusBadRequest},
		testCase{"blueprint move courses with invalid type should return 400",
			"PATCH", "/blueprint/courses?type=lorem", http.StatusBadRequest},
		testCase{"blueprint move courses without year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&semester=1", http.StatusBadRequest},
		testCase{"blueprint move courses with empty year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=&semester=1", http.StatusBadRequest},
		testCase{"blueprint move courses with invalid year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=lorem&semester=1", http.StatusBadRequest},
		testCase{"blueprint move courses with negative year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=-1&semester=1", http.StatusBadRequest},
		testCase{"blueprint move courses with non-existent year should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=5&semester=1", http.StatusBadRequest},
		testCase{"blueprint move courses without semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1", http.StatusBadRequest},
		testCase{"blueprint move courses with empty semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=", http.StatusBadRequest},
		testCase{"blueprint move courses with invalid semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=lorem", http.StatusBadRequest},
		testCase{"blueprint move courses with negative semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=-1", http.StatusBadRequest},
		testCase{"blueprint move courses with non-existent semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=1&semester=5", http.StatusBadRequest},
		testCase{"blueprint move courses without position should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1", http.StatusBadRequest},
		testCase{"blueprint move courses with empty position should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1&position=", http.StatusBadRequest},
		testCase{"blueprint move courses with invalid position should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1&position=lorem", http.StatusBadRequest},
		testCase{"blueprint move courses without selected should return 400",
			"PATCH", "/blueprint/courses?type=move-courses&year=1&semester=1&position=-1", http.StatusBadRequest},
		testCase{"blueprint unassign empty second year first semester should return 400",
			"PATCH", "/blueprint/courses?type=semester-unassign&year=2&semester=1", http.StatusBadRequest},
		testCase{"blueprint delete empty second year second semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=2&semester=2", http.StatusBadRequest},
		testCase{"blueprint remove course with bad ID should return 400",
			"DELETE", "/blueprint/course/987654321?year=1&semester=1&position=-1", http.StatusBadRequest},
		testCase{"blueprint remove courses without type should return 400",
			"DELETE", "/blueprint/courses", http.StatusBadRequest},
		testCase{"blueprint remove courses with empty type should return 400",
			"DELETE", "/blueprint/courses?type=", http.StatusBadRequest},
		testCase{"blueprint remove courses with invalid type should return 400",
			"DELETE", "/blueprint/courses?type=lorem", http.StatusBadRequest},
		testCase{"blueprint remove courses without year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&semester=1", http.StatusBadRequest},
		testCase{"blueprint remove courses with empty year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=&semester=1", http.StatusBadRequest},
		testCase{"blueprint remove courses with invalid year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=lorem&semester=1", http.StatusBadRequest},
		testCase{"blueprint remove courses with negative year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=-1&semester=1", http.StatusBadRequest},
		testCase{"blueprint remove courses with non-existent year should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=5&semester=1", http.StatusBadRequest},
		testCase{"blueprint remove courses without semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1", http.StatusBadRequest},
		testCase{"blueprint remove courses with empty semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=", http.StatusBadRequest},
		testCase{"blueprint remove courses with invalid semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=lorem", http.StatusBadRequest},
		testCase{"blueprint remove courses with negative semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=-1", http.StatusBadRequest},
		testCase{"blueprint remove courses with non-existent semester should return 400",
			"DELETE", "/blueprint/courses?type=semester-remove&year=1&semester=5", http.StatusBadRequest},
		testCase{"blueprint remove courses without selected should return 400",
			"DELETE", "/blueprint/courses?type=selected-remove", http.StatusBadRequest},
	}

	runTests(t, tests)
}

func TestDegreePlanDetailServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"degree plan page without saved plan should return 200",
			"GET", "/degreeplan/", http.StatusOK},
		testCase{"degree plan page with cs language without saved plan should return 200",
			"GET", "/cs/degreeplan/", http.StatusOK},
		testCase{"degree plan page with en language without saved plan should return 200",
			"GET", "/en/degreeplan/", http.StatusOK},
		testCase{"degree plan existing code NIDAW19B should return 200",
			"GET", "/degreeplan/NIDAW19B", http.StatusOK},
		testCase{"degree plan existing code NISD23N should return 200",
			"GET", "/degreeplan/NISD23N", http.StatusOK},
		testCase{"degree plan existing code NIPP19B should return 200",
			"GET", "/degreeplan/NIPP19B", http.StatusOK},
		testCase{"degree plan existing code NMOM24B should return 200",
			"GET", "/degreeplan/NMOM24B", http.StatusOK},
		testCase{"degree plan existing code NFEBCHF20N should return 200",
			"GET", "/degreeplan/NFEBCHF20N", http.StatusOK},
		testCase{"degree plan existing code NIPPA19B should return 200",
			"GET", "/degreeplan/NIPPA19B", http.StatusOK},
		testCase{"degree plan save NIPP19B should return 200",
			"PATCH", "/degreeplan/NIPP19B", http.StatusOK},
		testCase{"degree plan save NISD23N should return 200",
			"PATCH", "/degreeplan/NISD23N", http.StatusOK},
		testCase{"degree plan save NIDAW19B should return 200",
			"PATCH", "/degreeplan/NIDAW19B", http.StatusOK},
		testCase{"degree plan save NIPVS19B should return 200",
			"PATCH", "/degreeplan/NIPVS19B", http.StatusOK},
		testCase{"degree plan page with saved plan should return 200",
			"GET", "/degreeplan/", http.StatusOK},
		testCase{"degree plan page with cs language with saved plan should return 200",
			"GET", "/cs/degreeplan/", http.StatusOK},
		testCase{"degree plan page with en language with saved plan should return 200",
			"GET", "/en/degreeplan/", http.StatusOK},
		testCase{"degree plan remove saved plan should return 200",
			"DELETE", "/degreeplan/", http.StatusOK},
		testCaseWithReferer{"degree plan merge recommended plan with blueprint should return 200",
			"PATCH", "/degreeplan/plan-to-blueprint/NIPVS19B?maxYear=3", "/degreeplan/NIPVS19B", http.StatusOK},
		testCaseWithReferer{"degree plan rewrite blueprint with recommended plan should return 200",
			"PUT", "/degreeplan/plan-to-blueprint/NIPVS19B?maxYear=3", "/degreeplan/NIPVS19B", http.StatusOK},
		testCaseWithReferer{"degree plan merge empty recommended plan with blueprint should return 200",
			"PATCH", "/degreeplan/plan-to-blueprint/NIUI20N?maxYear=3", "/degreeplan/NIUI20N", http.StatusOK},
		testCaseWithReferer{"degree plan merge non-existing plan with blueprint should return 200",
			"PATCH", "/degreeplan/plan-to-blueprint/lorem?maxYear=3", "/degreeplan/NIPVS19B", http.StatusOK},
		testCaseWithReferer{"degree plan rewrite blueprint with one year plan should return 200",
			"PUT", "/degreeplan/plan-to-blueprint/NISD23N?maxYear=1", "/degreeplan/NISD23N", http.StatusOK},
		testCaseWithReferer{"degree plan add course to blueprint should return 200",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=0", "/degreeplan/NIPVS19B", http.StatusOK},
		testCaseWithReferer{"degree plan add courses to blueprint should return 200",
			"POST", "/degreeplan/blueprint?selected-courses=NSWI035&selected-courses=NPRG024&year=0&semester=0", "/degreeplan/NIPVS19B", http.StatusOK},

		// remove BP years
		testCase{"blueprint remove third year should return 200",
			"DELETE", "/blueprint/year?unassign=false", http.StatusOK},
		testCase{"blueprint remove second year should return 200",
			"DELETE", "/blueprint/year?unassign=false", http.StatusOK},
		testCase{"blueprint remove first year should return 200",
			"DELETE", "/blueprint/year?unassign=false", http.StatusOK},

		// Errors
		testCase{"degree plan page for non-existent plan should return 404",
			"GET", "/degreeplan/lorem", http.StatusNotFound},
		testCase{"degree plan save with invalid code should return 400",
			"PATCH", "/degreeplan/lorem", http.StatusBadRequest},
		testCase{"degree plan removing non-existent saved plan should return 404",
			"DELETE", "/degreeplan/", http.StatusNotFound},
		testCaseWithReferer{"degree plan add same course to blueprint should return 409",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=0", "/degreeplan/NIPVS19B", http.StatusConflict},
		testCaseWithReferer{"degree plan add course to blueprint without course should return 400",
			"POST", "/degreeplan/blueprint?year=0&semester=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint without year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&semester=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint without semester should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint with invalid course should return 400",
			"POST", "/degreeplan/blueprint?course=XXNOTXX&year=0&semester=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint with invalid year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=asdf&semester=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint with negative year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=-1&semester=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint with non-existent year should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=1&semester=0", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint with invalid semester should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=asdf", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan add course to blueprint with negative semester should return 400",
			"POST", "/degreeplan/blueprint?course=NSWI120&year=0&semester=-1", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan rewrite blueprint with empty recommended plan should return 400",
			"PUT", "/degreeplan/plan-to-blueprint/NIUI20N?maxYear=3", "/degreeplan/NIUI20N", http.StatusBadRequest},
		testCaseWithReferer{"degree plan rewrite blueprint with non-existing plan should return 400",
			"PUT", "/degreeplan/plan-to-blueprint/lorem?maxYear=3", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan merge recommended plan with no year should return 400",
			"PATCH", "/degreeplan/plan-to-blueprint/NIPVS19B", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan merge recommended plan with no int year return 400",
			"PATCH", "/degreeplan/plan-to-blueprint/NIPVS19B?maxYear=lorem", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan merge recommended plan with low year should return 400",
			"PATCH", "/degreeplan/plan-to-blueprint/NIPVS19B?maxYear=1", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan rewrite blueprint with no year should return 400",
			"PUT", "/degreeplan/plan-to-blueprint/NIPVS19B", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan rewrite blueprint with no int year return 400",
			"PUT", "/degreeplan/plan-to-blueprint/NIPVS19B?maxYear=lorem", "/degreeplan/NIPVS19B", http.StatusBadRequest},
		testCaseWithReferer{"degree plan rewrite blueprint with low year should return 400",
			"PUT", "/degreeplan/plan-to-blueprint/NIPVS19B?maxYear=1", "/degreeplan/NIPVS19B", http.StatusBadRequest},
	}
	runTests(t, tests)
}

func TestDegreePlansSearchServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"degree plan search page should return 200",
			"GET", "/degreeplans", http.StatusOK},
		testCase{"degree plan search page empty query should return 200",
			"GET", "/degreeplans/?search-dp-query=", http.StatusOK},
		testCase{"degree plan search page query 'software' should return 200",
			"GET", "/degreeplans/?search-dp-query=software", http.StatusOK},
		testCase{"degree plan search page query 'a' should return 200",
			"GET", "/degreeplans/?search-dp-query=a", http.StatusOK},
		testCase{"degree plan search page query 'NIPVS19B' should return 200",
			"GET", "/degreeplans/?search-dp-query=NIPVS19B", http.StatusOK},
		testCase{"degree plan search missing query should return 200",
			"GET", "/degreeplans/search", http.StatusOK},
		testCase{"degree plan search empty query should return 200",
			"GET", "/degreeplans/search?search-dp-query=", http.StatusOK},
		testCase{"degree plan search by code NIDAW19B should return 200",
			"GET", "/degreeplans/search?search-dp-query=NIDAW19B", http.StatusOK},
		testCase{"degree plan search by 'softwarove' should return 200",
			"GET", "/degreeplans/search?search-dp-query=softwarove", http.StatusOK},
		testCase{"degree plan search by 'a' should return 200",
			"GET", "/degreeplans/search?search-dp-query=a", http.StatusOK},

		// Errors
		testCase{"degree plan search with wrong path should return 404",
			"GET", "/degreeplans/lorem", http.StatusNotFound},
		testCase{"degree plan search with wrong path should return 404",
			"GET", "/degreeplans/search/lorem", http.StatusNotFound},
	}

	runTests(t, tests)
}

func TestDegreePlansCompareServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"degree plan compare NIDAW19B and NISD23N should return 200",
			"GET", "/degreeplans/compare/NIDAW19B/NISD23N", http.StatusOK},
		testCase{"degree plan compare NIPP19B and NMOM24B should return 200",
			"GET", "/degreeplans/compare/NIPP19B/NMOM24B", http.StatusOK},
		testCase{"degree plan compare NFEBCHF20N and NIPPA19B should return 200",
			"GET", "/degreeplans/compare/NFEBCHF20N/NIPPA19B", http.StatusOK},
		testCase{"degree plan compare same plans NIPVS19B and NIPVS19B should return 200",
			"GET", "/degreeplans/compare/NIPVS19B/NIPVS19B", http.StatusOK},

		// Errors
		testCase{"degree plan compare with wrong path should return 404",
			"GET", "/degreeplans/compare/NIDAW19B", http.StatusNotFound},
		testCase{"degree plan compare with wrong path should return 404",
			"GET", "/degreeplans/compare", http.StatusNotFound},
		testCase{"degree plan compare non-existent plans should return 404",
			"GET", "/degreeplans/compare/NOTEXIST1/NOTEXIST2", http.StatusNotFound},
		testCase{"degree plan compare existing and non-existent plan should return 404",
			"GET", "/degreeplans/compare/NIDAW19B/NOTEXIST2", http.StatusNotFound},
		testCase{"degree plan compare non-existent and existing plan should return 404",
			"GET", "/degreeplans/compare/NOTEXIST1/NISD23N", http.StatusNotFound},
	}

	runTests(t, tests)
}

func TestPageServer(t *testing.T) {
	tests := []testRunner{
		// Happy path
		testCase{"quicksearching 'NSWI153' should return 200",
			"GET", "/page/quicksearch?q=NSWI153", http.StatusOK},
		testCase{"quicksearching 'NSWI153' with cs should return 200",
			"GET", "/cs/page/quicksearch?q=NSWI153", http.StatusOK},
		testCase{"quicksearching 'NSWI153' with en should return 200",
			"GET", "/en/page/quicksearch?q=NSWI153", http.StatusOK},
		testCase{"quicksearching 'NPRG024' should return 200",
			"GET", "/page/quicksearch?q=NPRG024", http.StatusOK},
		testCase{"quicksearching 'peska' should return 200",
			"GET", "/page/quicksearch?q=peska", http.StatusOK},
		testCase{"quicksearching random text should return 200",
			"GET", "/page/quicksearch?q=asagafe45", http.StatusOK},
		testCase{"quicksearching empty query should return 200",
			"GET", "/page/quicksearch?q=", http.StatusOK},
		testCase{"quicksearching with missing query equals empty query should return 200",
			"GET", "/page/quicksearch", http.StatusOK},
	}

	runTests(t, tests)
}

//================================================================================
// Test Utilities
//================================================================================

type testRunner interface {
	getName() string
	getMethod() string
	getURL() string
	getReferer() string
	getWant() int
}

type testCase struct {
	name   string
	method string
	url    string
	want   int
}

func (tc testCase) getName() string    { return tc.name }
func (tc testCase) getMethod() string  { return tc.method }
func (tc testCase) getURL() string     { return tc.url }
func (tc testCase) getReferer() string { return "" }
func (tc testCase) getWant() int       { return tc.want }

type testCaseWithReferer struct {
	name    string
	method  string
	url     string
	referer string
	want    int
}

func (tc testCaseWithReferer) getName() string    { return tc.name }
func (tc testCaseWithReferer) getMethod() string  { return tc.method }
func (tc testCaseWithReferer) getURL() string     { return tc.url }
func (tc testCaseWithReferer) getReferer() string { return tc.referer }
func (tc testCaseWithReferer) getWant() int       { return tc.want }

func runTests(t *testing.T, tests []testRunner) {
	ts := setupTestServer(t)
	defer ts.Close()

	client := ts.Client()
	sessionCookie := setupTestUser(ts, t)

	for _, test := range tests {
		t.Run(test.getName(), func(t *testing.T) {
			req, err := http.NewRequest(test.getMethod(), ts.URL+test.getURL(), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.AddCookie(sessionCookie)

			if test.getReferer() != "" {
				req.Header.Set("Referer", ts.URL+test.getReferer())
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, test.getWant(), resp.StatusCode)
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
