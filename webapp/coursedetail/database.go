package coursedetail

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

type courseDetail struct {
	dbds.Course
	dbds.OverallRating
	CategoryRatings    []dbds.CourseCategoryRating
	requisites         []dbds.Requisite
	BlueprintSemesters pq.BoolArray `db:"semesters"`
	InDegreePlan       bool         `db:"in_degree_plan"`
}

func (reader DBManager) course(userID string, code string, lang language.Language) (*course, error) {
	var result courseDetail
	if err := reader.DB.Get(&result, sqlquery.Course, userID, code, lang); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorx.NewHTTPErr(
				errorx.AddContext(errors.New("course not found"), errorx.P("code", code), errorx.P("lang", lang)),
				http.StatusNotFound,
				fmt.Sprintf("%s%s%s", texts[lang].errCourseNotFoundPre, code, texts[lang].errCourseNotFoundSuf),
			)
		}
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Course: %w", err), errorx.P("code", code), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetCourse,
		)
	}
	if err := reader.DB.Select(&result.CategoryRatings, sqlquery.Rating, userID, code, lang); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Rating: %w", err), errorx.P("code", code), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetCourseRatings,
		)
	}
	if err := reader.DB.Select(&result.requisites, sqlquery.Requisites, code); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Requisites: %w", err), errorx.P("code", code)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetRequisites,
		)
	}
	course := intoCourse(&result)
	return &course, nil
}

func (db DBManager) rateCategory(userID string, code string, category string, rating int, lang language.Language) ([]courseCategoryRating, error) {
	var updatedRating []dbds.CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.RateCategory, userID, code, category, rating)
	if err != nil {
		return []courseCategoryRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.RateCategory: %w", err), errorx.P("code", code), errorx.P("category", category), errorx.P("rating", rating), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotRateCategory,
		)
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, userID, code, lang); err != nil {
		return []courseCategoryRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Rating: %w", err), errorx.P("code", code), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUpdatedRatings,
		)
	}
	return intoCategoryRatingSlice(updatedRating), nil
}

func (db DBManager) deleteCategoryRating(userID string, code string, category string, lang language.Language) ([]courseCategoryRating, error) {
	var updatedRating []dbds.CourseCategoryRating
	res, err := db.DB.Exec(sqlquery.DeleteCategoryRating, userID, code, category)
	if err != nil {
		return []courseCategoryRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DeleteCategoryRating: %w", err), errorx.P("code", code), errorx.P("category", category), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotDeleteRating,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return []courseCategoryRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DeleteCategoryRating: %w", err), errorx.P("code", code), errorx.P("category", category), errorx.P("lang", lang)),
			http.StatusBadRequest,
			texts[lang].errCannotDeleteRating,
		)
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, userID, code, lang); err != nil {
		return []courseCategoryRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Rating: %w", err), errorx.P("code", code), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUpdatedRatings,
		)
	}
	return intoCategoryRatingSlice(updatedRating), nil
}

func (db DBManager) rate(userID string, code string, value int, lang language.Language) (courseRating, error) {
	var rating dbds.OverallRating
	_, err := db.DB.Exec(sqlquery.Rate, userID, code, value)
	if err != nil {
		return courseRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.Rate: %w", err), errorx.P("code", code), errorx.P("value", value)),
			http.StatusInternalServerError,
			texts[lang].errUnableToRateCourse,
		)
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, userID, code); err != nil {
		return courseRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.CourseOverallRating: %w", err), errorx.P("code", code)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUpdatedRatings,
		)
	}
	return intoCourseRating(rating), nil
}

func (db DBManager) deleteRating(userID string, code string, lang language.Language) (courseRating, error) {
	var rating dbds.OverallRating
	res, err := db.DB.Exec(sqlquery.DeleteRating, userID, code)
	if err != nil {
		return courseRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DeleteRating: %w", err), errorx.P("code", code)),
			http.StatusInternalServerError,
			texts[lang].errCannotDeleteRating,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return courseRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DeleteRating: %w", err), errorx.P("code", code)),
			http.StatusBadRequest,
			texts[lang].errCannotDeleteRating,
		)
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, userID, code); err != nil {
		return courseRating{}, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.CourseOverallRating: %w", err), errorx.P("code", code)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUpdatedRatings,
		)
	}
	return intoCourseRating(rating), nil
}

