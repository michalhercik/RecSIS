package coursedetail

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"sort"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/filters"
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
	Component(semesters []bool, lang language.Language, course ...string) templ.Component
	PartialComponent(lang language.Language) PartialBlueprintAdd
	Action(userID string, year int, semester int, course ...string) ([]int, error)
	ParseRequest(r *http.Request) ([]string, int, int, error)
}
type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, semesters []bool, course ...string) templ.Component

const (
	positiveRating   = 1
	negativeRating   = 0
	numberOfComments = 20
	searchQuery      = "survey-search"
	surveyOffset     = "survey-offset"
)

type SurveyViewModel struct {
	lang   language.Language
	code   string
	query  string
	survey []Comment
	offset int
	isEnd  bool
	facets iter.Seq[filters.FacetIterator] // TODO
}

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
	BlueprintSemesters   []bool
	InDegreePlan         bool
	CategoryRatings      []CourseCategoryRating
	//Comments             []Comment
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

type Comment struct {
	Student
	CommentTarget
	AcademicYear int `json:"academic_year"`
	// Semester string `json:"semester"`
	Content string `json:"content"`
}

func (c Comment) AcademicYearString() string {
	return strconv.Itoa(c.AcademicYear)
}

type CommentTarget struct {
	Type          string  `json:"target_type"` // Lecture or Seminar
	CourseCode    string  `json:"course_code"`
	TargetTeacher Teacher `json:"teacher"`
}

type Teacher struct {
	SISID       string `json:"KOD"`
	LastName    string `json:"PRIJMENI"`
	FirstName   string `json:"JMENO"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

type Student struct {
	StudyYear  int       `json:"study_year"`
	StudyField string    `json:"study_field"`
	Study      StudyType `json:"study_type"`
}

func (c Comment) StudiesYearString() string {
	return strconv.Itoa(c.StudyYear)
}

type StudyType struct {
	// Code string `json:"code"`
	// Abbr string `json:"abbr"`
	Name string `json:"name"`
}

func (st *StudyType) UnmarshalJSON(val []byte) error {
	var tmp struct {
		Code string `json:"code"`
		// Abbr   string `json:"abbr"`
		NameCs string `json:"name_cs"`
		NameEn string `json:"name_en"`
	}
	if err := json.Unmarshal(val, &tmp); err != nil {
		return err
	}
	// st.Code = tmp.Code
	// st.Abbr = tmp.Abbr
	st.Name = tmp.NameCs
	if len(st.Name) == 0 {
		st.Name = tmp.NameEn
	}
	return nil
}

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

type TeacherSlice []Teacher

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
