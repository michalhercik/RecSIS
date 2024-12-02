package blueprint

import "fmt"

// type Database interface {
// 	GetData(user int) BlueprintData
// 	RemoveUnassigned(user int, courseCode string)
// 	RemoveYear(user int, year int)
// 	AddYear(user int)
// }

// var db Database

// func SetDatabase(newDB Database) {
// 	db = newDB
// }

// type Course struct {
// 	Id      int
// 	Code    string
// 	NameCze string
// 	NameEng string
// 	//ValidFrom            int
// 	//ValidTo              int
// 	//Faculty              string
// 	//Guarantor            string
// 	//State                string
// 	Semester string
// 	//SemesterCount        int
// 	//Language             string
// 	LectureHoursWinter int
// 	SeminarHoursWinter int
// 	LectureHoursSummer int
// 	SeminarHoursSummer int
// 	ExamWinter         string
// 	ExamSummer         string
// 	Credits            int
// 	Teachers           []string
// 	// MinEnrollment        int // -1 means no limit
// 	// Capacity             int // -1 means no limit
// }

type BlueprintData struct {
	unassigned []Course
	years      map[int][]Course
}

type DataManager interface {
	BluePrint(user int) (*Blueprint, error)
	AddCourse(user int, course string, year int, semester int, position int)
	RemoveCourse(user int, course string, year int, semester int)
	AddYear(user int)
	RemoveYear(user int, year int)
}

var db DataManager

func SetDataManager(newDb DataManager) {
	db = newDb
}

type Teacher struct {
	SisId       int
	FirstName   string
	LastName    string
	TitleBefore string
	TitleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s %s",
		t.TitleBefore, t.FirstName, t.LastName, t.TitleAfter)
}

type Semester int

const (
	Winter Semester = iota
	Summer
)

var SemesterNameEn = map[Semester]string{
	Winter: "Winter",
	Summer: "Summer",
}

func (s Semester) String() string {
	return SemesterNameEn[s]
}

type Course struct {
	Position int
	Code     string
	NameCs   string
	NameEn   string
	// ValidFrom       int
	// ValidTo         int
	// Faculty         Faculty
	// Guarantor       string
	// State           string
	Start Semester
	// SemesterCount   int
	// Language        string
	LectureRange1 int
	SeminarRange1 int
	LectureRange2 int
	SeminarRange2 int
	ExamType      string
	Credits       int
	Teachers      []Teacher
	// MinEnrollment   int // -1 means no limit
	// Capacity        int // -1 means no limit
	// AnnotationCs    string
	// AnnotationEn    string
	// SylabusCs       string
	// SylabusEn       string
	// Classifications []string
	// Classes         []string
	// Link            string // link to course webpage (not SIS)
}

type AcademicYear struct {
	position int
	winter   []Course
	summer   []Course
}

func (ay AcademicYear) winterCredits() int {
	sum := 0
	for _, course := range ay.winter {
		sum += course.Credits
	}
	return sum
}

func (ay AcademicYear) summerCredits() int {
	sum := 0
	for _, course := range ay.summer {
		sum += course.Credits
	}
	return sum
}

func (ay AcademicYear) credits() int {
	return ay.winterCredits() + ay.summerCredits()
}

type Blueprint struct {
	years []AcademicYear
}
