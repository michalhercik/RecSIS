package blueprint

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DataManager interface {
	Blueprint(userID string, lang DBLang) (*Blueprint, error)
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
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
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
