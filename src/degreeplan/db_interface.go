package degreeplan

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/language"
)

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
	View(main templ.Component, lang language.Language, title string) templ.Component
}

type PartialBlueprintAdd = func(course, hxSwap, hxTarget string) templ.Component

type DegreePlan struct {
	blocs []Bloc
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

// type CourseStatus string

// func (c CourseStatus) String() string {
// 	return string(c)
// }

// func (c *CourseStatus) UnmarshalJSON(b []byte) error {
// 	var r string
// 	if err := json.Unmarshal(b, &r); err != nil {
// 		return err
// 	}
// 	*c = CourseStatus(r)
// 	return nil
// }

type Course struct {
	Code           string
	Title          string
	Note           string
	Credits        int
	Start          TeachingSemester
	Guarantors     TeacherSlice
	LectureRange1  int
	SeminarRange1  int
	LectureRange2  int
	SeminarRange2  int
	SemesterCount  int
	ExamType       string
	BlueprintYears []int64
}

// if is (un)assigned in the blueprint
func (c *Course) isInBlueprint() bool {
	return len(c.BlueprintYears) > 0
}

// if is assigned in the blueprint
func (c *Course) isAssigned() bool {
	for _, year := range c.BlueprintYears {
		if year > 0 {
			return true
		}
	}
	return false
}

// if is unassigned in the blueprint
func (c *Course) isUnassigned() bool {
	for _, year := range c.BlueprintYears {
		if year == 0 {
			return true
		}
	}
	return false
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
