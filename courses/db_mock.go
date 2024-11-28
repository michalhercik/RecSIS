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

func (db *MockDB) GetData() CoursesData {
	return CoursesData{courses: db.courses}
}