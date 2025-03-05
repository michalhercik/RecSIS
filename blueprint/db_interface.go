package blueprint

import (
	"fmt"
	"strings"
)

type DataManager interface {
	BluePrint(sessionID string) (*Blueprint, error)
	NewCourse(sessionID string, course string, year int, semester SemesterAssignment) (int, error)
	AppendCourses(sessionID string, year int, semester SemesterAssignment, courses ...int) error
	InsertCourses(sessionID string, year int, semester SemesterAssignment, position int, courses ...int) error
	UnassignYear(sessionID string, year int) error
	UnassignSemester(sessionID string, year int, semester SemesterAssignment) error
	RemoveCourses(sessionID string, courses ...int) error
	RemoveCoursesBySemester(sessionID string, year int, semester SemesterAssignment) error
	RemoveCoursesByYear(sessionID string, year int) error
	AddYear(sessionID string) error
	RemoveYear(sessionID string) error
}

type Teacher struct {
	sisId       int
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%c. %s",
		t.firstName[0], t.lastName)
}

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

type SemesterAssignment int

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

type Teachers []Teacher

func (t Teachers) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.String())
	}
	if len(names) == 0 {
		return "---"
	}
	return strings.Join(names, ", ")
}

func (t *Teachers) trim() {
	result := Teachers{}
	for _, teacher := range *t {
		if teacher.sisId != -1 {
			result = append(result, teacher)
		}
	}
	*t = result
}

type Course struct {
	ID               int
	position         int
	code             string
	nameCs           string
	nameEn           string
	start            TeachingSemester
	semesterCount    int
	semesterPosition SemesterAssignment
	lectureRange1    int
	seminarRange1    int
	lectureRange2    int
	seminarRange2    int
	examType         string
	credits          int
	teachers         Teachers
}

type AcademicYear struct {
	position   int
	winter     []Course
	summer     []Course
	unassigned []Course
}

func sumCredits(courses []Course) int {
	sum := 0
	for _, course := range courses {
		sum += course.credits
	}
	return sum
}

func (ay AcademicYear) winterCredits() int {
	return sumCredits(ay.winter)
}

func (ay AcademicYear) summerCredits() int {
	return sumCredits(ay.summer)
}

func (ay AcademicYear) credits() int {
	return ay.winterCredits() + ay.summerCredits()
}

type Blueprint struct {
	years []AcademicYear
}

func (b *Blueprint) assign(year int, course *Course) error {
	if year < 0 {
		return fmt.Errorf("year must be non-negative %d", year)
	}
	if year > len(b.years) {
		return fmt.Errorf("invalid year %d > %d", year, len(b.years))
	}
	target := &b.years[year]
	switch course.semesterPosition {
	case assignmentWinter:
		target.winter = append(target.winter, *course)
	case assignmentSummer:
		target.summer = append(target.summer, *course)
	case assignmentNone:
		target.unassigned = append(target.unassigned, *course)
	default:
		return fmt.Errorf("unknown semester assignment %d", course.semesterPosition)
	}
	return nil
}
