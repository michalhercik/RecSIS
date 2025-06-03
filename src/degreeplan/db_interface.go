package degreeplan

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Interfaces
//================================================================================

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	PartialComponentSecond(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request) ([]string, int, int, error)
	Action(userID string, year int, semester int, course ...string) ([]int, error)
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course ...string) templ.Component

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Data Types and Methods
//================================================================================

// all data needed for the degree plan page
type degreePlanPage struct {
	blocs []bloc
}

func (dp *degreePlanPage) bpNumberOfSemesters() int {
	if len(dp.blocs) == 0 {
		return 0
	}
	if len(dp.blocs[0].Courses) == 0 {
		return 0
	}
	return len(dp.blocs[0].Courses[0].BlueprintSemesters)
}

// bloc info with courses
type bloc struct {
	Name         string
	Code         int
	Note         string
	Limit        int
	IsCompulsory bool
	Courses      []course
}

func (b *bloc) hasLimit() bool {
	return b.Limit > -1
}

func (b *bloc) isAssigned() bool {
	// ignores courses unassigned in the blueprint
	if b.hasLimit() && b.assignedCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *bloc) assignedCredits() int {
	credits := 0
	for _, c := range b.Courses {
		if c.isAssigned() {
			credits += c.Credits
		}
	}
	return credits
}

func (b *bloc) isCompleted() bool {
	if b.hasLimit() && b.completedCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *bloc) completedCredits() int {
	credits := 0
	for _, c := range b.Courses {
		// TODO: add course completion status -> change `false` to `course.Completed`
		if false {
			credits += c.Credits
		}
	}
	return credits
}

func (b *bloc) isInBlueprint() bool {
	if b.hasLimit() && b.blueprintCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *bloc) blueprintCredits() int {
	credits := 0
	for _, c := range b.Courses {
		if c.isInBlueprint() {
			credits += c.Credits
		}
	}
	return credits
}

// course representation
type course struct {
	Code               string
	Title              string
	Note               string
	Credits            int
	Semester           TeachingSemester
	Guarantors         TeacherSlice
	LectureRangeWinter sql.NullInt64
	SeminarRangeWinter sql.NullInt64
	LectureRangeSummer sql.NullInt64
	SeminarRangeSummer sql.NullInt64
	ExamType           string
	BlueprintSemesters []bool
}

// if is (un)assigned in the blueprint
func (c *course) isInBlueprint() bool {
	for _, isIn := range c.BlueprintSemesters {
		if isIn {
			return true
		}
	}
	return false
}

// if is assigned in the blueprint
func (c *course) isAssigned() bool {
	if len(c.BlueprintSemesters) < 2 {
		return false
	}
	for _, isIn := range c.BlueprintSemesters[1:] {
		if isIn {
			return true
		}
	}
	return false
}

// if is unassigned in the blueprint
func (c *course) isUnassigned() bool {
	if len(c.BlueprintSemesters) < 1 {
		return false
	}
	return c.BlueprintSemesters[0]
}

func (c *course) statusBackgroundColor() string {
	// TODO: add course completion status -> change `false` to `course.Completed`
	if false {
		// most important is completion status
		return "bg-success"
	} else if c.isAssigned() {
		// if not completed, check if assigned
		return "bg-blueprint"
	} else if c.isUnassigned() {
		// if not even assigned, check if unassigned
		return "bg-blueprint" // same color, but the row has a warning icon
	} else {
		// if nothing else, then it is not completed
		return "bg-danger"
	}
}

// semester type - winter, summer, or both
type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts TeachingSemester) string(t text) string {
	switch ts {
	case teachingBoth:
		return t.Both
	case teachingWinterOnly:
		return t.Winter
	case teachingSummerOnly:
		return t.Summer
	default:
		return ""
	}
}

// wrapper for teacher slice
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

// teacher type
type Teacher struct {
	SISID       string
	LastName    string
	FirstName   string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	if len(t.FirstName) > 0 {
		var initial rune
		for _, r := range t.FirstName {
			initial = r
			break
		}
		return fmt.Sprintf("%c. %s", initial, t.LastName)
	}
	return t.LastName
}

//================================================================================
// Helper Functions
//================================================================================

func winterString(course *course) string {
	winterText := ""
	if course.Semester == teachingWinterOnly || course.Semester == teachingBoth {
		winterText = fmt.Sprintf("%d/%d, %s", course.LectureRangeWinter.Int64, course.SeminarRangeWinter.Int64, course.ExamType)
	} else {
		winterText = "---"
	}
	return winterText
}

func summerString(course *course) string {
	summerText := ""
	if course.Semester == teachingSummerOnly {
		summerText = fmt.Sprintf("%d/%d, %s", course.LectureRangeSummer.Int64, course.SeminarRangeSummer.Int64, course.ExamType)
	} else if course.Semester == teachingBoth {
		summerText = fmt.Sprintf("%d/%d, %s", course.LectureRangeWinter.Int64, course.SeminarRangeWinter.Int64, course.ExamType)
	} else {
		summerText = "---"
	}
	return summerText
}

func hoursString(course *course, t text) string {
	result := ""
	winter := course.LectureRangeWinter.Valid && course.SeminarRangeWinter.Valid
	summer := course.LectureRangeSummer.Valid && course.SeminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", course.LectureRangeWinter.Int64, course.SeminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", course.LectureRangeSummer.Int64, course.SeminarRangeSummer.Int64)
	}
	return result
}
