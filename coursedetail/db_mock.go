package coursedetail

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
			ValidFrom:          course.ValidFrom,
			ValidTo:            course.ValidTo,
			Faculty:            course.Faculty,
			Guarantor:          course.Guarantor,
			State:              course.State,
			Semester:           course.Semester,
			SemesterCount:      course.SemesterCount,
			Language:           course.Language,
			LectureHoursWinter: course.LectureHoursWinter,
			SeminarHoursWinter: course.SeminarHoursWinter,
			LectureHoursSummer: course.LectureHoursSummer,
			SeminarHoursSummer: course.SeminarHoursSummer,
			ExamWinter:        	course.ExamWinter,
			ExamSummer:        	course.ExamSummer,
			Credits:           	course.Credits,
			Teachers:          	course.Teachers,
			MinEnrollment:     	course.MinEnrollment,
			Capacity:          	course.Capacity,
			Link: 			    "here should be link to course webpage (not SIS)",
		})
	}
	return &MockDB{courses: courses}
}

func (db *MockDB) GetData(id int) CourseData {
	for _, course := range db.courses {
		if course.Id == id {
			return CourseData{course: course}
		}
	}
	return CourseData{}
}