package mockdb

import (
	"fmt"
)

const user = 42 // TODO get user from session

// interfaces

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

// init

type DB struct {
	Courses []Course
	BlueprintCourses map[int][]int
	CoursesByYear map[int]map[int][]int
}

var db DB

func NewDB() *DB {
	newDB := &DB{
		Courses: getListOfCourses(),
		BlueprintCourses: getBlueprintCourses(),
		CoursesByYear: getCoursesByYears(),
	}
	db = *newDB
	return newDB
}

func getListOfCourses() []Course {
	result := make([]Course, 50)
	result[0] = Course{
		Id:                   4950171,
		Code:                 "NAIL025",
		NameCze:              "Evoluční algoritmy 1",
		NameEng:              "Evolutionary Algorithms 1",
		ValidFrom:            2020,
		ValidTo:              9999,
		Faculty:              "MFF",
		Guarantor:            "32-KTIML",
		State:                "taught",
		Semester:   		  "Winter",
		SemesterCount:        1,
		Language:             "CZE",
		LectureHoursWinter:   2,
		SeminarHoursWinter:   2,
		LectureHoursSummer:   0,
		SeminarHoursSummer:   0,
		ExamWinter:           "C+Ex",
		ExamSummer:           "",
		Credits:              5,
		Teachers:             []string{ "Mgr. Roman Neruda, CSc.", "doc. Mgr. Martin Pilát, Ph.D." },
		MinEnrollment:        -1,
		Capacity:             50,
	}
	result[1] = Course{
		Id:                   4950172,
		Code:                 "NAIL026",
		NameCze:              "Evoluční algoritmy 2",
		NameEng:              "Evolutionary Algorithms 2",
		ValidFrom:            2021,
		ValidTo:              9999,
		Faculty:              "MFF",
		Guarantor:            "32-KTIML",
		State:                "taught",
		Semester:     		  "Summer",
		SemesterCount:        1,
		Language:             "CZE",
		LectureHoursWinter:   0,
		SeminarHoursWinter:   0,
		LectureHoursSummer:   3,
		SeminarHoursSummer:   2,
		ExamWinter:           "",
		ExamSummer:           "C+Ex",
		Credits:              5,
		Teachers:             []string{ "Mgr. Roman Neruda, CSc.", "doc. Mgr. Martin Pilát, Ph.D." },
		MinEnrollment:        -1,
		Capacity:             -1,
	}
	result[2] = Course{
		Id:                   4950173,
		Code:                 "NAIL027",
		NameCze:              "Strojové učení",
		NameEng:              "Machine Learning",
		ValidFrom:            2020,
		ValidTo:              9999,
		Faculty:              "MFF",
		Guarantor:            "32-KTIML",
		State:                "taught",
		Semester:     		  "Winter",
		SemesterCount:        1,
		Language:             "ENG",
		LectureHoursWinter:   4,
		SeminarHoursWinter:   2,
		LectureHoursSummer:   0,
		SeminarHoursSummer:   0,
		ExamWinter:           "C+Ex",
		ExamSummer:           "",
		Credits:              6,
		Teachers:             []string{ "Dr. John Doe",  "Prof. Jane Roe" },
		MinEnrollment:        -1,
		Capacity:             100,
	}
	result[3] = Course{
		Id:                   4950174,
		Code:                 "NAIL028",
		NameCze:              "Zpracování obrazu",
		NameEng:              "Image Processing",
		ValidFrom:            2021,
		ValidTo:              9999,
		Faculty:              "MFF",
		Guarantor:            "32-KTIML",
		State:                "taught",
		Semester:    		  "Summer",
		SemesterCount:        2,
		Language:             "ENG",
		LectureHoursWinter:   3,
		SeminarHoursWinter:   2,
		LectureHoursSummer:   1,
		SeminarHoursSummer:   1,
		ExamWinter:           "C+Ex",
		ExamSummer:           "C+Ex",
		Credits:              5,
		Teachers:             []string{ "Dr. Richard Roe" },
		MinEnrollment:        -1,
		Capacity:             10,
	}
	result[4] = Course{
		Id:                   4950175,
		Code:                 "NAIL029",
		NameCze:              "Návrh algoritmů",
		NameEng:              "Algorithm Design",
		ValidFrom:            2020,
		ValidTo:              9999,
		Faculty:              "MFF",
		Guarantor:            "32-KTIML",
		State:                "taught",
		Semester:             "Winter",
		SemesterCount:        2,
		Language:             "CZE",
		LectureHoursWinter:	  2,
		SeminarHoursWinter:   2,
		LectureHoursSummer:   1,
		SeminarHoursSummer:   1,
		ExamWinter:           "C+Ex",
		ExamSummer:           "C+Ex",
		Credits:              5,
		Teachers:             []string{ "Dr. Jan Novak", "Prof. Eva Brown", "prof. RNDr. Michal Hercik PhD." },
		MinEnrollment:        -1,
		Capacity:             -1,
	}

	// Continue filling up to 50 records with similar structure
	// Placeholder data is used here for illustrative purposes
	for i := 5; i < 50; i++ {
		semester := 1 + (i % 2)
		semesterStr := "Winter"
		if semester == 2 {
			semesterStr = "Summer"
		}
		semesterCount := 1 + (i % 2)
		lectureHoursWinter := 0
		seminarHoursWinter := 0
		lectureHoursSummer := 0
		seminarHoursSummer := 0

		if semesterCount == 1 {
			if semester == 1 {
				lectureHoursWinter = 1 + (i % 3)
				seminarHoursWinter = 1 + (i % 3)
			} else {
				lectureHoursSummer = 1 + (i % 3)
				seminarHoursSummer = 1 + (i % 3)
			}
		} else {
			lectureHoursWinter = 1 + (i % 3)
			seminarHoursWinter = 1 + (i % 3)
			lectureHoursSummer = 1 + (i % 3)
			seminarHoursSummer = 1 + (i % 3)
		}

		result[i] = Course{
			Id:                   4950171 + i,
			Code:                 fmt.Sprintf("NAIL%03d", 25+i),
			NameCze:              fmt.Sprintf("Název kurzu %d", i),
			NameEng:              fmt.Sprintf("Course Name %d", i),
			ValidFrom:            2020,
			ValidTo:              9999,
			Faculty:              "MFF",
			Guarantor:            "32-KTIML",
			State:                "taught",
			Semester:             semesterStr,
			SemesterCount:        semesterCount,
			Language:             "ENG",
			LectureHoursWinter:   lectureHoursWinter,
			SeminarHoursWinter:   seminarHoursWinter,
			LectureHoursSummer:   lectureHoursSummer,
			SeminarHoursSummer:   seminarHoursSummer,
			ExamWinter:           "C+Ex",
			ExamSummer:           "C+Ex",
			Credits:              4 + (i % 2),
			Teachers:             []string{ fmt.Sprintf("Dr. Lecturer %d", i), fmt.Sprintf("Prof. Assistant %d", i) },
			MinEnrollment:        -1,
			Capacity:             -1,
		}
	}

	return result
}

func getBlueprintCourses() map[int][]int {
	return map[int][]int{
		user: {4950171, 4950172, 4950173, 4950174, 4950175},
	}
}

func getCoursesByYears() map[int]map[int][]int {
	return map[int]map[int][]int{
		user: {
			1: {4950182, 4950213, 4950214},
			2: {4950186, 4950205},
		},
	}
}