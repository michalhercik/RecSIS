package blueprint

import (
	"github.com/michalhercik/RecSIS/mockdb"
)

type MockDB struct {
	courses []Course
	blueprintCourses map[int][]int
	coursesByYear map[int]map[int][]int
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
	return &MockDB{
		courses: courses,
		blueprintCourses: db.BlueprintCourses,
		coursesByYear: db.CoursesByYear,
	}
}

func (db *MockDB) GetData(user int) BlueprintData {
	unassigned := db.blueprintCourses[user]
	unassignedCourses := make([]Course, len(unassigned))
	for i, idx := range unassigned {
		for _, course := range db.courses {
			if course.Id == idx {
				unassignedCourses[i] = course
				break
			}
		}
	}
	years := db.coursesByYear[user]
	yearsCourses := make(map[int][]Course, len(years)) 
	for i := 1; i <= len(years); i++ {
		yearCourses := make([]Course, len(years[i]))
		for j, idx := range years[i] {
			for _, course := range db.courses {
				if course.Id == idx {
					yearCourses[j] = course
					break
				}
			}
		}
		yearsCourses[i] = yearCourses
	} 
	return BlueprintData{unassigned: unassignedCourses, years: yearsCourses}
}

func (db *MockDB) RemoveUnassigned(user int, courseCode string) {
	blueprintCourses := db.blueprintCourses[user]
	var courseId int
	for _, c := range db.courses {
		if c.Code == courseCode {
			courseId = c.Id
			break
		}
	}
	for i, idx := range blueprintCourses {
		if idx == courseId {
			db.blueprintCourses[user] = append(blueprintCourses[:i], blueprintCourses[i+1:]...)
			return
		}
	}
}

func (db *MockDB) RemoveYear(user int, year int) {
	courses, ok := db.coursesByYear[user][year]
	delete(db.coursesByYear[user], year)
	if !ok {
		return
	}
	db.blueprintCourses[user] = append(db.blueprintCourses[user], courses...)
}

func (db *MockDB) AddYear(user int) {
	db.coursesByYear[user][len(db.coursesByYear[user])+1] = []int{}
}
