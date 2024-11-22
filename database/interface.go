package database

// DB defines the methods for interacting with the database
type DB interface {
	// courses
	GetAllCourses() []Course
	// blueprint
	GetUnassignedBlueprint(user int) []Course
	GetAssignedBlueprint(user int) map[int][]Course
	//AddBlueprintCourse(user int, course int)
	RemoveBlueprintCourse(user int, course int)
	//AddBlueprintYear(user int)
	//RemoveBlueprintYear(user int, year int)
	//AddCourseToYear(user int, year int, course int)
	//RemoveCourseFromYear(user int, year int, course int)
	// TODO add more
}

// Course representation in system
// TODO:
//	- check validity of field types (consider enum types)
//  - consider removing additional or adding missing fields
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

type BlueprintData struct {
	Unassigned []Course
	Years map[int][]Course
}

type CoursesData struct {
	Courses []Course
}
