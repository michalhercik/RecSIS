package blueprint

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DataManager interface {
	Blueprint(sessionID string, lang DBLang) (*Blueprint, error)
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

type DBLang string

const (
	DBLangCS DBLang = "cs"
	DBLangEN DBLang = "en"
)

type Teacher struct {
	SisId       int    `json:"KOD"`
	FirstName   string `json:"JMENO"`
	LastName    string `json:"PRIJMENI"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

func (t Teacher) String() string {
	var result string
	if len(t.FirstName) > 0 {
		var initial rune
		for _, r := range t.FirstName {
			initial = r
			break
		}
		result = fmt.Sprintf("%c. %s", initial, t.LastName)
	}
	return result
}

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

type SemesterAssignment int

func (sa *SemesterAssignment) Scan(val interface{}) error {
	switch v := val.(type) {
	case int64:
		*sa = SemesterAssignment(v)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

type TeacherSlice []Teacher

func (ts *TeacherSlice) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &ts)
		return nil
	case string:
		json.Unmarshal([]byte(v), &ts)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (t TeacherSlice) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.String())
	}
	if len(names) == 0 {
		return "---"
	}
	return strings.Join(names, ", ")
}

type Course struct {
	ID            int              `db:"id"`
	Code          string           `db:"code"`
	Title         string           `db:"title"`
	Start         TeachingSemester `db:"start_semester"`
	SemesterCount int              `db:"semester_count"`
	LectureRange1 int              `db:"lecture_range1"`
	SeminarRange1 int              `db:"seminar_range1"`
	LectureRange2 int              `db:"lecture_range2"`
	SeminarRange2 int              `db:"seminar_range2"`
	ExamType      string           `db:"exam_type"`
	Credits       int              `db:"credits"`
	Guarantors    TeacherSlice     `db:"guarantors"`
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
		sum += course.Credits
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

type insertedCourseInfo struct {
	courseID     int
	academicYear int
	semester     SemesterAssignment
}

type BlueprintRecordPosition struct {
	AcademicYear int                `db:"academic_year"`
	Semester     SemesterAssignment `db:"semester"`
}

type BlueprintRecord struct {
	BlueprintRecordPosition
	Course
}

type Blueprint struct {
	years []AcademicYear
}

func (b *Blueprint) assign(position BlueprintRecordPosition, course *Course) error {
	if position.AcademicYear < 0 {
		return fmt.Errorf("year must be non-negative %d", position.AcademicYear)
	}
	for position.AcademicYear >= len(b.years) {
		b.years = append(b.years, AcademicYear{position: position.AcademicYear})
	}
	target := &b.years[position.AcademicYear]
	switch position.Semester {
	case assignmentWinter:
		target.winter = append(target.winter, *course)
	case assignmentSummer:
		target.summer = append(target.summer, *course)
	case assignmentNone:
		target.unassigned = append(target.unassigned, *course)
	default:
		return fmt.Errorf("unknown semester assignment %d", position.Semester)
	}
	return nil
}

type courseAdditionRequestSource int

const (
	sourceNone courseAdditionRequestSource = iota
	sourceBlueprint
	sourceCourseDetail
	sourceDegreePlan
)

func (r courseAdditionRequestSource) String() string {
	switch r {
	case sourceBlueprint:
		return "blueprint"
	case sourceCourseDetail:
		return "courseDetail"
	case sourceDegreePlan:
		return "degreePlan"
	default:
		return fmt.Sprintf("unknown %d", r)
	}
}
