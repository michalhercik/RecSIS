package coursedetail

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

type courseDetail struct {
	dbds.Course
	dbds.OverallRating
	CategoryRatings    []dbds.CourseCategoryRating
	BlueprintSemesters pq.BoolArray `db:"semesters"`
	InDegreePlan       bool         `db:"in_degree_plan"`
	// blueprintAssignments []dbds.BlueprintAssignment
}

func (reader DBManager) Course(userID string, code string, lang language.Language) (*Course, error) {
	var result courseDetail
	if err := reader.DB.Get(&result, sqlquery.Course, userID, code, lang); err != nil {
		return nil, err
	}
	if err := reader.DB.Select(&result.CategoryRatings, sqlquery.Rating, userID, code, lang); err != nil {
		return nil, err
	}
	// if err := reader.DB.Select(&result.blueprintAssignments, sqlquery.BlueprintAssignments, userID, code); err != nil {
	// }
	course := intoCourse(&result)
	return &course, nil
}

func (db DBManager) RateCategory(userID string, code string, category string, rating int, lang language.Language) ([]CourseCategoryRating, error) {
	var updatedRating []dbds.CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.RateCategory, userID, code, category, rating)
	if err != nil {
		return []CourseCategoryRating{}, err
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, userID, code, lang); err != nil {
		return []CourseCategoryRating{}, err
	}
	return intoCategoryRatingSlice(updatedRating), err
}

func (db DBManager) DeleteCategoryRating(userID string, code string, category string, lang language.Language) ([]CourseCategoryRating, error) {
	var updatedRating []dbds.CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.DeleteCategoryRating, userID, code, category)
	if err != nil {
		return []CourseCategoryRating{}, err
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, userID, code, lang); err != nil {
		return []CourseCategoryRating{}, err
	}
	return intoCategoryRatingSlice(updatedRating), err
}

func (db DBManager) Rate(userID string, code string, value int) (CourseRating, error) {
	var rating dbds.OverallRating
	_, err := db.DB.Exec(sqlquery.Rate, userID, code, value)
	if err != nil {
		return CourseRating{}, err
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, userID, code); err != nil {
		return CourseRating{}, err
	}
	return intoCourseRating(rating), err
}

func (db DBManager) DeleteRating(userID string, code string) (CourseRating, error) {
	var rating dbds.OverallRating
	_, err := db.DB.Exec(sqlquery.DeleteRating, userID, code)
	if err != nil {
		return CourseRating{}, err
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, userID, code); err != nil {
		return CourseRating{}, err
	}
	return intoCourseRating(rating), err
}

func intoCourse(course *courseDetail) Course {
	return Course{
		Code:                   course.Code,
		Name:                   course.Title,
		Faculty:                course.Faculty,
		GuarantorDepartment:    course.Department,
		State:                  course.State,
		Start:                  TeachingSemester(course.Start),
		Language:               course.Language.String,
		LectureRangeWinter:     course.LectureRangeWinter,
		SeminarRangeWinter:     course.SeminarRangeWinter,
		LectureRangeSummer:     course.LectureRangeSummer,
		SeminarRangeSummer:     course.SeminarRangeSummer,
		ExamType:               course.ExamType,
		Credits:                course.Credits,
		Guarantors:             intoTeacherSlice(course.Guarantors),
		Teachers:               intoTeacherSlice(course.Teachers),
		MinEnrollment:          Capacity(course.MinOccupancy.Int64),
		Capacity:               course.MaxOccupancy.String,
		Annotation:             intoNullDesc(course.Annotation),
		Syllabus:               intoNullDesc(course.Syllabus),
		PassingTerms:           intoNullDesc(course.PassingTerms),
		Literature:             intoNullDesc(course.Literature),
		AssessmentRequirements: intoNullDesc(course.AssesmentRequirements),
		EntryRequirements:      intoNullDesc(course.EntryRequirements),
		Aim:                    intoNullDesc(course.Aim),
		Prereq:                 []string(course.Prereq),
		Coreq:                  []string(course.Coreq),
		Incompa:                []string(course.Incompa),
		Interchange:            []string(course.Interchange),
		Classes:                intoClassSlice(course.Classes),
		Classifications:        intoClassSlice(course.Classifications),
		CourseRating:           intoCourseRating(course.OverallRating),
		CategoryRatings:        intoCategoryRatingSlice(course.CategoryRatings),
		BlueprintAssignments:   intoBlueprintAssignmentSlice(course.BlueprintSemesters),
		BlueprintSemesters:     course.BlueprintSemesters,
		InDegreePlan:           course.InDegreePlan,
	}
}

func intoBlueprintAssignmentSlice(from pq.BoolArray) []Assignment {
	result := []Assignment{}
	if len(from) > 0 && from[0] {
		a := Assignment{year: 0, semester: SemesterAssignment(0)}
		result = append(result, a)
	}
	for j, assigned := range from[1:] {
		if assigned {
			a := Assignment{
				year:     (j / 2) + 1,
				semester: SemesterAssignment((j % 2) + 1),
			}
			result = append(result, a)
		}
	}
	return result
}

func intoCourseRating(from dbds.OverallRating) CourseRating {
	return CourseRating{
		UserRating:  NullInt64(from.UserRating),
		AvgRating:   NullFloat64(from.AvgRating),
		RatingCount: NullInt64(from.Count),
	}
}

func intoCategoryRatingSlice(from []dbds.CourseCategoryRating) []CourseCategoryRating {
	result := make([]CourseCategoryRating, len(from))
	for i, c := range from {
		result[i] = CourseCategoryRating{
			Code:  c.Code,
			Title: c.Title,
			CourseRating: CourseRating{
				UserRating:  NullInt64(c.UserRating),
				AvgRating:   NullFloat64(c.AvgRating),
				RatingCount: NullInt64(c.RatingCount),
			},
		}
	}
	return result
}

func intoClassSlice(from dbds.ClassSlice) []Class {
	result := make([]Class, len(from))
	for i, c := range from {
		result[i] = Class(c)
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
		result[i] = Teacher(t)
	}
	return result
}
