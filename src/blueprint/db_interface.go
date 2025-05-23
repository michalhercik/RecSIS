package blueprint

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

type DataManager interface {
	Blueprint(userID string, lang language.Language) (*Blueprint, error)
	NewCourse(userID string, course string, year int, semester SemesterAssignment) (int, error)
	AppendCourses(userID string, year int, semester SemesterAssignment, courses ...int) error
	InsertCourses(userID string, year int, semester SemesterAssignment, position int, courses ...int) error
	UnassignYear(userID string, year int) error
	UnassignSemester(userID string, year int, semester SemesterAssignment) error
	RemoveCourses(userID string, courses ...int) error
	RemoveCoursesBySemester(userID string, year int, semester SemesterAssignment) error
	RemoveCoursesByYear(userID string, year int) error
	AddYear(userID string) error
	RemoveYear(userID string) error
	FoldSemester(userID string, year int, semester SemesterAssignment, folded bool) error
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string) templ.Component
}

const (
	yearUnassign     string = "year-unassign"
	semesterUnassign string = "semester-unassign"
	selectedMove     string = "selected"
)

const (
	yearRemove     string = "year"
	semesterRemove string = "semester"
	selectedRemove string = "selected"
)

const (
	winterStr     = "winter"
	summerStr     = "summer"
	unassignedStr = "unassigned"
)

// type DBLang string

// const (
// 	DBLangCS DBLang = "cs"
// 	DBLangEN DBLang = "en"
// )

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

type TeacherSlice []Teacher

func (ts *TeacherSlice) Scan(val interface{}) error {
	switch v := val.(type) {
	case nil:
		return nil
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

func SemesterAssignmentFromString(s string) (SemesterAssignment, error) {
	switch s {
	case winterStr:
		return assignmentWinter, nil
	case summerStr:
		return assignmentSummer, nil
	case unassignedStr:
		return assignmentNone, nil
	default:
		return 0, fmt.Errorf("unknown semester assignment %s", s)
	}
}

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

type NullCourse struct {
	ID            sql.NullInt32  `db:"id"`
	Code          sql.NullString `db:"code"`
	Title         sql.NullString `db:"title"`
	Start         sql.NullInt32  `db:"start_semester"`
	SemesterCount sql.NullInt32  `db:"semester_count"`
	LectureRange1 sql.NullInt32  `db:"lecture_range1"`
	SeminarRange1 sql.NullInt32  `db:"seminar_range1"`
	LectureRange2 sql.NullInt32  `db:"lecture_range2"`
	SeminarRange2 sql.NullInt32  `db:"seminar_range2"`
	ExamType      sql.NullString `db:"exam_type"`
	Credits       sql.NullInt32  `db:"credits"`
	Guarantors    TeacherSlice   `db:"guarantors"`
}

func (nc NullCourse) Course() *Course {
	if !nc.ID.Valid {
		return nil
	}
	return &Course{
		ID:            int(nc.ID.Int32),
		Code:          nc.Code.String,
		Title:         nc.Title.String,
		Start:         TeachingSemester(nc.Start.Int32),
		SemesterCount: int(nc.SemesterCount.Int32),
		LectureRange1: int(nc.LectureRange1.Int32),
		SeminarRange1: int(nc.SeminarRange1.Int32),
		LectureRange2: int(nc.LectureRange2.Int32),
		SeminarRange2: int(nc.SeminarRange2.Int32),
		ExamType:      nc.ExamType.String,
		Credits:       int(nc.Credits.Int32),
		Guarantors:    nc.Guarantors,
	}
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
	Warnings      []string
}

type Semester struct {
	courses []Course
	folded  bool
}

type AcademicYear struct {
	position int
	winter   Semester
	summer   Semester
	// winter     []Course
	// summer     []Course
	unassigned Semester
}

func sumCredits(courses []Course) int {
	sum := 0
	for _, course := range courses {
		sum += course.Credits
	}
	return sum
}

func (ay AcademicYear) winterCredits() int {
	return sumCredits(ay.winter.courses)
}

func (ay AcademicYear) summerCredits() int {
	return sumCredits(ay.summer.courses)
}

func (ay AcademicYear) credits() int {
	return ay.winterCredits() + ay.summerCredits()
}

// type insertedCourseInfo struct {
// 	courseID     int
// 	academicYear int
// 	semester     SemesterAssignment
// }

type BlueprintRecordPosition struct {
	AcademicYear int                `db:"academic_year"`
	Semester     SemesterAssignment `db:"semester"`
	Folded       bool               `db:"folded"`
}

type BlueprintRecord struct {
	BlueprintRecordPosition
	NullCourse
}

type Blueprint struct {
	years []AcademicYear
}

func (b *Blueprint) totalCredits() int {
	total := 0
	for _, year := range b.years {
		total += year.credits()
	}
	return total
}

func (b *Blueprint) assign(position BlueprintRecordPosition, course *Course) error {
	if position.AcademicYear < 0 {
		return fmt.Errorf("year must be non-negative %d", position.AcademicYear)
	}
	for position.AcademicYear >= len(b.years) {
		b.years = append(b.years, AcademicYear{position: position.AcademicYear})
	}
	if course != nil {
		target := &b.years[position.AcademicYear]
		var semester *Semester
		switch position.Semester {
		case assignmentWinter:
			semester = &target.winter
			// target.winter.courses = append(target.winter.courses, *course)
		case assignmentSummer:
			semester = &target.summer
			// target.summer.courses = append(target.summer.courses, *course)
		case assignmentNone:
			semester = &target.unassigned
			// target.unassigned = append(target.unassigned, *course)
		default:
			return fmt.Errorf("unknown semester assignment %d", position.Semester)
		}
		if len(semester.courses) == 0 {
			semester.folded = position.Folded
		}
		semester.courses = append(semester.courses, *course)
	}
	return nil
}

// type courseAdditionRequestSource int

// const (
// 	sourceNone courseAdditionRequestSource = iota
// 	sourceBlueprint
// 	sourceCourseDetail
// 	sourceDegreePlan
// )

// func (r courseAdditionRequestSource) String() string {
// 	switch r {
// 	case sourceBlueprint:
// 		return "blueprint"
// 	case sourceCourseDetail:
// 		return "courseDetail"
// 	case sourceDegreePlan:
// 		return "degreePlan"
// 	default:
// 		return fmt.Sprintf("unknown %d", r)
// 	}
// }
