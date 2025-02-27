package degreeplan

import (
	"encoding/json"
)

type DataManager interface {
	// TODO: add lang parameter
	DegreePlan(uid string, lang DBLang) (*DegreePlan, error)
}

type DBLang int

const (
	cs DBLang = iota
	en
)

var DBLangs = map[DBLang]string{
	cs: "cs",
	en: "en",
}

func (l DBLang) String() string {
	return DBLangs[l]
}

type DegreePlan struct {
	blocs []Bloc
}

type Bloc struct {
	Name    string   `json:"BLOC_NAME"`
	Note    string   `json:"BLOC_NOTE"`
	Limit   int      `json:"BLOC_LIMIT"`
	Courses []Course `json:"COURSES"`
}

type BlueprintStatus int

const (
	NotInBlueprint BlueprintStatus = iota
	InBlueprint
)

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
	Code            string 			 `json:"CODE"`
	Name            string 			 `json:"NAME"`
	Note            string 			 `json:"NOTE"`
	Credits         int    			 `json:"CREDITS"`
	Start           TeachingSemester `json:"SEMESTER_PRIMARY"`
	LectureRange1   int    			 `json:"WORKLOAD_PRIMARY1"`
	SeminarRange1   int    			 `json:"WORKLOAD_SECONDARY1"`
	LectureRange2   int    			 `json:"WORKLOAD_PRIMARY2"`
	SeminarRange2   int    			 `json:"WORKLOAD_SECONDARY2"`
	SemesterCount   int    			 `json:"SEMESTER_COUNT"`
	ExamType 	    string 			 `json:"SUBJECT_TYPE"`
	BlueprintStatus BlueprintStatus
}
