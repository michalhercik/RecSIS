package courses

import (
	"github.com/michalhercik/RecSIS/mockdb"
)

type MockDB struct {
	courses []Course
}

func CreateDB(db *mockdb.DB) *MockDB {
	var courses []Course
	for _, course := range db.Courses {
		courses = append(courses, Course{
			Id:                 course.Id,
			Code:               course.Code,
			NameCze:            course.NameCze,
			NameEng:            course.NameEng,
			Semester:           course.Semester,
			LectureHoursWinter: course.LectureHoursWinter,
			SeminarHoursWinter: course.SeminarHoursWinter,
			LectureHoursSummer: course.LectureHoursSummer,
			SeminarHoursSummer: course.SeminarHoursSummer,
			ExamWinter:        	course.ExamWinter,
			ExamSummer:        	course.ExamSummer,
			Credits:           	course.Credits,
			Teachers:          	course.Teachers,
		})
	}
	return &MockDB{courses: courses}
}

func (db *MockDB) Courses(query query) coursesPage {
	var courses []Course
	startIndex := query.startIndex
	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex >= len(db.courses) {
		courses = []Course{}
	}

	endIndex := startIndex + query.maxCount
	if endIndex > len(db.courses) {
		endIndex = len(db.courses)
	}

	courses = db.courses[startIndex:endIndex]
	total := len(db.courses)
	count := len(courses)

	return coursesPage{
		courses:    courses,
		startIndex: startIndex,
		count:      count,
		total:      total,
		sorted:     query.sorted,
	}
}