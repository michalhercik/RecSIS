package blueprint

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/blueprint/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) Blueprint(userID string, lang DBLang) (*Blueprint, error) {
	var records []BlueprintRecord
	if err := m.DB.Select(&records, sqlquery.SelectCourses, userID, lang); err != nil {
		return nil, err
	}
	var bp Blueprint
	for _, record := range records {
		if err := bp.assign(record.BlueprintRecordPosition, record.NullCourse.Course()); err != nil {
			return nil, err
		}
	}
	return &bp, nil
}

func (m DBManager) NewCourse(userID string, course string, year int, semester SemesterAssignment) (int, error) {
	row := m.DB.QueryRow(sqlquery.InsertCourse, userID, year, int(semester), course)
	var courseID int
	err := row.Scan(&courseID)
	return courseID, err
}

func (m DBManager) InsertCourses(userID string, year int, semester SemesterAssignment, position int, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.MoveCourses, userID, pq.Array(courses), year, semester, position)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", count)
	}
	return nil
}

func (m DBManager) AppendCourses(userID string, year int, semester SemesterAssignment, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.AppendCourses, userID, year, int(semester), pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) UnassignYear(userID string, year int) error {
	_, err := m.DB.Exec(sqlquery.UnassignYear, userID, year)
	return err
}

func (m DBManager) UnassignSemester(userID string, year int, semester SemesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.UnassignSemester, userID, year, int(semester))
	return err
}

func (m DBManager) RemoveCourses(userID string, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesByID, userID, pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) RemoveCoursesBySemester(userID string, year int, semester SemesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesBySemester, userID, year, int(semester))
	return err
}

func (m DBManager) RemoveCoursesByYear(userID string, year int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesByYear, userID, year)
	return err
}

func (m DBManager) AddYear(userID string) error {
	fail := func(err error) error {
		return fmt.Errorf("AddYear: %v", err)
	}
	tx, err := m.DB.Beginx()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var newYearID int
	err = tx.Get(&newYearID, sqlquery.InsertYear, userID)
	if err != nil {
		return fail(err)
	}
	_, err = tx.Exec(sqlquery.InsertSemestersByYear, userID, newYearID)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

// TODO: this must remove all courses from the year but add them to unassigned
func (m DBManager) RemoveYear(userID string) error {
	fail := func(err error) error {
		return fmt.Errorf("RemoveYear: %v", err)
	}
	_, err := m.DB.Exec(sqlquery.DeleteYear, userID)
	if err != nil {
		return fail(err)
	}

	return nil
}

func (m DBManager) FoldSemester(userID string, year int, semester SemesterAssignment, folded bool) error {
	_, err := m.DB.Exec(sqlquery.FoldSemester, userID, year, int(semester), folded)
	if err != nil {
		return err
	}
	return nil
}
