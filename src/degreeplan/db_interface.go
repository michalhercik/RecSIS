package degreeplan

import (
	"encoding/json"
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

type DataManager interface {
	DegreePlan(uid string, lang language.Language) (*DegreePlan, error)
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type DBLang string

const (
	cs DBLang = "cs"
	en DBLang = "en"
)

type DegreePlan struct {
	blocs []Bloc
}

func (dp *DegreePlan) add(record DegreePlanRecord) {
	blocIndex := -1
	for i, b := range dp.blocs {
		if b.Code == record.BlocCode {
			blocIndex = i
			break
		}
	}
	if blocIndex == -1 {
		dp.blocs = append(dp.blocs, Bloc{
			Name:  record.BlocName,
			Code:  record.BlocCode,
			Note:  record.BlocNote,
			Limit: record.BlocLimit,
		})
		blocIndex = len(dp.blocs) - 1
	}
	dp.blocs[blocIndex].Courses = append(dp.blocs[blocIndex].Courses, record.Course)
}

type Bloc struct {
	Name    string
	Code    int
	Note    string
	Limit   int
	Courses []Course
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
		if c.InBlueprint {
			credits += c.Credits
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

type Course struct {
	Code          string           `db:"code"`
	Title         string           `db:"title"`
	Note          string           `db:"note"`
	Credits       int              `db:"credits"`
	Start         TeachingSemester `db:"start_semester"`
	LectureRange1 int              `db:"lecture_range1"`
	SeminarRange1 int              `db:"seminar_range1"`
	LectureRange2 int              `db:"lecture_range2"`
	SeminarRange2 int              `db:"seminar_range2"`
	SemesterCount int              `db:"semester_count"`
	ExamType      string           `db:"exam_type"`
	InBlueprint   bool             `db:"in_blueprint"`
}

type DegreePlanRecord struct {
	BlocCode  int    `db:"bloc_subject_code"`
	BlocLimit int    `db:"bloc_limit"`
	BlocName  string `db:"bloc_name"`
	BlocNote  string `db:"bloc_note"`
	Note      string `db:"note"`
	Course
}
