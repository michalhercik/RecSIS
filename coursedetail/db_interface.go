package coursedetail

import "fmt"

// TODO: change interface name if interface changes
type CourseDataManager interface {
	Course(code string) (*Course, error)
	AddComment(code string, commentContent string) error
	GetComments(code string) ([]Comment, error)
}

// TODO: rename db to something more descriptive (Mock data are not database...)
var db CourseDataManager

// TODO: rename to something more general but descriptive (doesn't have to be Database)
func SetDatabase(newDB CourseDataManager) {
	db = newDB
}

// TODO add more fields
type Rating struct {
	ID       int
	UserID   int
	Rating   int // 1..like -1..dislike
}

// TODO add more fields
type Comment struct {
	ID       int
	UserID   int
	Content  string
}

type Faculty struct {
	ID     int
	SisID  int
	NameCs string
	NameEn string
	Abbr   string
}

func (f Faculty) String() string {
	return f.NameEn
}

type Teacher struct {
	ID          int
	SisID       int
	Department  string
	Faculty     Faculty
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s, %s",
		t.TitleBefore, t.FirstName, t.LastName, t.TitleAfter)
}

type Course struct {
	ID              int
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
	Teachers        []Teacher
	MinEnrollment   int // -1 means no limit
	Capacity        int // -1 means no limit
	AnnotationCs    string
	AnnotationEn    string
	SylabusCs       string
	SylabusEn       string
	Classifications []string
	Classes         []string
	Link            string // link to course webpage (not SIS)
	Comments		[]Comment
	Ratings			[]Rating
}

func newCourse() *Course {
	return &Course{
		Faculty: Faculty{},
		Teachers: []Teacher{
			Teacher{Faculty: Faculty{}},
			Teacher{Faculty: Faculty{}},
		},
	}
}

type CourseData struct {
	course Course
	// TODO add comments and ratings
}
