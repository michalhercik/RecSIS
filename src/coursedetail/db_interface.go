package coursedetail

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

const (
	positiveRating = 1
	negativeRating = 0

	hideModal = false
	showModal = true
)

type DBLang string

// const (
// 	cs DBLang = "cs"
// 	en DBLang = "en"
// )

type Faculty struct {
	SisID int
	Name  string
	Abbr  string
}

type Semester int

const (
	winter Semester = iota + 1
	summer
	both
)

func (s Semester) String(lang string) string {
	l := language.Language(lang)
	switch s {
	case winter:
		return texts[l].Winter
	case summer:
		return texts[l].Summer
	case both:
		return texts[l].Both
	default:
		return "unknown"
	}
}

type Teacher struct {
	SisID       string `json:"KOD"`
	FirstName   string `json:"JMENO"`
	LastName    string `json:"PRIJMENI"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
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

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts *TeachingSemester) String(lang string) string {
	semester := ""
	l := language.Language(lang)
	switch *ts {
	case teachingWinterOnly:
		semester = texts[l].Winter
	case teachingSummerOnly:
		semester = texts[l].Summer
	case teachingBoth:
		semester = texts[l].Both
	default:
		semester = "unsupported"
	}
	return semester
}

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

type Description struct {
	Title   string `json:"TITLE"`
	Content string `json:"MEMO"`
}

func (d *Description) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &d)
		return nil
	case string:
		json.Unmarshal([]byte(v), &d)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (d Description) Value() (interface{}, error) {
	return json.Marshal(d)
}

type NullDescription struct {
	Description
	Valid bool
}

func (d *NullDescription) Scan(val interface{}) error {
	if val == nil {
		d.Valid = false
		return nil
	}
	if err := d.Description.Scan(val); err != nil {
		return err
	}
	d.Valid = true
	return nil
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

func (n *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	err := i.Scan(value)
	if err != nil {
		return err
	}
	*n = NullInt64(i)
	return nil
}

func (n NullInt64) String() string {
	if !n.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%d", n.Int64)
}

type NullFloat64 sql.NullFloat64

func (n *NullFloat64) Scan(value interface{}) error {
	var i sql.NullFloat64
	err := i.Scan(value)
	if err != nil {
		return err
	}
	*n = NullFloat64(i)
	return nil
}

func (n NullFloat64) String() string {
	if !n.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%f", n.Float64)
}

type CourseRating struct {
	UserRating  NullInt64   `db:"rating"`
	AvgRating   NullFloat64 `db:"avg_rating"`
	RatingCount NullInt64   `db:"rating_count"`
}

type CourseCategoryRating struct {
	Code  int    `db:"category_code"`
	Title string `db:"rating_title"`
	CourseRating
}

type CourseInfo struct {
	Code                  string           `db:"code"`
	Name                  string           `db:"title"`
	Faculty               string           `db:"faculty"`
	GuarantorDepartment   string           `db:"guarantor"`
	State                 string           `db:"taught"`
	Start                 TeachingSemester `db:"start_semester"`
	SemesterCount         int              `db:"semester_count"`
	Language              string           `db:"taught_lang"`
	LectureRange1         int              `db:"lecture_range1"`
	SeminarRange1         int              `db:"seminar_range1"`
	LectureRange2         int              `db:"lecture_range2"`
	SeminarRange2         int              `db:"seminar_range2"`
	ExamType              string           `db:"exam_type"`
	Credits               int              `db:"credits"`
	Guarantors            TeacherSlice     `db:"guarantors"`
	Teachers              TeacherSlice     `db:"teachers"`
	MinEnrollment         Capacity         `db:"min_number"`
	Capacity              string           `db:"capacity"`
	Annotation            NullDescription  `db:"annotation"`
	Syllabus              NullDescription  `db:"syllabus"`
	PassingTerms          NullDescription  `db:"terms_of_passing"`
	Literature            NullDescription  `db:"literature"`
	AssesmentRequirements NullDescription  `db:"requirements_for_assesment"`
	EntryRequirements     NullDescription  `db:"entry_requirements"`
	Aim                   NullDescription  `db:"aim"`
	Prereq                JSONStringArray  `db:"preqrequisities"`
	Coreq                 JSONStringArray  `db:"corequisities"`
	Incompa               JSONStringArray  `db:"incompatibilities"`
	Interchange           JSONStringArray  `db:"interchangebilities"`
	Classes               ClassSlice       `db:"classes"`
	Classifications       ClassSlice       `db:"classifications"`
}

type ClassSlice []Class

func (cs *ClassSlice) Scan(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*cs = nil
		return nil
	case []byte:
		*cs = nil
		err := json.Unmarshal(v, &cs)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &cs)
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type Class struct {
	Code string `json:"KOD"`
	Name string `json:"NAZEV"`
}

type JSONStringArray []string

func (jsa *JSONStringArray) Scan(val interface{}) error {
	switch v := val.(type) {
	case nil:
		jsa = nil
		return nil
	case []byte:
		*jsa = nil
		err := json.Unmarshal(v, &jsa)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &jsa)
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (jsa JSONStringArray) String() string {
	if len(jsa) == 0 {
		return ""
	}
	return strings.Join(jsa, ", ")
}

type Course struct {
	CourseInfo
	CourseRating
	// UserOverallRating    NullInt64   `db:"overall_rating"`
	// AvgOverallRating     NullFloat64 `db:"avg_overall_rating"`
	// OverallRatingCount   NullInt64   `db:"overall_rating_count"`
	Link                 string // link to course webpage (not SIS)
	BlueprintAssignments []Assignment
	CategoryRatings      []CourseCategoryRating
	Comments             search.SearchResult //[]course.Comment
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
