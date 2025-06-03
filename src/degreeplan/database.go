package degreeplan

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/degreeplan/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

type DBDegreePlanRecord struct {
	dbds.DegreePlanRecord
	BlueprintSemesters pq.BoolArray `db:"semesters"`
}

func (m DBManager) UserDegreePlan(uid string, lang language.Language) (*degreePlanPage, error) {
	var records []DBDegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.UserDegreePlan, uid, lang); err != nil {
		return nil, fmt.Errorf("degreePlan: %v", err)
	}
	var dp degreePlanPage
	for _, record := range records {
		add(&dp, record)
	}
	return &dp, nil
}

func (m DBManager) DegreePlan(uid, dpCode string, dpYear int, lang language.Language) (*degreePlanPage, error) {
	var records []DBDegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.DegreePlan, uid, dpCode, dpYear, lang); err != nil {
		return nil, fmt.Errorf("degreePlan: %v", err)
	}
	var dp degreePlanPage
	for _, record := range records {
		add(&dp, record)
	}
	return &dp, nil
}

func add(dp *degreePlanPage, record DBDegreePlanRecord) {
	blocIndex := -1
	for i, b := range dp.blocs {
		if b.Code == record.BlocCode {
			blocIndex = i
			break
		}
	}
	if blocIndex == -1 {
		// TODO: this is temporary FE hack, remove limit, if, and set Limit to record.BlocLimit
		limit := record.BlocLimit
		if record.IsBlocCompulsory {
			limit = 42
		}
		dp.blocs = append(dp.blocs, bloc{
			Name:         record.BlocName,
			Code:         record.BlocCode,
			Note:         record.BlocNote,
			Limit:        limit,
			IsCompulsory: record.IsBlocCompulsory,
		})
		blocIndex = len(dp.blocs) - 1
	}
	dp.blocs[blocIndex].Courses = append(dp.blocs[blocIndex].Courses, intoCourse(record))
}

func intoCourse(from DBDegreePlanRecord) course {
	return course{
		Code:               from.Code,
		Title:              from.Title,
		Credits:            from.Credits,
		Semester:           TeachingSemester(from.Start),
		LectureRangeWinter: from.LectureRangeWinter,
		SeminarRangeWinter: from.SeminarRangeWinter,
		LectureRangeSummer: from.LectureRangeSummer,
		SeminarRangeSummer: from.SeminarRangeSummer,
		ExamType:           from.ExamType,
		Guarantors:         intoTeacherSlice(from.Guarantors),
		Note:               from.Note,
		BlueprintSemesters: from.BlueprintSemesters,
	}
}

func intoTeacherSlice(from []dbds.Teacher) []Teacher {
	teachers := make([]Teacher, len(from))
	for i, t := range from {
		teachers[i] = Teacher{
			SISID:       t.SISID,
			LastName:    t.LastName,
			FirstName:   t.FirstName,
			TitleBefore: t.TitleBefore,
			TitleAfter:  t.TitleAfter,
		}
	}
	return teachers
}
