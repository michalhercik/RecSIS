package coursedetail

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/michalhercik/RecSIS/internal/course/comments/search"
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

const (
	positiveRating = 1
	negativeRating = 0
)

type Course struct {
	Code                string
	Name                string
	Faculty             string
	GuarantorDepartment string
	State               string
	Start               TeachingSemester
	// SemesterCount         int
	Language              string
	LectureRangeWinter    sql.NullInt64
	SeminarRangeWinter    sql.NullInt64
	LectureRangeSummer    sql.NullInt64
	SeminarRangeSummer    sql.NullInt64
	ExamType              string
	Credits               int
	Guarantors            TeacherSlice
	Teachers              TeacherSlice
	MinEnrollment         Capacity
	Capacity              string
	Annotation            NullDescription
	Syllabus              NullDescription
	PassingTerms          NullDescription
	Literature            NullDescription
	AssesmentRequirements NullDescription
	EntryRequirements     NullDescription
	Aim                   NullDescription
	Prereq                []string
	Coreq                 []string
	Incompa               []string
	Interchange           []string
	Classes               []Class
	Classifications       []Class
	CourseRating
	Link                 string // link to course webpage (not SIS)
	BlueprintAssignments []Assignment
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

type Teacher struct {
	SisID       string
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	if t.TitleBefore == "" && t.TitleAfter == "" {
		return fmt.Sprintf("%s %s", t.FirstName, t.LastName)
	}
	if t.TitleBefore == "" {
		return fmt.Sprintf("%s %s, %s",
			t.FirstName, t.LastName, t.TitleAfter)
	}
	if t.TitleAfter == "" {
		return fmt.Sprintf("%s %s %s",
			t.TitleBefore, t.FirstName, t.LastName)
	}
	return fmt.Sprintf("%s %s %s, %s",
		t.TitleBefore, t.FirstName, t.LastName, t.TitleAfter)
}

type TeacherSlice []Teacher

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

type Description struct {
	Title   string
	Content string
}

type NullDescription struct {
	Description
	Valid bool
}

// type Faculty struct {
// 	SisID int
// 	Name  string
// 	Abbr  string
// }

// type Semester int

// const (
// 	winter Semester = iota + 1
// 	summer
// 	both
// )

// func (s Semester) String(lang string) string {
// 	l := language.Language(lang)
// 	switch s {
// 	case winter:
// 		return texts[l].Winter
// 	case summer:
// 		return texts[l].Summer
// 	case both:
// 		return texts[l].Both
// 	default:
// 		return "unknown"
// 	}
// }

// func (ts *TeachingSemester) String(lang string) string {
// 	semester := ""
// 	l := language.Language(lang)
// 	switch *ts {
// 	case teachingWinterOnly:
// 		semester = texts[l].Winter
// 	case teachingSummerOnly:
// 		semester = texts[l].Summer
// 	case teachingBoth:
// 		semester = texts[l].Both
// 	default:
// 		semester = "unsupported"
// 	}
// 	return semester
// }

// func (ts *TeachingSemester) Color() string {
// 	switch *ts {
// 	case teachingWinterOnly:
// 		return "bg-winter"
// 	case teachingSummerOnly:
// 		return "bg-summer"
// 	case teachingBoth:
// 		return "bg-both"
// 	default:
// 		return "bg-text-secondary"
// 	}
// }

type Assignment struct {
	// year     int
	// semester Semester
}

func (a Assignment) String(lang string) string {
	// semester := ""
	// switch a.semester {
	// case assignmentNone:
	// 	semester = texts[lang].N
	// case assignmentWinter:
	// 	semester = texts[lang].W
	// case assignmentSummer:
	// 	semester = texts[lang].S
	// default:
	// 	semester = texts[lang].ER
	// }

	// result := fmt.Sprintf("%d%s", a.year, semester)
	// if a.year == 0 {
	// 	result = texts[lang].UN
	// }
	// return result
	return "TODO" // TODO NOT IMPLEMENTED
}

type Assignments []Assignment

func (a Assignments) String(lang string) string {
	assignments := []string{}
	for _, assignment := range a {
		assignments = append(assignments, assignment.String(lang))
	}
	if len(assignments) == 0 {
		return "TODO"
	}
	return strings.Join(assignments, " ")
}

// func (d *NullDescription) Scan(val interface{}) error {
// 	if val == nil {
// 		d.Valid = false
// 		return nil
// 	}
// 	if err := d.Description.Scan(val); err != nil {
// 		return err
// 	}
// 	d.Valid = true
// 	return nil
// }

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

// type CommentSlice []Comment

// func (cs *CommentSlice) Scan(val interface{}) error {
// 	switch v := val.(type) {
// 	case nil:
// 		*cs = nil
// 		return nil
// 	case []byte:
// 		*cs = nil
// 		err := json.Unmarshal(v, &cs)
// 		return err
// 	case string:
// 		err := json.Unmarshal([]byte(v), &cs)
// 		return err
// 	default:
// 		return fmt.Errorf("unsupported type: %T", v)
// 	}
// }

// type Comment struct {
// 	StudiesType   string  `json:"NAZEV"`
// 	StudiesYear   int     `json:"SROC"`
// 	StudiesField  string  `json:"SOBOR"`
// 	AcademicYear  int     `json:"SSKR"`
// 	TargetType    string  `json:"PRDMTYP"`
// 	TargetTeacher Teacher `json:"TEACHER"`
// 	Content       string  `json:"MEMO"`
// }

// func (c Comment) AcademicYearString() string {
// 	return strconv.Itoa(c.AcademicYear)
// }

// func (c Comment) StudiesYearString() string {
// 	return strconv.Itoa(c.StudiesYear)
// }

// func (c Comment) TargetTeacherString() string {
// 	if len(c.TargetTeacher.SisID) > 0 {
// 		return c.TargetTeacher.String()
// 	} else {
// 		return "Global"
// 	}
// }
