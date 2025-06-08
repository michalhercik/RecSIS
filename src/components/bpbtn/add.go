package bpbtn

import (
	"fmt"
	"iter"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/language"
)

type DoubleAdd struct {
	Add
	TemplSecond func(ViewModel, text) templ.Component
}

func (b DoubleAdd) PartialComponentSecond(lang language.Language) func(string, string, string, []bool, ...string) templ.Component {
	return func(hxSwap, hxTarget, hxInclude string, semesters []bool, course ...string) templ.Component {
		t := texts[lang]
		model := ViewModel{
			courses:    course,
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

func (b Add) Component(semesters []bool, lang language.Language, course ...string) templ.Component {
	t := texts[lang]
	model := ViewModel{
		courses:    course,
		semesters:  semesters,
		hxPostBase: b.Options.HxPostBase,
		hxSwap:     b.Options.HxSwap,
		hxTarget:   b.Options.HxTarget,
		hxInclude:  b.Options.HxInclude,
	}
	return b.Templ(model, t)
}

func (b Add) PartialComponent(lang language.Language) func(string, string, string, []bool, ...string) templ.Component {
	return func(hxSwap, hxTarget, hxInclude string, semesters []bool, course ...string) templ.Component {
		t := texts[lang]
		model := ViewModel{
			courses:    course,
			semesters:  semesters,
			hxPostBase: b.Options.HxPostBase,
			hxSwap:     hxSwap,
			hxTarget:   hxTarget,
			hxInclude:  hxInclude,
		}
		return b.Templ(model, t)
	}
}

func (b Add) ParseRequest(r *http.Request) ([]string, int, int, error) {
	var (
		courses  []string
		year     int
		semester int
		err      error
	)
	courses, err = b.CoursesFromRequest(r)
	if err != nil {
		return courses, year, semester, fmt.Errorf("failed to parse courses from request: %w", err)
	}
	year, err = b.YearFromRequest(r)
	if err != nil {
		return courses, year, semester, fmt.Errorf("failed to parse year from request: %w", err)
	}
	semester, err = b.SemesterFromRequest(r)
	if err != nil {
		return courses, year, semester, fmt.Errorf("failed to parse semester from request: %w", err)
	}
	return courses, year, semester, nil
}

func (b Add) YearFromRequest(r *http.Request) (int, error) {
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		return 0, fmt.Errorf("invalid year: %w", err)
	}
	return year, nil
}

func (b Add) SemesterFromRequest(r *http.Request) (int, error) {
	semester, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		return 0, fmt.Errorf("invalid semester: %w", err)
	}
	return semester, nil
}

func (b Add) CoursesFromRequest(r *http.Request) ([]string, error) {
	r.ParseForm()
	result := r.Form["courses"]
	if len(result) == 0 {
		return result, fmt.Errorf("no courses provided")
	}
	return result, nil
}

func (b Add) Action(userID string, year int, semester int, courses ...string) ([]int, error) {
	const InsertCourse = `--sql
		WITH target_position AS (
			SELECT bs.id AS blueprint_semester_id, COALESCE(bc.position, 0) + 1 AS last_position
			FROM blueprint_years y
			LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
			LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
			WHERE y.user_id=$1
			AND y.academic_year=$2
			AND bs.semester=$3
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
				SELECT code, valid_from FROM courses WHERE code=course_code ORDER BY valid_from DESC LIMIT 1
			) c ON TRUE
		RETURNING id;
		`
	var courseID []int
	err := b.DB.Select(&courseID, InsertCourse, userID, year, int(semester), pq.StringArray(courses))
	if err != nil {
		return []int{}, err
	}
	return courseID, nil
}

type ViewModel struct {
	courses    []string
	semesters  []bool
	hxPostBase string
	hxSwap     string
	hxTarget   string
	hxInclude  string
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
