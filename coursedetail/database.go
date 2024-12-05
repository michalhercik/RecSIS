package coursedetail

import (
	"database/sql"

	"github.com/michalhercik/RecSIS/coursedetail/sqlquery"
)

type DBManager struct {
	DB *sql.DB
}

func (reader DBManager) Course(code string) (*Course, error) {
	row := reader.DB.QueryRow(sqlquery.Course, code)
	course := newCourse()
	err := row.Scan(
		&course.ID,
		&course.Code,
		&course.NameCs,
		&course.NameEn,
		&course.ValidFrom,
		&course.ValidTo,
		&course.Faculty.ID,
		&course.Faculty.SisID,
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
		&course.Teachers[0].ID,
		&course.Teachers[0].SisID,
		&course.Teachers[0].Department,
		&course.Teachers[0].Faculty.ID,
		&course.Teachers[0].Faculty.SisID,
		&course.Teachers[0].Faculty.NameCs,
		&course.Teachers[0].Faculty.NameEn,
		&course.Teachers[0].Faculty.Abbr,
		&course.Teachers[0].FirstName,
		&course.Teachers[0].LastName,
		&course.Teachers[0].TitleBefore,
		&course.Teachers[0].TitleAfter,
		&course.Teachers[1].ID,
		&course.Teachers[1].SisID,
		&course.Teachers[1].Department,
		&course.Teachers[1].Faculty.ID,
		&course.Teachers[1].Faculty.SisID,
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

	// TODO: this is mock - change to real data
	course.Comments = []Comment{
		{ID: 1, UserID: 1, Content: "This is a comment"},
		{ID: 2, UserID: 2, Content: "This is another comment"},
		{ID: 3, UserID: 3, Content: "This is yet another comment"},
	}
	course.Ratings = []Rating{
		{ID: 1, UserID: 1, Rating: 1},
		{ID: 2, UserID: 2, Rating: 1},
		{ID: 3, UserID: 3, Rating: -1},
	}

	return course, err
}

// TODO: MOCK - implement
func (reader DBManager) AddComment(code, commentContent string) error {
	//_, err := reader.DB.Exec(sqlquery.AddComment, code, commentContent)
	//return err
	return nil
}

// TODO: MOCK - implement
func (reader DBManager) GetComments(code string) ([]Comment, error) {
	comments := []Comment{
		{ID: 1, UserID: 1, Content: "This is a comment"},
		{ID: 2, UserID: 2, Content: "This is another comment"},
		{ID: 3, UserID: 3, Content: "This is yet another comment"},
		{ID: 4, UserID: 4, Content: "I think that Michal is great name"},
	}
	return comments, nil
}
