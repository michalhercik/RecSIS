package bpbtn

import (
	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/dbcourse"
	"github.com/michalhercik/RecSIS/language"
)

type Add struct {
	DB      *sqlx.DB
	Templ   func(string, int, Text, Options) templ.Component
	Options Options
}

func (b Add) Component(course string, numberOfYears int, lang language.Language) templ.Component {
	t := texts[lang]
	return b.Templ(course, numberOfYears, t, b.Options)
}

func (b Add) PartialComponent(numberOfYears int, lang language.Language) func(string, string, string) templ.Component {
	return func(course, hxSwap, hxTarget string) templ.Component {
		t := texts[lang]
		options := b.Options.With(hxSwap, hxTarget)
		return b.Templ(course, numberOfYears, t, options)
	}
}

func (b Add) NumberOfYears(userID string) (int, error) {
	const NumberOfBlueprintYears = `SELECT COUNT(*) FROM blueprint_years WHERE user_id = $1`
	var numberOfYears int
	err := b.DB.Get(&numberOfYears, NumberOfBlueprintYears, userID)
	if err != nil {
		return 0, err
	}
	numberOfYears -= 1
	return numberOfYears, nil
}

func (b Add) Action(userID, course string, year int, semester dbcourse.SemesterAssignment) (int, error) {
	const InsertCourse = `--sql
		WITH target_position AS (
			SELECT bs.id AS blueprint_semester_id, COALESCE(bc.position, 0) + 1 AS position FROM blueprint_years y
			LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
			LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
			WHERE y.user_id=$1
			AND y.academic_year=$2
			AND bs.semester=$3
			ORDER BY bc.position DESC
			LIMIT 1
		),
		target_course AS (
			SELECT code, valid_from FROM courses WHERE code=$4 ORDER BY valid_from DESC LIMIT 1
		)
		INSERT INTO blueprint_courses(blueprint_semester_id, course_code, course_valid_from, position)
		VALUES (
			(SELECT blueprint_semester_id FROM target_position),
			(SELECT code FROM target_course),
			(SELECT valid_from FROM target_course),
			(SELECT position FROM target_position)
		)
		RETURNING id
		;
		`
	var courseID int
	err := b.DB.Get(&courseID, InsertCourse, userID, year, int(semester), course)
	if err != nil {
		return 0, err
	}
	return courseID, nil
}

type Options struct {
	HxPostBase string
	HxSwap     string
	HxTarget   string
}

func (o Options) With(hxSwap, hxTarget string) Options {
	return Options{
		HxPostBase: o.HxPostBase,
		HxSwap:     hxSwap,
		HxTarget:   hxTarget,
	}
}
