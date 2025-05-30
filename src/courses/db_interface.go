package courses

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
	//"log"
)

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request) ([]string, int, int, error)
	Action(userID string, year int, semester int, course ...string) ([]int, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, searchParam string) templ.Component
	SearchParam() string
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course ...string) templ.Component

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
	SisID       string
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	firstRune, _ := utf8.DecodeRuneInString(t.FirstName)
	return fmt.Sprintf("%c. %s", firstRune, t.LastName)
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
	Year     int
	Semester SemesterAssignment
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
	Title   string
	Content string
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
	Code                 string
	Name                 string
	Annotation           NullDescription
	Start                TeachingSemester
	SemesterCount        int
	LectureRange1        int
	SeminarRange1        int
	LectureRange2        int
	SeminarRange2        int
	ExamType             string
	Credits              int
	Guarantors           TeacherSlice
	BlueprintAssignments AssignmentSlice
	BlueprintSemesters   []bool
}
