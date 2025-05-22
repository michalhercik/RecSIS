package coursedetail

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/internal/course/comments/search"
	"github.com/michalhercik/RecSIS/internal/interface/teacher"
	"github.com/michalhercik/RecSIS/language"
)

// TODO: change interface name if interface changes
type DataManager interface {
	Course(sessionID string, code string, lang language.Language) (*Course, error)
	RateCategory(sessionID string, code string, category string, rating int, lang language.Language) ([]CourseCategoryRating, error)
	DeleteCategoryRating(sessionID string, code string, category string, lang language.Language) ([]CourseCategoryRating, error)
	Rate(sessionID string, code string, value int) (CourseRating, error)
	DeleteRating(sessionID string, code string) (CourseRating, error)
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Page interface {
	View(templ.Component, language.Language, string) templ.Component
}

type BlueprintAddButton interface {
	Component(course string, numberOfYears int, lang language.Language) templ.Component
	PartialComponent(numberOfYears int, lang language.Language) PartialBlueprintAdd
	NumberOfYears(userID string) (int, error)
	Action(userID, course string, year int, semester dbds.SemesterAssignment) (int, error)
}
type PartialBlueprintAdd = func(course, hxSwap, hxTarget string) templ.Component

const (
	positiveRating   = 1
	negativeRating   = 0
	numberOfComments = 20
	nOCommentsQuery  = "number-of-comments"
)

type Course struct {
	Code                   string
	Name                   string
	Faculty                string
	GuarantorDepartment    string
	State                  string
	Start                  TeachingSemester
	Language               string
	LectureRangeWinter     sql.NullInt64
	SeminarRangeWinter     sql.NullInt64
	LectureRangeSummer     sql.NullInt64
	SeminarRangeSummer     sql.NullInt64
	ExamType               string
	Credits                int
	Guarantors             TeacherSlice
	Teachers               TeacherSlice
	MinEnrollment          Capacity
	Capacity               string
	Annotation             NullDescription
	Syllabus               NullDescription
	PassingTerms           NullDescription
	Literature             NullDescription
	AssessmentRequirements NullDescription
	EntryRequirements      NullDescription
	Aim                    NullDescription
	Prereq                 []string
	Coreq                  []string
	Incompa                []string
	Interchange            []string
	Classes                []Class
	Classifications        []Class
	CourseRating
	Link                 string // link to course webpage (not SIS)
	BlueprintAssignments AssignmentSlice
	InDegreePlan         bool
	CategoryRatings      []CourseCategoryRating
	Comments             search.SearchResult
}

func (c Course) IsTaughtInWinter() bool {
	return c.LectureRangeWinter.Valid && c.SeminarRangeWinter.Valid
}

func (c Course) IsTaughtInSummer() bool {
	return c.LectureRangeSummer.Valid && c.SeminarRangeSummer.Valid
}

func (c Course) IsTaughtBoth() bool {
	return c.IsTaughtInWinter() && c.IsTaughtInSummer()
}

func (c Course) courseStyleClass() string {
	switch c.Start {
	case teachingBoth:
		return "bg-both"
	case teachingWinterOnly:
		return "bg-winter"
	case teachingSummerOnly:
		return "bg-summer"
	default:
		return ""
	}
}

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

type TeacherSlice []teacher.Teacher

type Description struct {
	Title   string
	Content string
}

type NullDescription struct {
	Description
	Valid bool
}

type SemesterAssignment int

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

func (sa SemesterAssignment) IDstring() string {
	switch sa {
	case assignmentNone:
		return "none"
	case assignmentWinter:
		return "winter"
	case assignmentSummer:
		return "summer"
	default:
		return "unsupported"
	}
}

type Assignment struct {
	id       int
	year     int
	semester SemesterAssignment
}

func (a Assignment) String(lang language.Language) string {
	t := texts[lang]
	semester := ""
	switch a.semester {
	case assignmentNone:
		semester = "unsupported"
	case assignmentWinter:
		semester = t.WinterAssign
	case assignmentSummer:
		semester = t.SummerAssign
	default:
		semester = "unsupported"
	}

	result := fmt.Sprintf("%s %s", t.YearStr(a.year), semester)
	if a.year == 0 {
		result = t.Unassigned
	}
	return result
}

type AssignmentSlice []Assignment

func (a AssignmentSlice) Sort() AssignmentSlice {
	sort.Slice(a, func(i, j int) bool {
		if a[i].year == a[j].year {
			return a[i].semester < a[j].semester
		}
		return a[i].year < a[j].year
	})
	return a
}

type Capacity int

func (c Capacity) String(lang string) string {
	l := language.Language(lang)
	if c == -1 { // -1 means no limit
		return texts[l].CapacityNoLimit
	}
	return fmt.Sprintf("%d", c)
}

type NullInt64 sql.NullInt64

func (n NullInt64) String() string {
	if !n.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%d", n.Int64)
}

type NullFloat64 sql.NullFloat64

func (n NullFloat64) String() string {
	if !n.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%f", n.Float64)
}

type CourseRating struct {
	UserRating  NullInt64
	AvgRating   NullFloat64
	RatingCount NullInt64
}

type CourseCategoryRating struct {
	Code  int
	Title string
	CourseRating
}

type Class struct {
	Code string
	Name string
}
