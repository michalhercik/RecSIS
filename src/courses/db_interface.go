package courses

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"sort"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
	//"log"
)

type DataManager interface {
	Courses(userID string, courseCodes []string, lang language.Language) ([]Course, error)
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	Component(course string, numberOfYears int, lang language.Language) templ.Component
	PartialComponent(numberOfYears int, lang language.Language) PartialBlueprintAdd
	NumberOfYears(userID string) (int, error)
	Action(userID, course string, year int, semester dbds.SemesterAssignment) (int, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, searchParam string) templ.Component
	SearchParam() string
}

type PartialBlueprintAdd = func(course, hxSwap, hxTarget string) templ.Component

type coursesPage struct {
	courses     []Course
	page        int
	pageSize    int
	totalPages  int
	search      string
	facets      iter.Seq[filters.FacetIterator] // func(func(filter.FacetIterator) bool) //filter.Filters //FacetDistribution
	searchParam string
	templ       PartialBlueprintAdd
}

type Teacher struct {
	SisID       int    `json:"KOD"`
	FirstName   string `json:"JMENO"`
	LastName    string `json:"PRIJMENI"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

func (t Teacher) String() string {
	return fmt.Sprintf("%c. %s", t.FirstName[0], t.LastName)
}

type TeacherSlice []Teacher

func (ts TeacherSlice) string(t text) string {
	names := []string{}
	for _, teacher := range ts {
		names = append(names, teacher.String())
	}
	if len(names) == 0 {
		return t.NoGuarantors
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

func (ts *TeachingSemester) String(lang language.Language) string {
	t := texts[lang]
	semester := ""
	switch *ts {
	case teachingWinterOnly:
		semester = t.WinterAssign
	case teachingSummerOnly:
		semester = t.SummerAssign
	case teachingBoth:
		semester = t.Both
	default:
		semester = "unsupported"
	}
	return semester
}

type Assignment struct {
	Year     int                `json:"year"`
	Semester SemesterAssignment `json:"semester"`
}

func (a Assignment) String(lang language.Language) string {
	t := texts[lang]
	semester := ""
	switch a.Semester {
	case assignmentNone:
		semester = t.N
	case assignmentWinter:
		semester = t.W
	case assignmentSummer:
		semester = t.S
	default:
		semester = t.ER
	}

	result := fmt.Sprintf("%d. %s", a.Year, semester)
	if a.Year == 0 {
		result = t.UN
	}
	return result
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

func (a AssignmentSlice) Sort() AssignmentSlice {
	sort.Slice(a, func(i, j int) bool {
		if a[i].Year == a[j].Year {
			return a[i].Semester < a[j].Semester
		}
		return a[i].Year < a[j].Year
	})
	return a
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
	Code                 string           `db:"code"`
	Name                 string           `db:"title"`
	Annotation           NullDescription  `db:"annotation"`
	Start                TeachingSemester `db:"start_semester"`
	SemesterCount        int              `db:"semester_count"`
	LectureRange1        int              `db:"lecture_range1"`
	SeminarRange1        int              `db:"seminar_range1"`
	LectureRange2        int              `db:"lecture_range2"`
	SeminarRange2        int              `db:"seminar_range2"`
	ExamType             string           `db:"exam_type"`
	Credits              int              `db:"credits"`
	Guarantors           TeacherSlice     `db:"guarantors"`
	BlueprintAssignments AssignmentSlice  `db:"assignment"`
}
