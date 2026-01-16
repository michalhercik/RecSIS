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

type DBManager struct {
	DB *sqlx.DB
}

type dbDegreePlanRecord struct {
	DegreePlanCode      string `db:"degree_plan_code"`
	DegreePlanTitle     string `db:"degree_plan_title"`
	FieldCode           string `db:"field_code"`
	FieldTitle          string `db:"field_title"`
	DegreePlanValidFrom int    `db:"degree_plan_valid_from"`
	DegreePlanValidTo   int    `db:"degree_plan_valid_to"`
	BlocCode            string `db:"bloc_subject_code"`
	BlocLimit           int    `db:"bloc_limit"`
	BlocName            string `db:"bloc_name"`
	BlocType            string `db:"bloc_type"`
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
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.SaveDegreePlan: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotSaveDP,
		)
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
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.MergeRecommendedPlanCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotMergeToBlueprint,
		)
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
	_, err = tx.Exec(sqlquery.InsertRecommendedPlanCourses, uid, dpCode, lang)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.InsertRecommendedPlanCourses: %w", err), errorx.P("dpCode", dpCode), errorx.P("lang", lang)),
			http.StatusInternalServerError,
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
	dp.code = records[0].DegreePlanCode
	dp.title = records[0].DegreePlanTitle
	dp.fieldCode = records[0].FieldCode
	dp.fieldTitle = records[0].FieldTitle
	dp.validFrom = records[0].DegreePlanValidFrom
	dp.validTo = records[0].DegreePlanValidTo
	dp.isUserPlan = isUserPlan
	for _, record := range records {
		add(&dp, record)
	}
	dp.recommendedPlan = createRecommendedPlan(records, &dp)
	fixLimits(&dp)
	return dp
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
			isCompulsory: record.BlocType == "A",
			isOptional:   record.BlocType == "C",
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

func fixLimits(dp *degreePlanPage) {
	for i := range dp.blocs {
		if dp.blocs[i].isCompulsory {
			creditSum := 0
			for _, c := range dp.blocs[i].courses {
				creditSum += c.credits
			}
			dp.blocs[i].limit = creditSum
		} else if dp.blocs[i].isOptional {
			dp.blocs[i].limit = 0
		}
	}
}

func createRecommendedPlan(records []dbDegreePlanRecord, dp *degreePlanPage) recommendedPlan {
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
