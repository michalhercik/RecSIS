package courses

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

const (
	Taught = iota
)

var courseStateName = map[CourseState]string{
	Taught: "taught",
}

func (cs CourseState) String() string {
	return courseStateName[cs]
}

type ExamType int

const (
	ExamTypeAll = iota
)

var examTypeName = map[ExamType]string{
	ExamTypeAll: "C+Ex",
}

func (et ExamType) String() string {
	return examTypeName[et]
}
