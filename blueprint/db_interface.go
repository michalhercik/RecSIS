package blueprint

type Database interface {
	GetData(user int) BlueprintData
	RemoveUnassigned(user int, course int)
	RemoveYear(user int, year int)
	AddYear(user int)
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

type BlueprintData struct {
	unassigned []Course
	years map[int][]Course
}