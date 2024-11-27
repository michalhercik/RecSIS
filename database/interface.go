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
	Faculty              string
	Guarantor            string
	State                string
	Semester     	     string
	SemesterCount        int
	Language             string
	LectureHoursWinter   int
	SeminarHoursWinter   int
	LectureHoursSummer   int
	SeminarHoursSummer   int
	ExamWinter           string
	ExamSummer           string
	Credits              int
	Teachers             []string
	MinEnrollment        int // -1 means no limit
	Capacity             int // -1 means no limit
}