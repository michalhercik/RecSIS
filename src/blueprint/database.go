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

func (m DBManager) Blueprint(sessionID string, lang DBLang) (*Blueprint, error) {
	var records []BlueprintRecord
	if err := m.DB.Select(&records, sqlquery.SelectCourses, sessionID, lang); err != nil {
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

func (m DBManager) NewCourse(sessionID string, course string, year int, semester SemesterAssignment) (int, error) {
	row := m.DB.QueryRow(sqlquery.InsertCourse, sessionID, year, int(semester), course)
	var courseID int
	err := row.Scan(&courseID)
	return courseID, err
}

func (m DBManager) InsertCourses(sessionID string, year int, semester SemesterAssignment, position int, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.MoveCourses, sessionID, pq.Array(courses), year, semester, position)
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

func (m DBManager) AppendCourses(sessionID string, year int, semester SemesterAssignment, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.AppendCourses, sessionID, year, int(semester), pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) UnassignYear(sessionID string, year int) error {
	_, err := m.DB.Exec(sqlquery.UnassignYear, sessionID, year)
	return err
}

func (m DBManager) UnassignSemester(sessionID string, year int, semester SemesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.UnassignSemester, sessionID, year, int(semester))
	return err
}

func (m DBManager) RemoveCourses(sessionID string, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesByID, sessionID, pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) RemoveCoursesBySemester(sessionID string, year int, semester SemesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesBySemester, sessionID, year, int(semester))
	return err
}

func (m DBManager) RemoveCoursesByYear(sessionID string, year int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesByYear, sessionID, year)
	return err
}

func (m DBManager) AddYear(sessionID string) error {
	fail := func(err error) error {
		return fmt.Errorf("AddYear: %v", err)
	}
	tx, err := m.DB.Beginx()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var newYearID int
	err = tx.Get(&newYearID, sqlquery.InsertYear, sessionID)
	if err != nil {
		return fail(err)
	}
	_, err = tx.Exec(sqlquery.InsertSemestersByYear, sessionID, newYearID)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

// TODO: this must remove all courses from the year but add them to unassigned
func (m DBManager) RemoveYear(sessionID string) error {
	fail := func(err error) error {
		return fmt.Errorf("RemoveYear: %v", err)
	}
	_, err := m.DB.Exec(sqlquery.DeleteYear, sessionID)
	if err != nil {
		return fail(err)
	}

	return nil
}
