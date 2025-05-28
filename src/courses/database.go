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
}

func (m DBManager) Courses(userID string, courseCodes []string, lang language.Language) ([]Course, error) {
	dbResult := courses{}
	if err := m.DB.Select(&dbResult, sqlquery.Courses, userID, pq.Array(courseCodes), lang); err != nil {
		return nil, fmt.Errorf("failed to fetch courses: %w", err)
	}
	result := intoCourses(dbResult)
	return result, nil
}

func intoCourses(from courses) []Course {
	result := make([]Course, len(from))
	for i, course := range from {
		result[i].Code = course.Code
		result[i].Name = course.Title
		result[i].Annotation = intoNullDesc(course.Annotation)
		result[i].Start = TeachingSemester(course.Start)
		result[i].LectureRange1 = int(course.LectureRangeWinter.Int64)
		result[i].SeminarRange1 = int(course.SeminarRangeWinter.Int64)
		result[i].LectureRange2 = int(course.LectureRangeSummer.Int64)
		result[i].SeminarRange2 = int(course.SeminarRangeSummer.Int64)
		result[i].ExamType = course.ExamType
		result[i].Credits = course.Credits
		result[i].Guarantors = intoTeacherSlice(course.Guarantors)
		result[i].BlueprintSemesters = course.BlueprintSemesters
		result[i].BlueprintAssignments = intoBlueprintAssignmentSlice(course.BlueprintSemesters)

	}
	return result
}

func intoNullDesc(from dbds.NullDescription) NullDescription {
	return NullDescription{
		Description: Description(from.Description),
		Valid:       from.Valid,
	}
}

func intoTeacherSlice(from dbds.TeacherSlice) []Teacher {
	result := make([]Teacher, len(from))
	for i, t := range from {
		result[i] = Teacher{
			SisID:       t.SISID,
			FirstName:   t.FirstName,
			LastName:    t.LastName,
			TitleBefore: t.TitleBefore,
			TitleAfter:  t.TitleAfter,
		}
	}
	return result
}

func intoBlueprintAssignmentSlice(from pq.BoolArray) []Assignment {
	result := []Assignment{}
	if len(from) > 0 && from[0] {
		a := Assignment{Year: 0, Semester: SemesterAssignment(0)}
		result = append(result, a)
	}
	for j, assigned := range from[1:] {
		if assigned {
			a := Assignment{
				Year:     (j / 2) + 1,
				Semester: SemesterAssignment((j % 2) + 1),
			}
			result = append(result, a)
		}
	}
	return result
}
