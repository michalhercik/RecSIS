package courses

type Database interface {
	GetData() CoursesData
}

var db Database

func SetDatabase(newDB Database) {
	db = newDB
}

type Course struct {
	Id                   int
	Code                 string
	NameCze              string
	NameEng              string
	//ValidFrom            int
	//ValidTo              int
	//Faculty              string
	//Guarantor            string
	//State                string
	Semester     	     string
	//SemesterCount        int
	//Language             string
	LectureHoursWinter   int
	SeminarHoursWinter   int
	LectureHoursSummer   int
	SeminarHoursSummer   int
	ExamWinter           string
	ExamSummer           string
	Credits              int
	Teachers             []string
	// MinEnrollment        int // -1 means no limit
	// Capacity             int // -1 means no limit
}

type CoursesData struct {
	courses []Course
}