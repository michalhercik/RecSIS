package degreeplan

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	PartialComponentSecond(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request) ([]string, int, int, error)
	Action(userID string, year int, semester int, course ...string) ([]int, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string) templ.Component
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course ...string) templ.Component

type DegreePlan struct {
	blocs []Bloc
}

func (dp *DegreePlan) bpNumberOfSemesters() int {
	if len(dp.blocs) == 0 {
		return 0
	}
	if len(dp.blocs[0].Courses) == 0 {
		return 0
	}
	return len(dp.blocs[0].Courses[0].BlueprintSemesters)
}

type Bloc struct {
	Name         string
	Code         int
	Note         string
	Limit        int
	IsCompulsory bool
	Courses      []Course
}

func (b *Bloc) hasLimit() bool {
	return b.Limit > -1
}

// ignores courses unassigned in the blueprint
func (b *Bloc) isAssigned() bool {
	if b.hasLimit() && b.assignedCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *Bloc) assignedCredits() int {
	credits := 0
	for _, c := range b.Courses {
		if c.isAssigned() {
			credits += c.Credits
		}
	}
	return credits
}

func (b *Bloc) isCompleted() bool {
	if b.hasLimit() && b.completedCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *Bloc) completedCredits() int {
	credits := 0
	for _, c := range b.Courses {
		// TODO: add course completion status -> change `false` to `course.Completed`
		if false {
			credits += c.Credits
		}
	}
	return credits
}

// gets all courses in the blueprint
func (b *Bloc) isInBlueprint() bool {
	if b.hasLimit() && b.blueprintCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *Bloc) blueprintCredits() int {
	credits := 0
	for _, c := range b.Courses {
		if c.isInBlueprint() {
			credits += c.Credits
		}
	}
	return credits
}

type Course struct {
	Code               string
	Title              string
	Note               string
	Credits            int
	Start              TeachingSemester
	Guarantors         TeacherSlice
	LectureRange1      int
	SeminarRange1      int
	LectureRange2      int
	SeminarRange2      int
	SemesterCount      int
	ExamType           string
	BlueprintSemesters []bool
}

// if is (un)assigned in the blueprint
func (c *Course) isInBlueprint() bool {
	for _, isIn := range c.BlueprintSemesters {
		if isIn {
			return true
		}
	}
	return false
}

// if is assigned in the blueprint
func (c *Course) isAssigned() bool {
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
func (c *Course) isUnassigned() bool {
	if len(c.BlueprintSemesters) < 1 {
		return false
	}
	return c.BlueprintSemesters[0]
}

func (c *Course) statusBackgroundColor() string {
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

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

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