func intoCourse(from *courseDetail) course {
	return course{
		code:                   from.Code,
		title:                  from.Title,
		faculty:                intoFaculty(from.Faculty),
		guarantorDepartment:    intoDepartment(from.Department),
		state:                  from.State,
		semester:               teachingSemester(from.Start),
		language:               from.Language.String,
		lectureRangeWinter:     from.LectureRangeWinter,
		seminarRangeWinter:     from.SeminarRangeWinter,
		lectureRangeSummer:     from.LectureRangeSummer,
		seminarRangeSummer:     from.SeminarRangeSummer,
		rangeUnit:              intoRangeUnit(from.RangeUnit),
		examType:               from.ExamType,
		credits:                from.Credits,
		guarantors:             intoTeacherSlice(from.Guarantors),
		teachers:               intoTeacherSlice(from.Teachers),
		capacity:               from.MaxOccupancy.String,
		url:                    from.URL,
		annotation:             intoNullDesc(from.Annotation),
		syllabus:               intoNullDesc(from.Syllabus),
		passingTerms:           intoNullDesc(from.PassingTerms),
		literature:             intoNullDesc(from.Literature),
		assessmentRequirements: intoNullDesc(from.AssessmentRequirements),
		entryRequirements:      intoNullDesc(from.EntryRequirements),
		aim:                    intoNullDesc(from.Aim),
		prerequisitesRoot:      intoRequisiteTree(from.requisites, from.Code, "P"),
		corequisitesRoot:       intoRequisiteTree(from.requisites, from.Code, "K"),
		incompatiblesRoot:      intoRequisiteTree(from.requisites, from.Code, "N"),
		interchangesRoot:       intoRequisiteTree(from.requisites, from.Code, "Z"),
		classes:                []string(from.Classes),
		classifications:        []string(from.Classifications),
		overallRating:          intoCourseRating(from.OverallRating),
		categoryRatings:        intoCategoryRatingSlice(from.CategoryRatings),
		blueprintAssignments:   intoBlueprintAssignmentSlice(from.BlueprintSemesters),
		blueprintSemesters:     from.BlueprintSemesters,
		inDegreePlan:           from.InDegreePlan,
	}
}

func intoRangeUnit(from dbds.NullRangeUnit) nullRangeUnit {
	return nullRangeUnit{
		rangeUnit: rangeUnit{
			abbr: from.Abbr,
			name: from.Name,
		},
		valid: from.Valid,
	}
}

func intoFaculty(from dbds.Faculty) faculty {
	return faculty{
		abbr: from.Abbr,
		name: from.Name,
	}
}

func intoDepartment(from dbds.Department) department {
	return department{
		id:   from.ID,
		name: from.Name,
	}
}

func intoRequisiteTree(from []dbds.Requisite, rootCourse string, reqType string) *requisiteNode {
	if from == nil {
		return nil
	}

	filteredByType := []dbds.Requisite{}
	for _, req := range from {
		if req.Type == reqType {
			filteredByType = append(filteredByType, req)
		}
	}

	nodes := map[string]*requisiteNode{}
	getNode := func(course string) *requisiteNode {
		if _, ok := nodes[course]; !ok {
			nodes[course] = &requisiteNode{courseCode: course}
		}
		return nodes[course]
	}

	var root *requisiteNode
	for _, r := range filteredByType {
		parentNode := getNode(r.Parent)
		childNode := getNode(r.Child)
		childNode.groupType = r.Group
		parentNode.children = append(parentNode.children, childNode)
		if r.Parent == rootCourse {
			root = parentNode
		}
	}

	return root
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
			SisID:       t.SisID,
			LastName:    t.LastName,
			FirstName:   t.FirstName,
			TitleBefore: t.TitleBefore,
			TitleAfter:  t.TitleAfter,
		}
	}
	return result
}
