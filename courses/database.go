package courses

import (
	"database/sql"
	"fmt"

	"github.com/michalhercik/RecSIS/courses/sqlquery"
)

type DBManager struct {
	DB *sql.DB
}

func scan(rows *sql.Rows, course *Course) error {
	err := rows.Scan(
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
	)
	return err
}

func (m DBManager) Courses(q query) (coursesPage, error) {
	// TODO: Take in account the search and sorted fields

	// Execute count query to get total number of courses
	var total int
	err := m.DB.QueryRow(sqlquery.CountCourses, ).Scan(&total)
	if err != nil {
		return coursesPage{}, fmt.Errorf("failed to count courses: %w", err)
	}

	// Execute main query to get courses
	rows, err := m.DB.Query(sqlquery.GetCourses, q.startIndex, q.maxCount)
	if err != nil {
		return coursesPage{}, fmt.Errorf("failed to fetch courses: %w", err)
	}
	defer rows.Close()

	// Parse rows into courses slice
	var courses []Course
	for rows.Next() {
		course := newCourse()
		if err := scan(rows, course); err != nil {
			return coursesPage{}, fmt.Errorf("failed to scan course: %w", err)
		}
		course.teachers.trim()
		// TODO: This is a temporary solution, remove it later
		course.rating = 42
		courses = append(courses, *course)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return coursesPage{}, fmt.Errorf("error during rows iteration: %w", err)
	}

	// Build the coursesPage result
	result := coursesPage{
		courses:    courses,
		startIndex: q.startIndex,
		count:      len(courses),
		total:      total,
		search:     q.search,
		sorted:     q.sorted,
	}

	return result, nil
}

func (m DBManager) AddCourseToBlueprint(user int, code string) ([]Assignment, error) {
	// TODO: Implement this method
	// year=0, semester=course.semester, position=-1
	// this must be done in a transaction, must return the all assignments
	return nil, nil
}

func (m DBManager) RemoveCourseFromBlueprint(user int, code string) error {
	// TODO: Implement this method
	// it is expected that there is only one course with the given code
	// if not return an error
	// by that we do not have to return the assignments - there should be none after the removal
	return nil
}
