package courses

import (
	"encoding/json"
	"fmt"
	"strings"
	//"log"
)

type DataManager interface {
	Courses(query query) (coursesPage, error)
	Blueprint(user int, courses []string) (map[string][]Assignment, error)
	AddCourseToBlueprint(user int, code string) ([]Assignment, error)
	RemoveCourseFromBlueprint(user int, code string) error
}

type SemesterAssignment int

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

type coursesPage struct {
	courses    []Course
	startIndex int
	count      int
	total      int
	search     string
	sorted     sortType
}

type query struct {
	user       int
	startIndex int
	maxCount   int
	search     string
	sorted     sortType
}

const (
	relevance = iota
	recommended
	rating
	mostPopular
	newest
)

type sortType int

var sortTypeName = map[sortType]string{
	relevance:   "Relevant",
	recommended: "Recommended",
	rating:      "By rating",
	mostPopular: "Most popular",
	newest:      "Newest",
}

func (st sortType) String() string {
	return sortTypeName[st]
}

type Teacher struct {
	sisId       int
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s %s",
		t.titleBefore, t.firstName, t.lastName, t.titleAfter)
}

type Semester int

const (
	winter Semester = iota + 1
	summer
	both
)

var semesterNameEn = map[Semester]string{
	winter: "Winter",
	summer: "Summer",
	both:   "Both",
}

func (s Semester) String() string {
	return semesterNameEn[s]
}

type Teachers []Teacher

func (t Teachers) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.String())
	}
	return strings.Join(names, "<br>")
}

func (t *Teachers) trim() {
	result := Teachers{}
	for _, teacher := range *t {
		if teacher.sisId != -1 {
			result = append(result, teacher)
		}
	}
	*t = result
}

type Assignment struct {
	year     int
	semester SemesterAssignment
}

func (a Assignment) String() string {
	result := fmt.Sprintf("Year %d, semester %d", a.year, a.semester)
	if a.year == 0 {
		result = "Not assigned"
	}
	return result
}

type Course struct {
	code                 string
	name                 string
	nameCs               string // TODO: delete after transition to name
	nameEn               string // TODO: delete after transition to name
	start                Semester
	semesterCount        int
	lectureRange1        int
	seminarRange1        int
	lectureRange2        int
	seminarRange2        int
	examType             string
	credits              int
	teachers             Teachers // TODO: delete after transition to teacher1, teacher2, teaacher3
	rating               int
	blueprintAssignments []Assignment
}

func (c *Course) UnmarshalJSON(data []byte) error {
	raw := struct {
		Code              string   `json:"code"`
		NameCS            string   `json:"nameCs"`
		NameEN            string   `json:"nameEn"`
		Start             Semester `json:"start"`
		SemesterCount     int      `json:"semesterCount"`
		LectureRange1     int      `json:"lectureRange1"`
		SeminarRange1     int      `json:"seminarRange1"`
		LectureRange2     int      `json:"lectureRange2"`
		SeminarRange2     int      `json:"seminarRange2"`
		ExamType          string   `json:"examType"`
		Credits           int      `json:"credits"`
		Teacher1Id        int      `json:"teacher1Id"`
		Teacher1Firstname string   `json:"teacher1Firstname"`
		Teacher1Lastname  string   `json:"teacher1Lastname"`
		Teacher2Id        int      `json:"teacher2Id"`
		Teacher2Firstname string   `json:"teacher2Firstname"`
		Teacher2Lastname  string   `json:"teacher2Lastname"`
		Teacher3Id        int      `json:"teacher3Id"`
		Teacher3Firstname string   `json:"teacher3Firstname"`
		Teacher3Lastname  string   `json:"teacher3Lastname"`
		Rating            int      `json:"rating"`
	}{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	c.code = raw.Code
	if raw.NameCS != "" {
		c.name = raw.NameCS
	} else {
		c.name = raw.NameEN
	}
	c.start = raw.Start
	c.semesterCount = raw.SemesterCount
	c.lectureRange1 = raw.LectureRange1
	c.seminarRange1 = raw.SeminarRange1
	c.lectureRange2 = raw.LectureRange2
	c.seminarRange2 = raw.SeminarRange2
	c.examType = raw.ExamType
	c.credits = raw.Credits
	c.teachers = Teachers{
		Teacher{
			sisId:     raw.Teacher1Id,
			firstName: raw.Teacher1Firstname,
			lastName:  raw.Teacher1Lastname,
		},
		Teacher{
			sisId:     raw.Teacher2Id,
			firstName: raw.Teacher2Firstname,
			lastName:  raw.Teacher2Lastname,
		},
		Teacher{
			sisId:     raw.Teacher3Id,
			firstName: raw.Teacher3Firstname,
			lastName:  raw.Teacher3Lastname,
		},
	}
	c.rating = raw.Rating
	return nil
}
