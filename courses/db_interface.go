package courses

import (
	"encoding/json"
	"fmt"
	"strings"
	//"log"
)

type DataManager interface {
	//Courses(query query) (coursesPage, error)
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
	page       int
	pageSize   int
	totalPages int
	search     string
	sortedBy   sortType
	semester   TeachingSemester
}

const (
	relevance = iota
	recommended
	rating
	mostPopular
	newest
)

type sortType int


func (st sortType) String(lang string) string {
    switch st {
    case relevance:
        return texts[lang].Relevance
    case recommended:
        return texts[lang].Recommended
    case rating:
        return texts[lang].Rating
    case mostPopular:
        return texts[lang].MostPopular
    case newest:
        return texts[lang].Newest
    default:
		return "unknown"
    }
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

type Teachers []Teacher

func (t Teachers) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.String())
	}
	if len(names) == 0 {
		return "---"
	}
	return strings.Join(names, ", ")
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

type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

type Assignment struct {
	year     int
	semester SemesterAssignment
}

// TODO: this string is broken, year and semester is swapped
func (a Assignment) String(lang string) string {
	semester := ""
	switch a.semester {
	case assignmentNone:
		semester = texts[lang].N
	case assignmentWinter:
		semester = texts[lang].W
	case assignmentSummer:
		semester = texts[lang].S
	default:
		semester = texts[lang].ER
	}

	result := fmt.Sprintf("%d%s", a.year, semester)
	if a.year == 0 {
		result = texts[lang].UN
	}
	return result
}

type Assignments []Assignment

func (a Assignments) String(lang string) string {
    assignments := []string{}
    for _, assignment := range a {
        assignments = append(assignments, assignment.String(lang))
    }
    if len(assignments) == 0 {
        return ""
    }
    return strings.Join(assignments, " ")
}

type Course struct {
	code                 string
	name                 string
	annotation           string
	nameCs               string // TODO: delete after transition to name
	nameEn               string // TODO: delete after transition to name
	start                TeachingSemester
	semesterCount        int
	lectureRange1        int
	seminarRange1        int
	lectureRange2        int
	seminarRange2        int
	examType             string
	credits              int
	teachers             Teachers // TODO: delete after transition to teacher1, teacher2, teaacher3
	rating               int
	blueprintAssignments Assignments
}

func (c *Course) UnmarshalJSON(data []byte) error {
	raw := struct {
		Code              string           `json:"code"`
		NameCS            string           `json:"nameCs"`
		NameEN            string           `json:"nameEn"`
		AnnotationCs      string           `json:"annotationCs"`
		AnnotationEn      string           `json:"annotationEn"`
		Start             TeachingSemester `json:"start"`
		SemesterCount     int              `json:"semesterCount"`
		LectureRange1     int              `json:"lectureRange1"`
		SeminarRange1     int              `json:"seminarRange1"`
		LectureRange2     int              `json:"lectureRange2"`
		SeminarRange2     int              `json:"seminarRange2"`
		ExamType          string           `json:"examType"`
		Credits           int              `json:"credits"`
		Teacher1Id        int              `json:"teacher1Id"`
		Teacher1Firstname string           `json:"teacher1Firstname"`
		Teacher1Lastname  string           `json:"teacher1Lastname"`
		Teacher2Id        int              `json:"teacher2Id"`
		Teacher2Firstname string           `json:"teacher2Firstname"`
		Teacher2Lastname  string           `json:"teacher2Lastname"`
		Teacher3Id        int              `json:"teacher3Id"`
		Teacher3Firstname string           `json:"teacher3Firstname"`
		Teacher3Lastname  string           `json:"teacher3Lastname"`
		Rating            int              `json:"rating"`
	}{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	c.code = raw.Code
	if raw.NameCS != "" {
		c.name = raw.NameCS
		c.annotation = raw.AnnotationCs
	} else {
		c.name = raw.NameEN
		c.annotation = raw.AnnotationEn
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
