package courses

type Database interface {
	Courses(query query) coursesPage
}

var db Database

func SetDatabase(newDB Database) {
	db = newDB
}

type coursesPage struct {
	courses    []Course
	startIndex int
	count      int
	total      int
	search     string
	sorted     sortType
}

type query struct {
	user       int
	startIndex int
	maxCount   int
	search     string
	sorted     sortType 
}

const (
	relevance    = iota
	recommended  = iota
	rating       = iota
	mostPopular  = iota
	newest       = iota
)

type sortType int

var sortTypeName = map[sortType]string{
	relevance:   "Relevant",
	recommended: "Recommended",
	rating:      "By rating",
	mostPopular: "Most popular",
	newest:		 "Newest",
}

func (st sortType) String() string {
	return sortTypeName[st]
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