package degreeplandetail

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/degreeplandetail/internal/sqlquery"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

const (
	foreignKeyViolationCode = "23503"
	notNullViolationCode    = "23502"
)

type DBManager struct {
	DB *sqlx.DB
}

type dbDegreePlanRecord struct {
	Plan       dbds.DegreePlan `db:"plan"`
	BlocCode   string          `db:"bloc_subject_code"`
	BlocLimit  int             `db:"bloc_limit"`
	BlocName   string          `db:"bloc_name"`
	IsRequired bool            `db:"is_required"`
	IsElective bool            `db:"is_elective"`
	dbds.Course
	RecommendedYearFrom sql.NullInt64 `db:"recommended_year_from"`
	RecommendedYearTo   sql.NullInt64 `db:"recommended_year_to"`
	RecommendedSemester sql.NullInt64 `db:"recommended_semester"`
	CourseIsSupported   bool          `db:"course_is_supported"`
	BlueprintSemesters  pq.BoolArray  `db:"semesters"`
}

func (m DBManager) userHasSelectedDegreePlan(uid string) bool {
	var userPlan sql.NullString
	err := m.DB.Get(&userPlan, sqlquery.UserDegreePlanCode, uid)
	if err != nil {
		return false
	}
	return userPlan.Valid
}

func (m DBManager) userDegreePlan(uid string, lang language.Language) (*degreePlanPage, error) {
	var records []dbDegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.UserDegreePlan, uid, lang); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.UserDegreePlan: %w", err), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUserDP,
		)
	}
	if len(records) == 0 {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.UserDegreePlan: no degree plan returned"), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUserDP,
		)
	}
	dp := buildDegreePlanPage(records, true)
	return &dp, nil
}

func (m DBManager) degreePlan(uid, dpCode string, lang language.Language) (*degreePlanPage, error) {
	var records []dbDegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.DegreePlan, uid, dpCode, lang); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DegreePlan: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetDP,
		)
	}
	if len(records) == 0 {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(errors.New("no records found"), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusNotFound,
			texts[lang].errDPNotFound,
		)
	}
	dp := buildDegreePlanPage(records, false)
	return &dp, nil
}

func (m DBManager) saveDegreePlan(uid, dpCode string, lang language.Language) error {
	_, err := m.DB.Exec(sqlquery.SaveDegreePlan, uid, dpCode)
	if err != nil {
		// Handle foreign key violation (invalid user_id or degree_plan_code)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == foreignKeyViolationCode {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.SaveDegreePlan: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusBadRequest,
				texts[lang].errCannotSaveDP,
			)
		} else {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.SaveDegreePlan: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotSaveDP,
			)
		}
	}
	return nil
}

func (m DBManager) deleteSavedDegreePlan(uid string, lang language.Language) (string, error) {
	var planCode sql.NullString
	err := m.DB.QueryRow(sqlquery.DeleteSavedDegreePlan, uid).Scan(&planCode)
	if err != nil {
		return "", errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DeleteSavedDegreePlan: %w", err), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotDeleteSavedDP,
		)
	}
	return planCode.String, nil
}

func (m DBManager) mergeRecommendedPlanWithBlueprint(uid, dpCode string, maxYear int, lang language.Language) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to begin transaction: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotMergeToBlueprint,
		)
	}
	defer tx.Rollback()

	// get current count of blueprint years
	var bpYearCount int
	err = tx.Get(&bpYearCount, sqlquery.CountBlueprintYears, uid)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.CountBlueprintYears: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotMergeToBlueprint,
		)
	}

	if bpYearCount < maxYear {
		// insert missing years
		_, err = tx.Exec(sqlquery.InsertMissingBlueprintYears, uid, bpYearCount+1, maxYear)
		if err != nil {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.InsertMissingBlueprintYears: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotMergeToBlueprint,
			)
		}

		// insert missing semesters
		_, err = tx.Exec(sqlquery.InsertMissingBlueprintSemesters, uid, bpYearCount+1, maxYear)
		if err != nil {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.InsertMissingBlueprintSemesters: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotMergeToBlueprint,
			)
		}
	}

	// insert recommended plan courses to blueprint (merging - duplicates silently ignored via ON CONFLICT)
	_, err = tx.Exec(sqlquery.MergeRecommendedPlanCourses, uid, dpCode, lang)
	if err != nil {
		// Handle NOT NULL constraint violation (missing blueprint year/semester)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == notNullViolationCode {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.MergeRecommendedPlanCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusBadRequest,
				texts[lang].errCannotMergeToBlueprint,
			)
		} else {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.MergeRecommendedPlanCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotMergeToBlueprint,
			)
		}
	}

	if err = tx.Commit(); err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to commit transaction: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotMergeToBlueprint,
		)
	}

	return nil
}

