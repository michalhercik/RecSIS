package degreeplan

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/degreeplan/internal/sqlquery"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

type dbDegreePlanRecord struct {
	DegreePlanCode string `db:"degree_plan_code"`
	DegreePlanYear int    `db:"start_year"`
	BlocCode       string `db:"bloc_subject_code"`
	BlocLimit      int    `db:"bloc_limit"`
	BlocName       string `db:"bloc_name"`
	BlocNote       string `db:"bloc_note"`
	BlocType       string `db:"bloc_type"`
	Note           string `db:"note"`
	dbds.Course
	BlueprintSemesters pq.BoolArray `db:"semesters"`
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
	var dp degreePlanPage
	dp.degreePlanCode = records[0].DegreePlanCode
	dp.degreePlanYear = records[0].DegreePlanYear
	dp.canSave = false
	for _, record := range records {
		add(&dp, record)
	}
	fixLimits(&dp)
	return &dp, nil
}

func (m DBManager) degreePlan(uid, dpCode string, dpYear int, lang language.Language) (*degreePlanPage, error) {
	var records []dbDegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.DegreePlan, uid, dpCode, dpYear, lang); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DegreePlan: %w", err), errorx.P("dpCode", dpCode), errorx.P("dpYear", dpYear), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetDP,
		)
	}
	if len(records) == 0 {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(errors.New("no records found"), errorx.P("dpCode", dpCode), errorx.P("dpYear", dpYear), errorx.P("lang", lang)),
			http.StatusNotFound,
			texts[lang].errDPNotFound,
		)
	}
	var dp degreePlanPage
	dp.degreePlanCode = records[0].DegreePlanCode
	dp.degreePlanYear = dpYear
	dp.canSave = true
	for _, record := range records {
		add(&dp, record)
	}
	fixLimits(&dp)
	return &dp, nil
}

func (m DBManager) saveDegreePlan(uid, dpCode string, dpYear int, lang language.Language) error {
	_, err := m.DB.Exec(sqlquery.SaveDegreePlan, uid, dpCode, dpYear)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.SaveDegreePlan: %w", err), errorx.P("dpCode", dpCode), errorx.P("dpYear", dpYear), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotSaveDP,
		)
	}
	return nil
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
			note:         record.BlocNote,
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
		note:               from.Note,
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
