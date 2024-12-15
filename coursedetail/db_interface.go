package coursedetail

import "fmt"

// TODO: change interface name if interface changes
type DataManager interface {
	Course(code string) (*Course, error)
	AddComment(code string, commentContent string) error
	GetComments(code string) ([]Comment, error)
	AddCourseToBlueprint(user int, code string) ([]Assignment, error)
	RemoveCourseFromBlueprint(user int, code string) error
}

var db DataManager

func SetDataManager(newDB DataManager) {
	db = newDB
}

// TODO add more fields
type Rating struct {
	ID     int
	UserID int
	Rating int // 1..like -1..dislike
}

// TODO add more fields
type Comment struct {
	ID      int
	UserID  int
	Content string
}

type Faculty struct {
	SisID int
	Name  string
	Abbr  string
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

type Teacher struct {
	SisID       int
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s, %s",
		t.TitleBefore, t.FirstName, t.LastName, t.TitleAfter)
}

type Assignment struct {
	year     int
	semester Semester
}

func (a Assignment) String() string {
	result := fmt.Sprintf("Year %d, semester %s", a.year, a.semester)
	if a.year == 0 {
		result = "Not assigned"
	}
	return result
}

type Description struct {
	title   string
	content string
}

type Course struct {
	Code                     string
	Name                     string
	Faculty                  Faculty
	GuarantorDepartment      string
	State                    string
	Start                    Semester
	SemesterCount            int
	Language                 string
	LectureRange1            int
	SeminarRange1            int
	LectureRange2            int
	SeminarRange2            int
	ExamType                 string
	Credits                  int
	Guarantors               []Teacher
	Teachers                 []Teacher
	MinEnrollment            int // -1 means no limit
	Capacity                 int // -1 means no limit
	Annotation               Description
	CompletitionRequirements Description
	ExamRequirements         Description
	Sylabus                  Description
	Classifications          []string
	Classes                  []string
	Link                     string // link to course webpage (not SIS)
	Comments                 []Comment
	Ratings                  []Rating
	BlueprintAssignments     []Assignment
}
