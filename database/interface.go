package database

const User = 42

// DB defines the methods for interacting with the database
type DB interface {
	// courses
	GetAllCourses() []Course
	GetCourse(id int) (Course, error)
	// blueprint
	BlueprintGetUnassigned(user int) []Course
	BlueprintGetAssigned(user int) map[int][]Course
	BlueprintAddUnassigned(user int, course int)
	BlueprintRemoveUnassigned(user int, course int)
	BlueprintAddYear(user int)
	BlueprintRemoveYear(user int, year int) []Course
	//BlueprintAddCourseToYear(user int, year int, course int)
	//BlueprintRemoveCourseFromYear(user int, year int, course int)
	// TODO add more
}

type BlueprintData struct {
	Unassigned []Course
	Years map[int][]Course
}

type CoursesData struct {
	Courses []Course
}

type CourseData struct {
	Course Course
	// TODO add comments and ratings
}

// Course representation in system
// TODO:
//	- check validity of field types (consider enum types)
//  - consider removing irrelevant or adding missing fields
type Course struct {
	Id                   int
	Code                 string
	NameCze              string
	NameEng              string
	ValidFrom            int
	ValidTo              int
	Faculty              int
	Guarantor            string
	State                CourseState
	StartingSemester     int
	SemesterCount        int
	Language             string
	LectureHoursPerWeek1 int
	SeminarHoursPerWeek1 int
	LectureHoursPerWeek2 int
	SeminarHoursPerWeek2 int
	Exam                 ExamType
	Credits              int
	Teacher1             string
	Teacher2             string
	MinEnrollment        int
	Capacity             int
}

// TODO: complete enum definition for other states
type CourseState int
type ExamType int

var courseStateName = map[CourseState]string{
	Taught: "taught",
}

func (cs CourseState) String() string {
	return courseStateName[cs]
}

const (
	Taught = iota
)

const (
	ExamTypeAll = iota
)

var examTypeName = map[ExamType]string{
	ExamTypeAll: "C+Ex",
}

func (et ExamType) String() string {
	return examTypeName[et]
}
