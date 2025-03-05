package coursedetail

import "fmt"

// TODO: change interface name if interface changes
type DataManager interface {
	Course(code string, lang DBLang) (*Course, error)
	AddComment(code string, commentContent string) error
	GetComments(code string) ([]Comment, error)
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

func (s Semester) String(lang string) string {
	switch s {
    case winter:
        return texts[lang].Winter
	case summer:
		return texts[lang].Summer
	case both:
		return texts[lang].Both
    default:
		return "unknown"
    }
}

type Teacher struct {
	SisID       int
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	if t.TitleAfter == "" {
		return fmt.Sprintf("%s %s %s",
			t.TitleBefore, t.FirstName, t.LastName)
	}
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

type Capacity int

func (c Capacity) String(lang string) string {
	if c == -1 {  // -1 means no limit
		return texts[lang].CapacityNoLimit
	}
	return fmt.Sprintf("%d", c)
}

type Course struct {
	Code                     string
	// TODO here is missing the name in the other language, but it is not necessary
	Name                     string
	Faculty                  Faculty
	GuarantorDepartment      string
	// TODO this must be saved in both languages
	// TODO this is either N or V, but mostly just V, should be Taught, Not .. and respectively in czech
	State                    string
	Start                    Semester
	SemesterCount            int
	// TODO in some cases is both CZ and EN but here is only one
	// TODO it is string and it ignores languages (is in english only)
	Language                 string
	LectureRange1            int
	SeminarRange1            int
	LectureRange2            int
	SeminarRange2            int
	// TODO this must be saved in both languages, and as Z+Zk, ... not as *
	ExamType                 string
	Credits                  int
	// TODO there are some empty guarantors, and those are rendering as empty strings
	Guarantors               []Teacher
	Teachers                 []Teacher
	MinEnrollment            Capacity
	Capacity                 Capacity
	Annotation               Description
	// TODO this is Cil predmetu, is it ok?
	CompletionRequirements   Description
	// TODO this is Pozadavky ke kontrole studia, is it ok?
	ExamRequirements         Description
	Sylabus                  Description
	// TODO what is this
	Classifications          []string
	// TODO what is this
	Classes                  []string
	Link                     string // link to course webpage (not SIS)
	Comments                 []Comment
	Ratings                  []Rating
	BlueprintAssignments     []Assignment
}
