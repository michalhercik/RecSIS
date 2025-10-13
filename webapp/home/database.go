package home

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/home/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type courses []struct {
	dbds.Course
	InDegreePlan bool `db:"in_degree_plan"`
}

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) courses(userID string, courseCodes []string, lang language.Language) ([]course, error) {
	var result courses
	if err := m.DB.Select(&result, sqlquery.Courses, userID, pq.Array(courseCodes), lang); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Courses: %w", err), errorx.P("courseCodes", strings.Join(courseCodes, ",")), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotLoadCourses,
		)
	}
	courses := intoCourses(result)
	return courses, nil
}

func (m DBManager) testAccounts() ([]string, error) {
	var result []string
	if err := m.DB.Select(&result, sqlquery.TestAccounts); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.TestAccounts: %w", err)),
			http.StatusInternalServerError,
			"Cannot load test accounts",
		)
	}
	return result, nil
}

func intoCourses(from courses) []course {
	result := make([]course, len(from))
	for i, course := range from {
		result[i].Code = course.Code
		result[i].Title = course.Title
		result[i].Semester = teachingSemester(course.Start)
		result[i].LectureRangeWinter = course.LectureRangeWinter
		result[i].SeminarRangeWinter = course.SeminarRangeWinter
		result[i].LectureRangeSummer = course.LectureRangeSummer
		result[i].SeminarRangeSummer = course.SeminarRangeSummer
		result[i].ExamType = course.ExamType
		result[i].Credits = course.Credits
		result[i].Guarantors = intoTeacherSlice(course.Guarantors)
		// result[i].InDegreePlan = course.InDegreePlan
	}
	return result
}

func intoTeacherSlice(from dbds.TeacherSlice) []teacher {
	result := make([]teacher, len(from))
	for i, t := range from {
		result[i] = teacher{
			SisID:       t.SisID,
			FirstName:   t.FirstName,
			LastName:    t.LastName,
			TitleBefore: t.TitleBefore,
			TitleAfter:  t.TitleAfter,
		}
	}
	return result
}
