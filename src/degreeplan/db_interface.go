package degreeplan

import (
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/language"
)

type DataManager interface {
	DegreePlan(userID string, lang language.Language) (*DegreePlan, error)
	Course(userID, courseCode string, lang language.Language) (Course, error)
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

func (bloc *Bloc) hasLimit() bool {
	return bloc.Limit > -1
}

func (b *Bloc) inBlueprint() bool {
	if b.hasLimit() && b.inBlueprintCredits() >= b.Limit {
		return true
	}
	return false
}

func (b *Bloc) inBlueprintCredits() int {
	credits := 0
	for _, c := range b.Courses {
		for _, year := range c.BlueprintYears {
			if year > 0 {
				credits += c.Credits
				break
			}
		}
	}
	return credits
}

func (b *Bloc) completed() bool {
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

type CourseStatus string

func (c CourseStatus) String() string {
	return string(c)
}

func (c *CourseStatus) UnmarshalJSON(b []byte) error {
	var r string
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	*c = CourseStatus(r)
	return nil
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

type Course struct {
	Code           string
	Title          string
	Note           string
	Credits        int
	Start          TeachingSemester
	Guarantors     []Teacher
	LectureRange1  int
	SeminarRange1  int
	LectureRange2  int
	SeminarRange2  int
	SemesterCount  int
	ExamType       string
	BlueprintYears []int64
}

func (c *Course) assignedInBlueprint() bool {
	for _, year := range c.BlueprintYears {
		if year > 0 {
			return true
		}
	}
	return false
}
