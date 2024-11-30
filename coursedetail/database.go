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
		&course.Teacher1.Id,
		&course.Teacher1.SisId,
		&course.Teacher1.Department,
		&course.Teacher1.Faculty.Id,
		&course.Teacher1.Faculty.SisId,
		&course.Teacher1.Faculty.NameCs,
		&course.Teacher1.Faculty.NameEn,
		&course.Teacher1.Faculty.Abbr,
		&course.Teacher1.FirstName,
		&course.Teacher1.LastName,
		&course.Teacher1.TitleBefore,
		&course.Teacher1.TitleAfter,
		&course.Teacher2.Id,
		&course.Teacher2.SisId,
		&course.Teacher2.Department,
		&course.Teacher2.Faculty.Id,
		&course.Teacher2.Faculty.SisId,
		&course.Teacher2.Faculty.NameCs,
		&course.Teacher2.Faculty.NameEn,
		&course.Teacher2.Faculty.Abbr,
		&course.Teacher2.FirstName,
		&course.Teacher2.LastName,
		&course.Teacher2.TitleBefore,
		&course.Teacher2.TitleAfter,
		&course.MinEnrollment,
		&course.Capacity,
		&course.AnnotationCs,
		&course.AnnotationEn,
		&course.SylabusCs,
		&course.SylabusEn,
	)
	return course, err
}
