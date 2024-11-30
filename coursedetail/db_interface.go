package coursedetail

import "fmt"

// TODO: change interface name if interface changes
type Courser interface {
	Course(code string) (*Course, error)
}

// TODO: rename db to something more descriptive (Mock data are not database...)
var db Courser

// TODO: rename to something more general but descriptive (doesn't have to be Database)
func SetDatabase(newDB Courser) {
	db = newDB
}

type Faculty struct {
	Id     int
	SisId  int
	NameCs string
	NameEn string
	Abbr   string
}

func (f Faculty) String() string {
	return f.NameCs
}

type Teacher struct {
	Id          int
	SisId       int
	Department  string
	Faculty     Faculty
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s %s",
		t.TitleBefore, t.FirstName, t.LastName, t.TitleAfter)
}

type Course struct {
	Id              int
	Code            string
	NameCs          string
	NameEn          string
	ValidFrom       int
	ValidTo         int
	Faculty         Faculty
	Guarantor       string
	State           string
	Start           int
	SemesterCount   int
	Language        string
	LectureRange1   int
	SeminarRange1   int
	LectureRange2   int
	SeminarRange2   int
	ExamType        string
	Credits         int
	Teacher1        Teacher
	Teacher2        Teacher
	MinEnrollment   int // -1 means no limit
	Capacity        int // -1 means no limit
	AnnotationCs    string
	AnnotationEn    string
	SylabusCs       string
	SylabusEn       string
	Classifications []string
	Classes         []string
	Link            string // link to course webpage (not SIS)
}

func newCourse() *Course {
	return &Course{
		Faculty: Faculty{},
		Teacher1: Teacher{
			Faculty: Faculty{},
		},
		Teacher2: Teacher{
			Faculty: Faculty{},
		},
	}
}

type CourseData struct {
	course Course
	// TODO add comments and ratings
}
