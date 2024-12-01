package coursedetail

import (
	"database/sql"

	"github.com/michalhercik/RecSIS/coursedetail/sqlquery"
)

type DbCourseReader struct {
	Db *sql.DB
}

func (reader DbCourseReader) Course(code string) (*Course, error) {
	row := reader.Db.QueryRow(sqlquery.Course, code)
	course := newCourse()
	err := row.Scan(
		&course.Id,
		&course.Code,
		&course.NameCs,
		&course.NameEn,
		&course.ValidFrom,
		&course.ValidTo,
		&course.Faculty.Id,
		&course.Faculty.SisId,
		&course.Faculty.NameCs,
		&course.Faculty.NameEn,
		&course.Faculty.Abbr,
		&course.Guarantor,
		&course.State,
		&course.Start,
		&course.SemesterCount,
		&course.Language,
		&course.LectureRange1,
		&course.SeminarRange1,
		&course.LectureRange2,
		&course.SeminarRange2,
		&course.ExamType,
		&course.Credits,
		&course.Teachers[0].Id,
		&course.Teachers[0].SisId,
		&course.Teachers[0].Department,
		&course.Teachers[0].Faculty.Id,
		&course.Teachers[0].Faculty.SisId,
		&course.Teachers[0].Faculty.NameCs,
		&course.Teachers[0].Faculty.NameEn,
		&course.Teachers[0].Faculty.Abbr,
		&course.Teachers[0].FirstName,
		&course.Teachers[0].LastName,
		&course.Teachers[0].TitleBefore,
		&course.Teachers[0].TitleAfter,
		&course.Teachers[1].Id,
		&course.Teachers[1].SisId,
		&course.Teachers[1].Department,
		&course.Teachers[1].Faculty.Id,
		&course.Teachers[1].Faculty.SisId,
		&course.Teachers[1].Faculty.NameCs,
		&course.Teachers[1].Faculty.NameEn,
		&course.Teachers[1].Faculty.Abbr,
		&course.Teachers[1].FirstName,
		&course.Teachers[1].LastName,
		&course.Teachers[1].TitleBefore,
		&course.Teachers[1].TitleAfter,
		&course.MinEnrollment,
		&course.Capacity,
		&course.AnnotationCs,
		&course.AnnotationEn,
		&course.SylabusCs,
		&course.SylabusEn,
	)
	return course, err
}
