package bpbtn

/** PACKAGE DESCRIPTION

The bpbtn package provides reusable components and logic for adding courses to a user's blueprint (own study plan) in the application. Its main purpose is to encapsulate the UI and backend logic for the "add to blueprint" button, including request parsing, validation, error handling, and database operations. This package is designed to be injected into servers (such as courses or degree plan) so that the add button can be rendered and its actions handled consistently across different parts of the app.

Typical usage:
You create an instance of Add (or AddWithTwoTemplComponents for more complex UI needs), configure it with a database connection, a templating function for rendering the button, and options for HTMX integration. The server calls PartialComponent() to get a function that renders the add button, passing in flag slice with available semesters to assign to and course code. When a user clicks the button, the server uses ParseRequest to extract the course(s), year, and semester to assign to from the HTTP request, then calls Action to add the course(s) to the blueprint in the database. The package handles errors (such as duplicates or missing data) and returns localized messages using the errorx package.

This approach allows developers to easily add blueprint functionality to any course listing or detail page, keeping the UI and backend logic consistent and maintainable.

*/

import (
	"fmt"
	"iter"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/components/bpbtn/internal/sqlquery"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type AddWithTwoTemplComponents struct {
	Add
	TemplSecond func(ViewModel, text) templ.Component
}

func (b AddWithTwoTemplComponents) PartialComponentSecond(lang language.Language) func(string, string, string, []bool, string) templ.Component {
	return func(hxSwap, hxTarget, hxInclude string, semesters []bool, course string) templ.Component {
		t := texts[lang]
		model := ViewModel{
			course:     course,
			semesters:  semesters,
			hxPostBase: b.HxPostBase,
			hxSwap:     hxSwap,
			hxTarget:   hxTarget,
			hxInclude:  hxInclude,
		}
		return b.TemplSecond(model, t)
	}
}

type Add struct {
	DB         *sqlx.DB
	Templ      func(ViewModel, text) templ.Component
	HxPostBase string
}

func (b Add) Endpoint() string {
	return "POST /" + endpointPath
}

func (b Add) PartialComponent(lang language.Language) func(string, string, string, []bool, string) templ.Component {
	return func(hxSwap, hxTarget, hxInclude string, semesters []bool, course string) templ.Component {
		t := texts[lang]
		model := ViewModel{
			course:     course,
			semesters:  semesters,
			hxPostBase: b.HxPostBase,
			hxSwap:     hxSwap,
			hxTarget:   hxTarget,
			hxInclude:  hxInclude,
		}
		return b.Templ(model, t)
	}
}

func (b Add) ParseRequest(r *http.Request, additionalCourses []string) ([]string, int, int, error) {
	var (
		courses  []string
		year     int
		semester int
		err      error
	)
	t := texts[language.FromContext(r.Context())]
	courses, err = b.coursesFromRequest(r, additionalCourses, t)
	if err != nil {
		return courses, year, semester, errorx.AddContext(err)
	}
	year, err = b.yearFromRequest(r, t)
	if err != nil {
		return courses, year, semester, errorx.AddContext(err)
	}
	semester, err = b.semesterFromRequest(r, t)
	if err != nil {
		return courses, year, semester, errorx.AddContext(err)
	}
	return courses, year, semester, nil
}

func (b Add) coursesFromRequest(r *http.Request, additionalCourses []string, t text) ([]string, error) {
	courses := additionalCourses
	course := r.FormValue(courseParam)
	if course != "" {
		courses = append(courses, course)
	}
	if len(courses) == 0 {
		return courses, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("no courses provided in request"), errorx.P("additionalCourses", strings.Join(additionalCourses, ","))),
			http.StatusBadRequest,
			t.errNoCoursesProvided,
		)
	}
	return courses, nil
}

func (b Add) yearFromRequest(r *http.Request, t text) (int, error) {
	yearString := r.FormValue(yearParam)
	if yearString == "" {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("year not provided in request")),
			http.StatusBadRequest,
			t.errNoYearProvided,
		)
	}
	year, err := strconv.Atoi(yearString)
	if err != nil {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid year: %w", err), errorx.P(yearParam, yearString)),
			http.StatusBadRequest,
			t.errInvalidYear,
		)
	}
	return year, nil
}

func (b Add) semesterFromRequest(r *http.Request, t text) (int, error) {
	semesterString := r.FormValue(semesterParam)
	if semesterString == "" {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("semester not provided in request")),
			http.StatusBadRequest,
			t.errNoSemesterProvided,
		)
	}
	semester, err := strconv.Atoi(semesterString)
	if err != nil {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid semester: %w", err), errorx.P(semesterParam, semesterString)),
			http.StatusBadRequest,
			t.errInvalidSemester,
		)
	}
	return semester, nil
}

// TODO: make package `services` with operation with BP
// TODO: is possible to add course to multiple years/semesters at once?
func (b Add) Action(userID string, year int, semester int, lang language.Language, courses ...string) ([]int, error) {
	var courseIDs []int
	err := b.DB.Select(&courseIDs, sqlquery.InsertCourse, userID, year, semester, pq.StringArray(courses))
	if err != nil {
		// Handle unique violation for blueprint_semester_id, course_code
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode && pqErr.Constraint == duplicateCoursesViolation {
			userErrMsg := texts[lang].errDuplicateCourseInBP
			if len(courses) > 1 {
				userErrMsg = texts[lang].errDuplicateCoursesInBP
			}
			return []int{}, errorx.NewHTTPErr(
				errorx.AddContext(err, errorx.P(yearParam, year), errorx.P(semesterParam, semester), errorx.P("courses", strings.Join(courses, ","))),
				http.StatusConflict,
				userErrMsg,
			)
		}
		return []int{}, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P(yearParam, year), errorx.P(semesterParam, semester), errorx.P("courses", strings.Join(courses, ","))),
			http.StatusInternalServerError,
			texts[lang].errAddCourseToBPFailed,
		)
	}
	if len(courseIDs) == 0 {
		return []int{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("no rows were changed"), errorx.P(yearParam, year), errorx.P(semesterParam, semester), errorx.P("courses", strings.Join(courses, ","))),
			http.StatusBadRequest,
			texts[lang].errAddCourseToBPFailed,
		)
	}
	return courseIDs, nil
}

func iterateOverAssignedYears(semesters []bool) iter.Seq2[int, overYearsIterator] {
	if len(semesters) == 0 {
		return func(func(int, overYearsIterator) bool) {}
	}
	semesters = semesters[1:]
	return func(yield func(int, overYearsIterator) bool) {
		for i := 0; i < len(semesters); i += 2 {
			overYearsIterator := overYearsIterator{
				disableWinter: semesters[i],
				disableSummer: false,
			}
			if i+1 < len(semesters) {
				overYearsIterator.disableSummer = semesters[i+1]
			}
			if !yield((i/2)+1, overYearsIterator) {
				break
			}
		}
	}
}
