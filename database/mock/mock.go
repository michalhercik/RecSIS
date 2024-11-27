package mock

import (
	"fmt"
    "github.com/michalhercik/RecSIS/database"
)

// init

type MockDB struct {
	courses []database.Course
	blueprintCourses map[int][]int
	coursesByYear map[int]map[int][]int
}

var db MockDB

func NewMockDB() *MockDB {
	newDB := &MockDB{
		courses: getListOfCourses(),
		blueprintCourses: getBlueprintCourses(),
		coursesByYear: getCoursesByYears(),
	}
	db = *newDB
	return newDB
}

// interface methods

func (db *MockDB) GetAllCourses() []database.Course {
	return db.courses
}

func (db *MockDB) GetCourse(id int) (database.Course, error) {
	for _, course := range db.courses {
		if course.Id == id {
			return course, nil
		}
	}
	return database.Course{}, fmt.Errorf("Course with id %d not found", id)
}

func (db *MockDB) BlueprintGetUnassigned(user int) []database.Course {
	blueprintCourses := db.blueprintCourses[user]
	result := make([]database.Course, len(blueprintCourses))
	for i, idx := range blueprintCourses {
		for _, course := range db.courses {
			if course.Id == idx {
				result[i] = course
				break
			}
		}
	}
	return result
}

func (db *MockDB) BlueprintGetAssigned(user int) map[int][]database.Course {
	coursesByYear := db.coursesByYear[user]
	result := make(map[int][]database.Course, len(coursesByYear)) 
	for i := 1; i <= len(coursesByYear); i++ {
		yearCourses := make([]database.Course, len(coursesByYear[i]))
		for j, idx := range coursesByYear[i] {
			for _, course := range db.courses {
				if course.Id == idx {
					yearCourses[j] = course
					break
				}
			}
		}
		result[i] = yearCourses
	} 
	return result
}

func (db *MockDB) BlueprintAddUnassigned(user int, course int) {
	db.blueprintCourses[user] = append(db.blueprintCourses[user], course)
}

func (db *MockDB) BlueprintRemoveUnassigned(user int, course int) {
	blueprintCourses := db.blueprintCourses[user]
	for i, idx := range blueprintCourses {
		if idx == course {
			db.blueprintCourses[user] = append(blueprintCourses[:i], blueprintCourses[i+1:]...)
			return
		}
	}
}

func (db *MockDB) BlueprintRemoveYear(user int, year int) []database.Course {
	courses, ok := db.coursesByYear[user][year]
	delete(db.coursesByYear[user], year)
	if !ok {
		return []database.Course{}
	}
	result := make([]database.Course, len(courses))
	for i, idx := range courses {
		for _, course := range db.courses {
			if course.Id == idx {
				result[i] = course
				break
			}
		}
	}
	return result
}

func (db *MockDB) BlueprintAddYear(user int) {
	db.coursesByYear[user][len(db.coursesByYear[user]) + 1] = []int{}
}

// database logic

func getListOfCourses() []database.Course {
	result := make([]database.Course, 50)
	result[0] = database.Course{
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
	result[1] = database.Course{
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
	result[2] = database.Course{
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
	result[3] = database.Course{
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
	result[4] = database.Course{
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

		result[i] = database.Course{
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
		database.User: {4950171, 4950172, 4950173, 4950174, 4950175},
	}
}

func getCoursesByYears() map[int]map[int][]int {
	return map[int]map[int][]int{
		database.User: {
			1: {4950182, 4950213, 4950214},
			2: {4950186, 4950205},
		},
	}
}

// func RemoveFromBlueprint(index int) {
// 	for i, idx := range BlueprintCourses {
// 		if idx == index {
// 			BlueprintCourses = append(BlueprintCourses[:i], BlueprintCourses[i+1:]...)
// 			return
// 		}
// 	}
// }

// func AddCourseToYear(year, course int) {
// 	courses, ok := CoursesByYear[year]
// 	if !ok {
// 		return
// 	}
// 	// remove duplicates
// 	for _, idx := range courses {
// 		if idx == course {
// 			return
// 		}
// 	}
// 	CoursesByYear[year] = append(courses, course)
// }

// func RemoveCourseFromYear(year, course int) {
// 	courses, ok := CoursesByYear[year]
// 	if !ok {
// 		return
// 	}
// 	for i, idx := range courses {
// 		if idx == course {
// 			CoursesByYear[year] = append(courses[:i], courses[i+1:]...)
// 			return
// 		}
// 	}
// }