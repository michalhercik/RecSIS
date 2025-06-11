package courses

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/courses/internal/sqlquery"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

type courses []struct {
	dbds.Course
	BlueprintSemesters pq.BoolArray `db:"semesters"`
	InDegreePlan       bool         `db:"in_degree_plan"`
}

func (m DBManager) courses(userID string, courseCodes []string, lang language.Language) ([]course, error) {
	var result courses
	if err := m.DB.Select(&result, sqlquery.Courses, userID, pq.Array(courseCodes), lang); err != nil {
		return nil, fmt.Errorf("failed to fetch courses: %w", err)
	}
	courses := intoCourses(result)
	return courses, nil
}

func intoCourses(from courses) []course {
	result := make([]course, len(from))
	for i, course := range from {
		result[i].code = course.Code
		result[i].title = course.Title
		result[i].annotation = intoNullDesc(course.Annotation)
		result[i].semester = teachingSemester(course.Start)
		result[i].lectureRangeWinter = course.LectureRangeWinter
		result[i].seminarRangeWinter = course.SeminarRangeWinter
		result[i].lectureRangeSummer = course.LectureRangeSummer
		result[i].seminarRangeSummer = course.SeminarRangeSummer
		result[i].examType = course.ExamType
		result[i].credits = course.Credits
		result[i].guarantors = intoTeacherSlice(course.Guarantors)
		result[i].blueprintSemesters = course.BlueprintSemesters
		result[i].blueprintAssignments = intoBlueprintAssignmentSlice(course.BlueprintSemesters)
		result[i].inDegreePlan = course.InDegreePlan
	}
	return result
}

func intoNullDesc(from dbds.NullDescription) nullDescription {
	return nullDescription{
		description: description{
			title:   from.Title,
			content: from.Content,
		},
		valid: from.Valid,
	}
}

func intoTeacherSlice(from dbds.TeacherSlice) []teacher {
	result := make([]teacher, len(from))
	for i, t := range from {
		result[i] = teacher{
			sisID:       t.SisID,
			firstName:   t.FirstName,
			lastName:    t.LastName,
			titleBefore: t.TitleBefore,
			titleAfter:  t.TitleAfter,
		}
	}
	return result
}

func intoBlueprintAssignmentSlice(from pq.BoolArray) []assignment {
	result := []assignment{}
	if len(from) > 0 && from[0] {
		a := assignment{year: 0, semester: semesterAssignment(0)}
		result = append(result, a)
	}
	for j, assigned := range from[1:] {
		if assigned {
			a := assignment{
				year:     (j / 2) + 1,
				semester: semesterAssignment((j % 2) + 1),
			}
			result = append(result, a)
		}
	}
	return result
}
