package blueprint

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/blueprint/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func scan(rows *sql.Rows, year *int, course *Course) error {
	err := rows.Scan(
		year,
		&course.ID,
		&course.position,
		&course.semesterPosition,
		&course.code,
		&course.nameCs,
		&course.nameEn,
		&course.start,
		&course.semesterCount,
		&course.lectureRange1,
		&course.lectureRange2,
		&course.seminarRange1,
		&course.seminarRange2,
		&course.examType,
		&course.credits,
		&course.teachers[0].sisId,
		&course.teachers[0].firstName,
		&course.teachers[0].lastName,
		&course.teachers[0].titleBefore,
		&course.teachers[0].titleAfter,
		&course.teachers[1].sisId,
		&course.teachers[1].firstName,
		&course.teachers[1].lastName,
		&course.teachers[1].titleBefore,
		&course.teachers[1].titleAfter,
		&course.teachers[2].sisId,
		&course.teachers[2].firstName,
		&course.teachers[2].lastName,
		&course.teachers[2].titleBefore,
		&course.teachers[2].titleAfter,
	)
	return err
}

func selectYears(tx *sql.Tx, user int) ([]AcademicYear, error) {
	rows, err := tx.Query(sqlquery.SelectYears, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	years := []AcademicYear{}
	for rows.Next() {
		var year AcademicYear
		if err := rows.Scan(&year.position); err != nil {
			return nil, err
		}
		years = append(years, year)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return years, nil
}

func trim(teachers []Teacher) Teachers {
	result := Teachers{}
	for _, teacher := range teachers {
		if teacher.sisId != -1 {
			result = append(result, teacher)
		}
	}
	return result
}

func selectCourses(tx *sql.Tx, user int, blueprint *Blueprint) error {
	rows, err := tx.Query(sqlquery.SelectCourses, user)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var year int
		course := &Course{
			teachers: []Teacher{{}, {}, {}},
		}
		if err := scan(rows, &year, course); err != nil {
			return err
		}
		course.teachers = trim(course.teachers)
		if err := blueprint.assign(year, course); err != nil {
			return err
		}
		if err := rows.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (m DBManager) BluePrint(user int) (*Blueprint, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	years, err := selectYears(tx, user)
	if err != nil {
		return nil, err
	}
	blueprint := &Blueprint{years: years}
	if err = selectCourses(tx, user, blueprint); err != nil {
		return blueprint, err
	}
	if err := tx.Commit(); err != nil {
		return blueprint, err
	}
	return blueprint, nil
}

func (m DBManager) NewCourse(user int, course string, year int, semester SemesterAssignment) (int, error) {
	row := m.DB.QueryRow(sqlquery.InsertCourse, user, year, int(semester), course)
	var courseID int
	err := row.Scan(&courseID)
	return courseID, err
}

func (m DBManager) InsertCourses(user int, year int, semester SemesterAssignment, position int, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.MoveCourses, semester, position, user, year, pq.Array(courses))
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

func (m DBManager) AppendCourses(user int, year int, semester SemesterAssignment, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.AppendCourses, user, year, int(semester), pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) UnassignYear(user int, year int) error {
	_, err := m.DB.Exec(sqlquery.UnassignYear, user, year)
	return err
}

func (m DBManager) UnassignSemester(user int, year int, semester SemesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.UnassignSemester, user, year, int(semester))
	return err
}

func (m DBManager) RemoveCourses(user int, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCourses, user, pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) RemoveCoursesBySemester(user, year int, semester SemesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesBySemester, user, year, int(semester))
	return err
}

func (m DBManager) RemoveCoursesByYear(user, year int) error {
	_, err := m.DB.Exec(sqlquery.DeleteYearCourses, user, year)
	return err
}

func (m DBManager) AddYear(user int) error {
	res, err := m.DB.Exec(sqlquery.InsertYear, user)
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

// TODO: this must remove all courses from the year but add them to unassigned
func (m DBManager) RemoveYear(user int) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(sqlquery.DeleteYearCourses, user)
	if err != nil {
		return err
	}
	res, err := tx.Exec(sqlquery.DeleteYear, user)
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

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
