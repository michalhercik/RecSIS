package coursedetail

import (
	"github.com/michalhercik/RecSIS/mockdb"
)

// TODO: Maybe change to better name
type MockCourseReader struct {
	courses []Course
}

func NewMockCourseReader(db *mockdb.DB) *MockCourseReader {
	var courses []Course
	// for _, course := range db.Courses {
	// 	courses = append(courses, Course{
	// 		Id:            course.Id,
	// 		Code:          course.Code,
	// 		NameCs:        course.NameCs,
	// 		NameEn:        course.NameEn,
	// 		ValidFrom:     course.ValidFrom,
	// 		ValidTo:       course.ValidTo,
	// 		Faculty:       course.Faculty,
	// 		Guarantor:     course.Guarantor,
	// 		State:         course.State,
	// 		Start:         course.Start,
	// 		SemesterCount: course.SemesterCount,
	// 		Language:      course.Language,
	// 		LectureRange1: course.LectureRange1,
	// 		SeminarRange1: course.SeminarRange1,
	// 		LectureRange2: course.LectureRange2,
	// 		SeminarRange2: course.SeminarRange2,
	// 		ExamType:      course.ExamType,
	// 		Credits:       course.Credits,
	// 		Teacher1:      course.Teacher1,
	// 		Teacher2:      course.Teacher2,
	// 		MinEnrollment: course.MinEnrollment,
	// 		Capacity:      course.Capacity,
	// 		Link:          "here should be link to course webpage (not SIS)",
	// 	})
	// }
	return &MockCourseReader{courses: courses}
}

func (db *MockCourseReader) Course(code string) CourseData {
	for _, course := range db.courses {
		if course.Code == code {
			return CourseData{course: course}
		}
	}
	return CourseData{}
}
