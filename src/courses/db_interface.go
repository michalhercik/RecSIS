package courses

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/michalhercik/RecSIS/courses/internal/filter"
	//"log"
)

type DataManager interface {
	Courses(sessionID string, courseCodes []string, lang Language) ([]Course, error)
	ParamLabels(lang Language) (map[string][]filter.ParamValue, error)
}

type SemesterAssignment int

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

type coursesPage struct {
	courses    []Course
	page       int
	pageSize   int
	totalPages int
	search     string
	sortedBy   sortType
	semester   TeachingSemester
	facets     filter.FacetDistribution
}

const (
	relevance = iota
	recommended
	rating
	mostPopular
	newest
)

type sortType int

func (st sortType) String(lang string) string {
	switch st {
	case relevance:
		return texts[lang].Relevance
	case recommended:
		return texts[lang].Recommended
	case rating:
		return texts[lang].Rating
	case mostPopular:
		return texts[lang].MostPopular
	case newest:
		return texts[lang].Newest
	default:
		return "unknown"
	}
}

type Language string

const (
	cs Language = "cs"
	en Language = "en"
)

type Teacher struct {
	SisID       int    `json:"KOD"`
	FirstName   string `json:"JMENO"`
	LastName    string `json:"PRIJMENI"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

// TODO: here should be this ↓ but DB is broken, first name always empty
// return fmt.Sprintf("%c. %s",
//
//	t.firstName[0], t.lastName)
func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s %s",
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

type Assignment struct {
	Year     int                `json:"year"`
	Semester SemesterAssignment `json:"semester"`
}

func (a Assignment) String(lang string) string {
	semester := ""
	switch a.Semester {
	case assignmentNone:
		semester = texts[lang].N
	case assignmentWinter:
		semester = texts[lang].W
	case assignmentSummer:
		semester = texts[lang].S
	default:
		semester = texts[lang].ER
	}

	result := fmt.Sprintf("%d%s", a.Year, semester)
	if a.Year == 0 {
		result = texts[lang].UN
	}
	return result
}

type AssignmentSlice []Assignment

func (a *AssignmentSlice) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*a = AssignmentSlice{}
	case []byte:
		json.Unmarshal(v, a)
	case string:
		json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (a AssignmentSlice) String(lang string) string {
	assignments := []string{}
	for _, assignment := range a {
		assignments = append(assignments, assignment.String(lang))
	}
	if len(assignments) == 0 {
		return ""
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

func (d NullDescription) string() string {
	if d.Valid {
		return d.Content
	}
	return ""
}

type Course struct {
	Code          string           `db:"code"`
	Name          string           `db:"title"`
	Annotation    NullDescription  `db:"annotation"`
	Start         TeachingSemester `db:"start_semester"`
	SemesterCount int              `db:"semester_count"`
	LectureRange1 int              `db:"lecture_range1"`
	SeminarRange1 int              `db:"seminar_range1"`
	LectureRange2 int              `db:"lecture_range2"`
	SeminarRange2 int              `db:"seminar_range2"`
	ExamType      string           `db:"exam_type"`
	Credits       int              `db:"credits"`
	Guarantors    TeacherSlice     `db:"guarantors"`
	// Rating               int
	BlueprintAssignments AssignmentSlice `db:"assignment"`
}
