package mock_data

import (
	"fmt"
)

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

func GetListOfCourses() []Course {
	result := make([]Course, 50)
	result[0] = Course{Id: 4950171, Code: "NAIL025", NameCze: "Evoluční algoritmy 1", NameEng: "Evolutionary Algorithms 1", ValidFrom: 2020, ValidTo: 9999, Faculty: 11320, Guarantor: "32-KTIML", State: Taught, StartingSemester: 1, SemesterCount: 1, Language: "CZE", LectureHoursPerWeek1: 2, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 0, SeminarHoursPerWeek2: 0, Exam: ExamTypeAll, Credits: 5, Teacher1: "Mgr. Roman Neruda, CSc.", Teacher2: "doc. Mgr. Martin Pilát, Ph.D.", MinEnrollment: -1, Capacity: -1}
	result[1] = Course{Id: 4950172, Code: "NAIL026", NameCze: "Evoluční algoritmy 2", NameEng: "Evolutionary Algorithms 2", ValidFrom: 2021, ValidTo: 9999, Faculty: 11321, Guarantor: "32-KTIML", State: Taught, StartingSemester: 2, SemesterCount: 2, Language: "CZE", LectureHoursPerWeek1: 3, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 0, SeminarHoursPerWeek2: 0, Exam: ExamTypeAll, Credits: 5, Teacher1: "Mgr. Roman Neruda, CSc.", Teacher2: "doc. Mgr. Martin Pilát, Ph.D.", MinEnrollment: -1, Capacity: -1}
	result[2] = Course{Id: 4950173, Code: "NAIL027", NameCze: "Strojové učení", NameEng: "Machine Learning", ValidFrom: 2020, ValidTo: 9999, Faculty: 11322, Guarantor: "32-KTIML", State: Taught, StartingSemester: 1, SemesterCount: 1, Language: "ENG", LectureHoursPerWeek1: 4, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 1, SeminarHoursPerWeek2: 1, Exam: ExamTypeAll, Credits: 6, Teacher1: "Dr. John Doe", Teacher2: "Prof. Jane Roe", MinEnrollment: -1, Capacity: -1}
	result[3] = Course{Id: 4950174, Code: "NAIL028", NameCze: "Zpracování obrazu", NameEng: "Image Processing", ValidFrom: 2021, ValidTo: 9999, Faculty: 11323, Guarantor: "32-KTIML", State: Taught, StartingSemester: 2, SemesterCount: 2, Language: "ENG", LectureHoursPerWeek1: 3, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 1, SeminarHoursPerWeek2: 1, Exam: ExamTypeAll, Credits: 5, Teacher1: "Dr. Richard Roe", Teacher2: "Dr. Alice Smith", MinEnrollment: -1, Capacity: -1}
	result[4] = Course{Id: 4950175, Code: "NAIL029", NameCze: "Návrh algoritmů", NameEng: "Algorithm Design", ValidFrom: 2020, ValidTo: 9999, Faculty: 11324, Guarantor: "32-KTIML", State: Taught, StartingSemester: 1, SemesterCount: 1, Language: "CZE", LectureHoursPerWeek1: 2, SeminarHoursPerWeek1: 2, LectureHoursPerWeek2: 0, SeminarHoursPerWeek2: 1, Exam: ExamTypeAll, Credits: 5, Teacher1: "Dr. Jan Novak", Teacher2: "Prof. Eva Brown", MinEnrollment: -1, Capacity: -1}

	// Continue filling up to 50 records with similar structure
	// Placeholder data is used here for illustrative purposes
	for i := 5; i < 50; i++ {
		result[i] = Course{
			Id:                   4950171 + i,
			Code:                 fmt.Sprintf("NAIL%03d", 25+i),
			NameCze:              fmt.Sprintf("Název kurzu %d", i),
			NameEng:              fmt.Sprintf("Course Name %d", i),
			ValidFrom:            2020,
			ValidTo:              9999,
			Faculty:              11000 + i,
			Guarantor:            "32-KTIML",
			State:                Taught,
			StartingSemester:     1 + (i % 2),
			SemesterCount:        1 + (i % 2),
			Language:             "ENG",
			LectureHoursPerWeek1: 3 + (i % 3),
			SeminarHoursPerWeek1: 2,
			LectureHoursPerWeek2: 1 + (i % 2),
			SeminarHoursPerWeek2: 1,
			Exam:                 ExamTypeAll,
			Credits:              5 + (i % 2),
			Teacher1:             fmt.Sprintf("Dr. Lecturer %d", i),
			Teacher2:             fmt.Sprintf("Prof. Assistant %d", i),
			MinEnrollment:        -1,
			Capacity:             -1,
		}
	}

	return result
}

var BlueprintCourses = []int{4950171, 4950172, 4950173, 4950174, 4950175}

func AddToBlueprint(index int) {
	BlueprintCourses = append(BlueprintCourses, index)
}

func RemoveFromBlueprint(index int) {
	for i, idx := range BlueprintCourses {
		if idx == index {
			BlueprintCourses = append(BlueprintCourses[:i], BlueprintCourses[i+1:]...)
			return
		}
	}
}

func GetBlueprintCourses() []Course {
	courses := GetListOfCourses()
	result := make([]Course, len(BlueprintCourses))
	for i, idx := range BlueprintCourses {
		for _, course := range courses {
			if course.Id == idx {
				result[i] = course
				break
			}
		}
	}
	return result
}

var CoursesByYear = map[int][]int{
	1: {4950182, 4950213, 4950214},
	2: {4950186, 4950205},
}

func AddYear() {
	CoursesByYear[len(CoursesByYear) + 1] = []int{}
}

func RemoveYear(year int) {
	delete(CoursesByYear, year)
}

func AddCourseToYear(year, course int) {
	courses, ok := CoursesByYear[year]
	if !ok {
		return
	}
	// remove duplicates
	for _, idx := range courses {
		if idx == course {
			return
		}
	}
	CoursesByYear[year] = append(courses, course)
}

func RemoveCourseFromYear(year, course int) {
	courses, ok := CoursesByYear[year]
	if !ok {
		return
	}
	for i, idx := range courses {
		if idx == course {
			CoursesByYear[year] = append(courses[:i], courses[i+1:]...)
			return
		}
	}
}

func GetCoursesByYears() map[int][]Course {
	courses := GetListOfCourses()
	result := make(map[int][]Course, len(CoursesByYear)) 
	for i := 1; i <= len(CoursesByYear); i++ {
		year_courses := make([]Course, len(CoursesByYear[i]))
		for j, idx := range CoursesByYear[i] {
			for _, course := range courses {
				if course.Id == idx {
					year_courses[j] = course
					break
				}
			}
		}
		result[i] = year_courses
	} 
	return result
}