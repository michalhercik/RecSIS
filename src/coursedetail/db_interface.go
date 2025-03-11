package coursedetail

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// TODO: change interface name if interface changes
type DataManager interface {
	Course(sessionID string, code string, lang DBLang) (*Course, error)
	OverallRating(sessionID string, code string, value int) error
}

type DBLang string

const (
	cs DBLang = "cs"
	en DBLang = "en"
)

// TODO add more fields
type Rating struct {
	ID     int
	UserID int
	Rating int // 1..like -1..dislike
}

// TODO add more fields
type Comment struct {
	ID      int
	UserID  int
	Content string
}

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
	switch s {
	case winter:
		return texts[lang].Winter
	case summer:
		return texts[lang].Summer
	case both:
		return texts[lang].Both
	default:
		return "unknown"
	}
}

type Teacher struct {
	SisID       int    `json:"KOD"`
	FirstName   string `json:"JMENO"`
	LastName    string `json:"PRIJMENI"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

func (t Teacher) String() string {
	if t.TitleAfter == "" {
		return fmt.Sprintf("%s %s %s",
			t.TitleBefore, t.FirstName, t.LastName)
	}
	return fmt.Sprintf("%s %s %s, %s",
		t.TitleBefore, t.FirstName, t.LastName, t.TitleAfter)
}

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

type Assignment struct {
	year     int
	semester Semester
}

func (a Assignment) String() string {
	result := fmt.Sprintf("Year %d, semester %s", a.year, a.semester)
	if a.year == 0 {
		result = "Not assigned"
	}
	return result
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

type Capacity int

func (c Capacity) String(lang string) string {
	if c == -1 { // -1 means no limit
		return texts[lang].CapacityNoLimit
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

type CourseRatings struct {
	Overall    NullInt64 `db:"overall_rating"`
	Difficulty NullInt64 `db:"difficulty_rating"`
	Workload   NullInt64 `db:"workload_rating"`
	Usefulness NullInt64 `db:"usefulness_rating"`
	Fun        NullInt64 `db:"fun_rating"`
}

type Course struct {
	Code                string `db:"code"`
	Name                string `db:"title"`
	Faculty             string `db:"faculty"`
	GuarantorDepartment string `db:"guarantor"`
	State               string `db:"taught"`
	Start               string `db:"semester_description"`
	SemesterCount       int    `db:"semester_count"`
	// TODO in some cases is both CZ and EN but here is only one
	Language      string       `db:"taught_lang"`
	LectureRange1 int          `db:"lecture_range1"`
	SeminarRange1 int          `db:"seminar_range1"`
	LectureRange2 int          `db:"lecture_range2"`
	SeminarRange2 int          `db:"seminar_range2"`
	ExamType      string       `db:"exam_type"`
	Credits       int          `db:"credits"`
	Guarantors    TeacherSlice `db:"guarantors"`
	Teachers      TeacherSlice `db:"teachers"`
	MinEnrollment Capacity     `db:"min_number"`
	Capacity      string       `db:"capacity"`
	Annotation    Description  `db:"annotation"`
	// TODO this is Cil predmetu, is it ok?
	CompletionRequirements Description `db:"aim"`
	// TODO this is Pozadavky ke kontrole studia, is it ok?
	ExamRequirements Description `db:"requirements"`
	// TODO this must be syllabus - broken
	Sylabus Description `db:"syllabus"`
	// TODO what is this
	Classifications []string
	// TODO what is this
	Classes              []string
	Link                 string // link to course webpage (not SIS)
	BlueprintAssignments []Assignment
	CourseRatings
}
