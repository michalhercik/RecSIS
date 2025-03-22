package degreeplan

import (
	"encoding/json"
)

type DataManager interface {
	// TODO: add lang parameter
	DegreePlan(uid string, lang DBLang) (*DegreePlan, error)
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