func (m DBManager) rewriteBlueprintWithRecommendedPlan(uid, dpCode string, maxYear int, lang language.Language) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to begin transaction: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotRewriteBlueprint,
		)
	}
	defer tx.Rollback()

	// remove all current blueprint courses
	_, err = tx.Exec(sqlquery.ClearBlueprintCourses, uid)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.ClearBlueprintCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotRewriteBlueprint,
		)
	}

	// get current count of blueprint years
	var bpYearCount int
	err = tx.Get(&bpYearCount, sqlquery.CountBlueprintYears, uid)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.CountBlueprintYears: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotRewriteBlueprint,
		)
	}

	if bpYearCount < maxYear {
		// insert missing years
		_, err = tx.Exec(sqlquery.InsertMissingBlueprintYears, uid, bpYearCount+1, maxYear)
		if err != nil {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.InsertMissingBlueprintYears: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotRewriteBlueprint,
			)
		}

		// insert missing semesters
		_, err = tx.Exec(sqlquery.InsertMissingBlueprintSemesters, uid, bpYearCount+1, maxYear)
		if err != nil {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.InsertMissingBlueprintSemesters: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotRewriteBlueprint,
			)
		}
	}

	// insert recommended plan courses to blueprint
	result, err := tx.Exec(sqlquery.InsertRecommendedPlanCourses, uid, dpCode, lang)
	if err != nil {
		// Handle NOT NULL constraint violation (missing blueprint year/semester)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == notNullViolationCode {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.InsertRecommendedPlanCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusBadRequest,
				texts[lang].errCannotRewriteBlueprint,
			)
		} else {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.InsertRecommendedPlanCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
				http.StatusInternalServerError,
				texts[lang].errCannotRewriteBlueprint,
			)
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.InsertRecommendedPlanCourses: no courses affected"), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusBadRequest,
			texts[lang].errCannotRewriteBlueprint,
		)
	}

	if err = tx.Commit(); err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to commit transaction: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotRewriteBlueprint,
		)
	}

	return nil
}

func buildDegreePlanPage(records []dbDegreePlanRecord, isUserPlan bool) degreePlanPage {
	var dp degreePlanPage
	insertMetadata(&dp, records[0].Plan)
	dp.isUserPlan = isUserPlan
	for _, record := range records {
		add(&dp, record)
	}
	dp.recommendedPlan = createRecommendedPlan(&dp, records)
	return dp
}

func insertMetadata(dp *degreePlanPage, from dbds.DegreePlan) {
	dp.code = from.Code
	dp.title = from.Title
	dp.fieldCode = from.FieldCode
	dp.fieldTitle = from.FieldTitle
	dp.validFrom = from.ValidFrom
	dp.validTo = from.ValidTo
	dp.reqGraphData = from.RequisiteGraphData
	dp.requiredCredits = from.RequiredCredits
	dp.requiredElectiveCredits = from.RequiredElectiveCredits
	dp.totalCredits = from.TotalCredits
	dp.studying = intoStudyingSlice(from.Studying)
	dp.graduates = intoGraduatesSlice(from.Graduates)
}

func intoStudyingSlice(from []dbds.Studying) StudyingSlice {
	studying := make(StudyingSlice, len(from))
	for i, s := range from {
		studying[i] = Studying{
			year:  s.Year,
			count: s.Count,
		}
	}
	return studying
}

