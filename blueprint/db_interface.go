package blueprint

import (
	"fmt"
	"strings"
)

type DataManager interface {
	BluePrint(user int) (*Blueprint, error)
	MoveCourse(user int, course string, year int, semester int, position int) error
	RemoveCourse(user int, course string, year int, semester int) error
	AddYear(user int) error
	RemoveYear(user int) error
}

var db DataManager

func SetDataManager(newDB DataManager) {
	db = newDB
}

type Teacher struct {
	sisId       int
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s %s",
		t.titleBefore, t.firstName, t.lastName, t.titleAfter)
}

type Semester int

const (
	winter Semester = iota + 1
	summer
	both
)

var semesterNameEn = map[Semester]string{
	winter: "Winter",
	summer: "Summer",
	both:   "Both",
}

func (s Semester) String() string {
	return semesterNameEn[s]
}

func (s Semester) isBoth() bool {
	switch s {
	case winter:
		return false
	case summer:
		return false
	default:
		return true
	}
}

type SemesterPosition int

const (
	none SemesterPosition = iota
	winterP
	summerP
)

type Teachers []Teacher

func (t Teachers) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.String())
	}
	return strings.Join(names, "<br>")
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
	position         int
	code             string
	nameCs           string
	nameEn           string
	start            Semester
	semesterCount    int
	semesterPosition SemesterPosition
	lectureRange1    int
	seminarRange1    int
	lectureRange2    int
	seminarRange2    int
	examType         string
	credits          int
	teachers         Teachers
}

func newCourse() *Course {
	return &Course{
		teachers: []Teacher{{}, {}},
	}
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
	case winterP:
		target.winter = append(target.winter, *course)
	case summerP:
		target.summer = append(target.summer, *course)
	case none:
		target.unassigned = append(target.unassigned, *course)
	default:
		return fmt.Errorf("unknown semester placement %d", course.semesterPosition)
	}
	return nil
}
