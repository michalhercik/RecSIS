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

func (reader DBManager) Course(userID string, code string, lang language.Language) (*course, error) {
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

func (db DBManager) RateCategory(userID string, code string, category string, rating int, lang language.Language) ([]courseCategoryRating, error) {
	var updatedRating []dbds.CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.RateCategory, userID, code, category, rating)
	if err != nil {
		return []courseCategoryRating{}, err
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, userID, code, lang); err != nil {
		return []courseCategoryRating{}, err
	}
	return intoCategoryRatingSlice(updatedRating), err
}

func (db DBManager) DeleteCategoryRating(userID string, code string, category string, lang language.Language) ([]courseCategoryRating, error) {
	var updatedRating []dbds.CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.DeleteCategoryRating, userID, code, category)
	if err != nil {
		return []courseCategoryRating{}, err
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, userID, code, lang); err != nil {
		return []courseCategoryRating{}, err
	}
	return intoCategoryRatingSlice(updatedRating), err
}

func (db DBManager) Rate(userID string, code string, value int) (courseRating, error) {
	var rating dbds.OverallRating
	_, err := db.DB.Exec(sqlquery.Rate, userID, code, value)
	if err != nil {
		return courseRating{}, err
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, userID, code); err != nil {
		return courseRating{}, err
	}
	return intoCourseRating(rating), err
}

func (db DBManager) DeleteRating(userID string, code string) (courseRating, error) {
	var rating dbds.OverallRating
	_, err := db.DB.Exec(sqlquery.DeleteRating, userID, code)
	if err != nil {
		return courseRating{}, err
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, userID, code); err != nil {
		return courseRating{}, err
	}
	return intoCourseRating(rating), err
}

func intoCourse(from *courseDetail) course {
	return course{
		code:                   from.Code,
		title:                  from.Title,
		faculty:                from.Faculty,
		guarantorDepartment:    from.Department,
		state:                  from.State,
		semester:               teachingSemester(from.Start),
		language:               from.Language.String,
		lectureRangeWinter:     from.LectureRangeWinter,
		seminarRangeWinter:     from.SeminarRangeWinter,
		lectureRangeSummer:     from.LectureRangeSummer,
		seminarRangeSummer:     from.SeminarRangeSummer,
		examType:               from.ExamType,
		credits:                from.Credits,
		guarantors:             intoTeacherSlice(from.Guarantors),
		teachers:               intoTeacherSlice(from.Teachers),
		capacity:               from.MaxOccupancy.String,
		annotation:             intoNullDesc(from.Annotation),
		syllabus:               intoNullDesc(from.Syllabus),
		passingTerms:           intoNullDesc(from.PassingTerms),
		literature:             intoNullDesc(from.Literature),
		assessmentRequirements: intoNullDesc(from.AssessmentRequirements),
		entryRequirements:      intoNullDesc(from.EntryRequirements),
		aim:                    intoNullDesc(from.Aim),
		prerequisites:          []string(from.Prereq),
		corequisites:           []string(from.Coreq),
		incompatible:           []string(from.Incompa),
		interchange:            []string(from.Interchange),
		classes:                intoClassSlice(from.Classes),
		classifications:        intoClassSlice(from.Classifications),
		overallRating:          intoCourseRating(from.OverallRating),
		categoryRatings:        intoCategoryRatingSlice(from.CategoryRatings),
		blueprintAssignments:   intoBlueprintAssignmentSlice(from.BlueprintSemesters),
		blueprintSemesters:     from.BlueprintSemesters,
		inDegreePlan:           from.InDegreePlan,
	}
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

func intoCourseRating(from dbds.OverallRating) courseRating {
	return courseRating{
		userRating:  from.UserRating,
		avgRating:   from.AvgRating,
		ratingCount: from.Count,
	}
}

func intoCategoryRatingSlice(from []dbds.CourseCategoryRating) []courseCategoryRating {
	result := make([]courseCategoryRating, len(from))
	for i, c := range from {
		result[i] = courseCategoryRating{
			code:  c.Code,
			title: c.Title,
			courseRating: courseRating{
				userRating:  c.UserRating,
				avgRating:   c.AvgRating,
				ratingCount: c.RatingCount,
			},
		}
	}
	return result
}

func intoClassSlice(from dbds.ClassSlice) []class {
	result := make([]class, len(from))
	for i, c := range from {
		result[i] = class{
			code: c.Code,
			name: c.Name,
		}
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
			SISID:       t.SISID,
			LastName:    t.LastName,
			FirstName:   t.FirstName,
			TitleBefore: t.TitleBefore,
			TitleAfter:  t.TitleAfter,
		}
	}
	return result
}