func intoGraduatesSlice(from []dbds.Graduates) GraduatesSlice {
	graduates := make(GraduatesSlice, len(from))
	for i, g := range from {
		graduates[i] = Graduates{
			year:     g.Year,
			count:    g.Count,
			avgYears: g.AvgYears,
		}
	}
	return graduates
}

func add(dp *degreePlanPage, record dbDegreePlanRecord) {
	blocIndex := -1
	for i, b := range dp.blocs {
		if b.code == record.BlocCode {
			blocIndex = i
			break
		}
	}
	if blocIndex == -1 {
		dp.blocs = append(dp.blocs, bloc{
			name:         record.BlocName,
			code:         record.BlocCode,
			limit:        record.BlocLimit,
			isCompulsory: record.IsRequired,
			isOptional:   record.IsElective,
		})
		blocIndex = len(dp.blocs) - 1
	}
	dp.blocs[blocIndex].courses = append(dp.blocs[blocIndex].courses, intoCourse(record))
}

func intoCourse(from dbDegreePlanRecord) course {
	return course{
		code:               from.Code,
		title:              from.Title,
		credits:            from.Credits,
		semester:           teachingSemester(from.Start),
		lectureRangeWinter: from.LectureRangeWinter,
		seminarRangeWinter: from.SeminarRangeWinter,
		lectureRangeSummer: from.LectureRangeSummer,
		seminarRangeSummer: from.SeminarRangeSummer,
		examType:           from.ExamType,
		guarantors:         intoTeacherSlice(from.Guarantors),
		isSupported:        from.CourseIsSupported,
		blueprintSemesters: from.BlueprintSemesters,
	}
}

func intoTeacherSlice(from []dbds.Teacher) []teacher {
	teachers := make([]teacher, len(from))
	for i, t := range from {
		teachers[i] = teacher{
			sisID:       t.SisID,
			lastName:    t.LastName,
			firstName:   t.FirstName,
			titleBefore: t.TitleBefore,
			titleAfter:  t.TitleAfter,
		}
	}
	return teachers
}

func createRecommendedPlan(dp *degreePlanPage, records []dbDegreePlanRecord) recommendedPlan {
	coursesMap := make(map[string]course)
	for _, bloc := range dp.blocs {
		for _, c := range bloc.courses {
			coursesMap[c.code] = c
		}
	}

	var recommendedCourseCodes map[int]map[int][]string = make(map[int]map[int][]string)
	for _, record := range records {
		if record.RecommendedYearFrom.Valid && record.RecommendedYearTo.Valid {
			yearFrom := int(record.RecommendedYearFrom.Int64)
			yearTo := int(record.RecommendedYearTo.Int64)
			semester := record.Start
			if semester == int(teachingBoth) {
				if record.RecommendedSemester.Valid {
					semester = int(record.RecommendedSemester.Int64)
				} else {
					semester = int(teachingWinterOnly)
				}
			}
			for y := yearFrom; y <= yearTo; y++ {
				if _, ok := recommendedCourseCodes[y]; !ok {
					recommendedCourseCodes[y] = make(map[int][]string)
				}
				if !slices.Contains(recommendedCourseCodes[y][semester], record.Code) {
					recommendedCourseCodes[y][semester] = append(recommendedCourseCodes[y][semester], record.Code)
				}
			}
		}
	}

	var rp recommendedPlan = recommendedPlan{
		years: []year{},
	}

	for y := 1; y <= len(recommendedCourseCodes); y++ {
		rp.years = append(rp.years, year{
			winter: []course{},
			summer: []course{},
		})

		winterCodes := recommendedCourseCodes[y][int(teachingWinterOnly)]
		for _, wc := range winterCodes {
			if c, ok := coursesMap[wc]; ok {
				rp.years[y-1].winter = append(rp.years[y-1].winter, c)
			}
		}

		summerCodes := recommendedCourseCodes[y][int(teachingSummerOnly)]
		for _, sc := range summerCodes {
			if c, ok := coursesMap[sc]; ok {
				rp.years[y-1].summer = append(rp.years[y-1].summer, c)
			}
		}
	}

	return rp
}
