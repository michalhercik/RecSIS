package coursedetail

import (
	"encoding/json"
	"fmt"
)

// TODO: change interface name if interface changes
type DataManager interface {
	Course(code string, lang DBLang) (*Course, error)
	AddComment(code string, commentContent string) error
	GetComments(code string) ([]Comment, error)
	AddCourseToBlueprint(user int, code string) ([]Assignment, error)
	RemoveCourseFromBlueprint(user int, code string) error
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

var semesterNameEn = map[Semester]string{
	winter: "Winter",
	summer: "Summer",
	both:   "Both",
}

func (s Semester) String() string {
	return semesterNameEn[s]
}

type Teacher struct {
	SisID       int    `json:"KOD"`
	FirstName   string `json:"JMENO"`
	LastName    string `json:"PRIJMENI"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

func (t Teacher) String() string {
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

type Course struct {
	Code                     string       `db:"code"`
	Name                     string       `db:"title"`
	Faculty                  string       `db:"faculty"`
	GuarantorDepartment      string       `db:"guarantor"`
	State                    string       `db:"taught"`
	Start                    string       `db:"semester_description"`
	SemesterCount            int          `db:"semester_count"`
	Language                 string       `db:"taught_lang"`
	LectureRange1            int          `db:"lecture_range1"`
	SeminarRange1            int          `db:"seminar_range1"`
	LectureRange2            int          `db:"lecture_range2"`
	SeminarRange2            int          `db:"seminar_range2"`
	ExamType                 string       `db:"exam_type"`
	Credits                  int          `db:"credits"`
	Guarantors               TeacherSlice `db:"guarantors"`
	Teachers                 TeacherSlice `db:"teachers"`
	MinEnrollment            int          `db:"min_number"` // -1 means no limit
	Capacity                 string       `db:"capacity"`   // -1 means no limit
	Annotation               Description  `db:"annotation"`
	CompletitionRequirements Description  `db:"aim"`
	ExamRequirements         Description  `db:"requirements"`
	Sylabus                  Description  `db:"syllabus"`
	Classifications          []string
	Classes                  []string
	Link                     string // link to course webpage (not SIS)
	Comments                 []Comment
	Ratings                  []Rating
	BlueprintAssignments     []Assignment
}
