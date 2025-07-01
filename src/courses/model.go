package courses

import (
	"database/sql"
	"fmt"
	"iter"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Constants
//================================================================================

const (
	coursesPerPage   = 24
	courseIndex      = "courses"
	pageParam        = "page"
	hitsPerPageParam = "hitsPerPage"
)

//================================================================================
// Data Types and Methods
//================================================================================

// all data needed for the courses page
type coursesPage struct {
	courses     []course
	page        int
	pageSize    int
	totalPages  int
	search      string
	facets      iter.Seq[filters.FacetIterator]
	searchParam string
	bpBtn       PartialBlueprintAdd
}

// course representation
type course struct {
	code                 string
	title                string
	annotation           nullDescription
	semester             teachingSemester
	lectureRangeWinter   sql.NullInt64
	seminarRangeWinter   sql.NullInt64
	lectureRangeSummer   sql.NullInt64
	seminarRangeSummer   sql.NullInt64
	examType             string
	credits              int
	guarantors           teacherSlice
	blueprintAssignments assignmentSlice
	blueprintSemesters   []bool
	inDegreePlan         bool
}

// wrapper for description that allows it to be nullable
type nullDescription struct {
	description
	valid bool
}

func (d nullDescription) string() string {
	if d.valid {
		return d.content
	}
	return ""
}

// description type - title and content
type description struct {
	title   string
	content string
}

// semester type - winter, summer, or both
type teachingSemester int

const (
	teachingWinterOnly teachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts teachingSemester) string(t text) string {
	switch ts {
	case teachingWinterOnly:
		return t.winterAssign
	case teachingSummerOnly:
		return t.summerAssign
	case teachingBoth:
		return t.both
	default:
		return ""
	}
}

// wrapper for teacher slice
type teacherSlice []teacher

func (ts teacherSlice) string(t text) string {
	names := []string{}
	for _, teacher := range ts {
		names = append(names, teacher.string())
	}
	if len(names) == 0 {
		return t.noGuarantors
	}
	return strings.Join(names, ", ")
}

// teacher type
type teacher struct {
	sisID       string
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t teacher) string() string {
	firstRune, _ := utf8.DecodeRuneInString(t.firstName)
	return fmt.Sprintf("%c. %s", firstRune, t.lastName)
}

// assignment slice for blueprint assignments
type assignmentSlice []assignment

func (a assignmentSlice) Sort() assignmentSlice {
	sort.Slice(a, func(i, j int) bool {
		if a[i].year == a[j].year {
			return a[i].semester < a[j].semester
		}
		return a[i].year < a[j].year
	})
	return a
}

// assignment type - year and semester
type assignment struct {
	year     int
	semester semesterAssignment
}

func (a assignment) string(lang language.Language) string {
	t := texts[lang]
	semester := ""
	switch a.semester {
	case assignmentNone:
		semester = t.n
	case assignmentWinter:
		semester = t.w
	case assignmentSummer:
		semester = t.s
	default:
		semester = t.er
	}

	result := fmt.Sprintf("%d. %s", a.year, semester)
	if a.year == 0 {
		result = t.un
	}
	return result
}

// semester assignment type - winter or summer (none = unassigned)
type semesterAssignment int

const (
	assignmentNone semesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

func (sa semesterAssignment) stringID() string {
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

// ================================================================================
// Helper Functions
// ================================================================================

func (c *course) hoursString() string {
	result := ""
	winter := c.lectureRangeWinter.Valid && c.seminarRangeWinter.Valid
	summer := c.lectureRangeSummer.Valid && c.seminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", c.lectureRangeSummer.Int64, c.seminarRangeSummer.Int64)
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func splitByLastSpace(s string) (string, string) {
	words := strings.Fields(s)
	rest, last := "", ""
	if len(words) > 1 {
		rest = strings.Join(words[:len(words)-1], " ")
		last = " " + words[len(words)-1]
	} else {
		rest = ""
		last = s
	}
	return rest, last
}
