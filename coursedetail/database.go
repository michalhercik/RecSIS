package coursedetail

import (
	"database/sql"
	"fmt"

	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
)

type DBLang int

const (
	cs DBLang = iota
	en
)

var DBLangs = map[DBLang]string{
	cs: "CZE",
	en: "ENG",
}

func (l DBLang) String() string {
	return DBLangs[l]
}

type DBManager struct {
	DB *sql.DB
}

func selectCourse(tx *sql.Tx, code string, lang DBLang) (*Course, error) {
	row := tx.QueryRow(sqlquery.Course, code, lang.String())
	course := &Course{
		Faculty:                  Faculty{},
		Guarantors:               []Teacher{{}, {}, {}},
		Annotation:               Description{},
		CompletionRequirements:   Description{},
		ExamRequirements:         Description{},
		Sylabus:                  Description{},
	}

	err := row.Scan(
		&course.Code,
		&course.Name,
		&course.Faculty.SisID,
		&course.Faculty.Name,
		&course.Faculty.Abbr,
		&course.GuarantorDepartment,
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
		&course.Guarantors[0].SisID,
		&course.Guarantors[0].FirstName,
		&course.Guarantors[0].LastName,
		&course.Guarantors[0].TitleBefore,
		&course.Guarantors[0].TitleAfter,
		&course.Guarantors[1].SisID,
		&course.Guarantors[1].FirstName,
		&course.Guarantors[1].LastName,
		&course.Guarantors[1].TitleBefore,
		&course.Guarantors[1].TitleAfter,
		&course.Guarantors[2].SisID,
		&course.Guarantors[2].FirstName,
		&course.Guarantors[2].LastName,
		&course.Guarantors[2].TitleBefore,
		&course.Guarantors[2].TitleAfter,
		&course.MinEnrollment,
		&course.Capacity,
		&course.Annotation.title,
		&course.Annotation.content,
		&course.CompletionRequirements.title,
		&course.CompletionRequirements.content,
		&course.ExamRequirements.title,
		&course.ExamRequirements.content,
		&course.Sylabus.title,
		&course.Sylabus.content,
	)
	if err != nil {
		return course, fmt.Errorf("selectCourse: %v", err)
	}
	return course, nil
}

func selectTeachers(tx *sql.Tx, course *Course) error {
	rows, err := tx.Query(sqlquery.CourseTeachers, course.Code)
	if err != nil {
		return fmt.Errorf("selectTeachers: %v", err)
	}
	defer rows.Close()
	course.Teachers = []Teacher{}
	for rows.Next() {
		var teacher Teacher
		if err := rows.Scan(
			&teacher.SisID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.TitleBefore,
			&teacher.TitleAfter,
		); err != nil {
			return fmt.Errorf("selectTeachers: %v", err)
		}
		course.Teachers = append(course.Teachers, teacher)
	}
	return nil
}

func (reader DBManager) Course(code string, lang DBLang) (*Course, error) {
	tx, err := reader.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	course, err := selectCourse(tx, code, lang)
	if err != nil {
		return nil, err
	}
	if err = selectTeachers(tx, course); err != nil {
		return course, err
	}
	if err = tx.Commit(); err != nil {
		return course, err
	}
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
		{ID: 4, UserID: 4, Content: "I think that Michal is a great name"},
	}
	return comments, nil
}