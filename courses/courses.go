package courses

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
