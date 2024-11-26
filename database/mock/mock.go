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
	result[0] = database.Course{Id: 4950171, Code: "NAIL025", NameCze: "Evoluční algoritmy 1", NameEng: "Evolutionary Algorithms 1", ValidFrom: 2020, ValidTo: 9999, Faculty: 11320, Guarantor: "32-KTIML", State: database.Taught, StartingSemester: 1, SemesterCount: 1, Language: "CZE", LectureHoursPerWeek1: 2, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 0, SeminarHoursPerWeek2: 0, Exam: database.ExamTypeAll, Credits: 5, Teacher1: "Mgr. Roman Neruda, CSc.", Teacher2: "doc. Mgr. Martin Pilát, Ph.D.", MinEnrollment: -1, Capacity: -1}
	result[1] = database.Course{Id: 4950172, Code: "NAIL026", NameCze: "Evoluční algoritmy 2", NameEng: "Evolutionary Algorithms 2", ValidFrom: 2021, ValidTo: 9999, Faculty: 11321, Guarantor: "32-KTIML", State: database.Taught, StartingSemester: 2, SemesterCount: 2, Language: "CZE", LectureHoursPerWeek1: 3, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 0, SeminarHoursPerWeek2: 0, Exam: database.ExamTypeAll, Credits: 5, Teacher1: "Mgr. Roman Neruda, CSc.", Teacher2: "doc. Mgr. Martin Pilát, Ph.D.", MinEnrollment: -1, Capacity: -1}
	result[2] = database.Course{Id: 4950173, Code: "NAIL027", NameCze: "Strojové učení", NameEng: "Machine Learning", ValidFrom: 2020, ValidTo: 9999, Faculty: 11322, Guarantor: "32-KTIML", State: database.Taught, StartingSemester: 1, SemesterCount: 1, Language: "ENG", LectureHoursPerWeek1: 4, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 1, SeminarHoursPerWeek2: 1, Exam: database.ExamTypeAll, Credits: 6, Teacher1: "Dr. John Doe", Teacher2: "Prof. Jane Roe", MinEnrollment: -1, Capacity: -1}
	result[3] = database.Course{Id: 4950174, Code: "NAIL028", NameCze: "Zpracování obrazu", NameEng: "Image Processing", ValidFrom: 2021, ValidTo: 9999, Faculty: 11323, Guarantor: "32-KTIML", State: database.Taught, StartingSemester: 2, SemesterCount: 2, Language: "ENG", LectureHoursPerWeek1: 3, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 1, SeminarHoursPerWeek2: 1, Exam: database.ExamTypeAll, Credits: 5, Teacher1: "Dr. Richard Roe", Teacher2: "Dr. Alice Smith", MinEnrollment: -1, Capacity: -1}
	result[4] = database.Course{Id: 4950175, Code: "NAIL029", NameCze: "Návrh algoritmů", NameEng: "Algorithm Design", ValidFrom: 2020, ValidTo: 9999, Faculty: 11324, Guarantor: "32-KTIML", State: database.Taught, StartingSemester: 1, SemesterCount: 1, Language: "CZE", LectureHoursPerWeek1: 2, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 0, SeminarHoursPerWeek2: 1, Exam: database.ExamTypeAll, Credits: 5, Teacher1: "Dr. Jan Novak", Teacher2: "Prof. Eva Brown", MinEnrollment: -1, Capacity: -1}

	// Continue filling up to 50 records with similar structure
	// Placeholder data is used here for illustrative purposes
	for i := 5; i < 50; i++ {
		result[i] = database.Course{
			Id:                   4950171 + i,
			Code:                 fmt.Sprintf("NAIL%03d", 25+i),
			NameCze:              fmt.Sprintf("Název kurzu %d", i),
			NameEng:              fmt.Sprintf("Course Name %d", i),
			ValidFrom:            2020,
			ValidTo:              9999,
			Faculty:              11000 + i,
			Guarantor:            "32-KTIML",
			State:                database.Taught,
			StartingSemester:     1 + (i % 2),
			SemesterCount:        1 + (i % 2),
			Language:             "ENG",
			LectureHoursPerWeek1: 3 + (i % 3),
			SeminarHoursPerWeek1: 2,
			LectureHoursPerWeek2: 1 + (i % 2),
			SeminarHoursPerWeek2: 1,
			Exam:                 database.ExamTypeAll,
			Credits:              5 + (i % 2),
			Teacher1:             fmt.Sprintf("Dr. Lecturer %d", i),
			Teacher2:             fmt.Sprintf("Prof. Assistant %d", i),
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