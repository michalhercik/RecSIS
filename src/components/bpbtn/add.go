package bpbtn

import (
	"fmt"
	"iter"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type ViewModel struct {
	course     string
	semesters  []bool
	hxPostBase string
	hxSwap     string
	hxTarget   string
	hxInclude  string
}

type DoubleAdd struct {
	Add
	TemplSecond func(ViewModel, text) templ.Component
}

func (b DoubleAdd) PartialComponentSecond(lang language.Language) func(string, string, string, []bool, string) templ.Component {
	return func(hxSwap, hxTarget, hxInclude string, semesters []bool, course string) templ.Component {
		t := texts[lang]
		model := ViewModel{
			course:     course,
			semesters:  semesters,
			hxPostBase: b.Options.HxPostBase,
			hxSwap:     hxSwap,
			hxTarget:   hxTarget,
			hxInclude:  hxInclude,
		}
		return b.TemplSecond(model, t)
	}
}

type Add struct {
	DB      *sqlx.DB
	Templ   func(ViewModel, text) templ.Component
	Options Options
}

func (b Add) Component(semesters []bool, lang language.Language, course string) templ.Component {
	t := texts[lang]
	model := ViewModel{
		course:     course,
		semesters:  semesters,
		hxPostBase: b.Options.HxPostBase,
		hxSwap:     b.Options.HxSwap,
		hxTarget:   b.Options.HxTarget,
		hxInclude:  b.Options.HxInclude,
	}
	return b.Templ(model, t)
}

func (b Add) PartialComponent(lang language.Language) func(string, string, string, []bool, string) templ.Component {
	return func(hxSwap, hxTarget, hxInclude string, semesters []bool, course string) templ.Component {
		t := texts[lang]
		model := ViewModel{
			course:     course,
			semesters:  semesters,
			hxPostBase: b.Options.HxPostBase,
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
	courses, err = b.CoursesFromRequest(r, additionalCourses, t)
	if err != nil {
		return courses, year, semester, errorx.AddContext(err)
	}
	year, err = b.YearFromRequest(r, t)
	if err != nil {
		return courses, year, semester, errorx.AddContext(err)
	}
	semester, err = b.SemesterFromRequest(r, t)
	if err != nil {
		return courses, year, semester, errorx.AddContext(err)
	}
	return courses, year, semester, nil
}

// if additionalCourses is nil (or empty) then always returns 1 course or error
func (b Add) CoursesFromRequest(r *http.Request, additionalCourses []string, t text) ([]string, error) {
	courses := additionalCourses
	course := r.FormValue("course")
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

func (b Add) YearFromRequest(r *http.Request, t text) (int, error) {
	yearString := r.FormValue("year")
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
			errorx.AddContext(fmt.Errorf("invalid year: %w", err), errorx.P("year", yearString)),
			http.StatusBadRequest,
			t.errInvalidYear,
		)
	}
	return year, nil
}

func (b Add) SemesterFromRequest(r *http.Request, t text) (int, error) {
	semesterString := r.FormValue("semester")
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
			errorx.AddContext(fmt.Errorf("invalid semester: %w", err), errorx.P("semester", semesterString)),
			http.StatusBadRequest,
			t.errInvalidSemester,
		)
	}
	return semester, nil
}

const (
	uniqueViolationCode       = "23505"
	duplicateCoursesViolation = "blueprint_courses_blueprint_semester_id_course_code_key"
)

func (b Add) Action(userID string, year int, semester int, lang language.Language, courses ...string) ([]int, error) {
	const insertCourse = `--sql
		WITH target_position AS (
			SELECT
				bs.id AS blueprint_semester_id,
				COALESCE(bc.position, 0) + 1 AS last_position
			FROM blueprint_years y
			LEFT JOIN blueprint_semesters bs
				ON y.id = bs.blueprint_year_id
			LEFT JOIN blueprint_courses bc
				ON bs.id = bc.blueprint_semester_id
			WHERE y.user_id = $1
				AND y.academic_year = $2
				AND bs.semester = $3
			ORDER BY bc.position DESC
			LIMIT 1
		)
		INSERT INTO blueprint_courses(blueprint_semester_id, course_code, course_valid_from, position)
		SELECT
			tp.blueprint_semester_id,
			c.code,
			c.valid_from,
			tp.last_position + ROW_NUMBER() OVER (ORDER BY c.code)
		FROM
			target_position tp,
			UNNEST($4::text[]) AS course_code
			JOIN LATERAL (
				SELECT code, valid_from FROM courses WHERE code = course_code ORDER BY valid_from DESC LIMIT 1
			) c ON TRUE
		RETURNING id;
		`
	var courseIDs []int
	err := b.DB.Select(&courseIDs, insertCourse, userID, year, semester, pq.StringArray(courses))
	if err != nil {
		// Handle unique violation for blueprint_semester_id, course_code
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode && pqErr.Constraint == duplicateCoursesViolation {
			userErrMsg := texts[lang].errDuplicateCourseInBP
			if len(courses) > 1 {
				userErrMsg = texts[lang].errDuplicateCoursesInBP
			}
			return []int{}, errorx.NewHTTPErr(
				errorx.AddContext(err, errorx.P("year", year), errorx.P("semester", semester), errorx.P("courses", strings.Join(courses, ","))),
				http.StatusConflict,
				userErrMsg,
			)
		}
		return []int{}, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("year", year), errorx.P("semester", semester), errorx.P("courses", strings.Join(courses, ","))),
			http.StatusInternalServerError,
			texts[lang].errAddCourseToBPFailed,
		)
	}
	return courseIDs, nil
}

type Options struct {
	HxPostBase string
	HxSwap     string
	HxTarget   string
	HxInclude  string
}

func (o Options) With(hxSwap, hxTarget, hxInclude string) Options {
	return Options{
		HxPostBase: o.HxPostBase,
		HxSwap:     hxSwap,
		HxTarget:   hxTarget,
		HxInclude:  hxInclude,
	}
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

type overYearsIterator struct {
	disableWinter bool
	disableSummer bool
}
