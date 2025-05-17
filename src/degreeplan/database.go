package degreeplan

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/degreeplan/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) Course(uid, courseCode string, lang language.Language) (Course, error) {
	var course Course
	err := m.DB.Get(&course, sqlquery.Course, uid, courseCode, lang)
	if err != nil {
		return course, err
	}
	return course, nil
}

func (m DBManager) DegreePlan(uid string, lang language.Language) (*DegreePlan, error) {
	var records []dbds.DegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.DegreePlan, uid, lang); err != nil {
		return nil, fmt.Errorf("degreePlan: %v", err)
	}
	var dp DegreePlan
	for _, record := range records {
		add(&dp, record)
	}
	return &dp, nil
}

func add(dp *DegreePlan, record dbds.DegreePlanRecord) {
	blocIndex := -1
	for i, b := range dp.blocs {
		if b.Code == record.BlocCode {
			blocIndex = i
			break
		}
	}
	if blocIndex == -1 {
		dp.blocs = append(dp.blocs, Bloc{
			Name:         record.BlocName,
			Code:         record.BlocCode,
			Note:         record.BlocNote,
			Limit:        record.BlocLimit,
			IsCompulsory: record.IsBlocCompulsory,
		})
		blocIndex = len(dp.blocs) - 1
	}
	dp.blocs[blocIndex].Courses = append(dp.blocs[blocIndex].Courses, intoCourse(record))
}

func intoCourse(from dbds.DegreePlanRecord) Course {
	return Course{
		Code:           from.Code,
		Title:          from.Title,
		Credits:        from.Credits,
		Start:          TeachingSemester(from.Start),
		LectureRange1:  int(from.LectureRangeWinter.Int64),
		LectureRange2:  int(from.LectureRangeSummer.Int64),
		SeminarRange1:  int(from.SeminarRangeWinter.Int64),
		SeminarRange2:  int(from.SeminarRangeSummer.Int64),
		ExamType:       from.ExamType,
		Guarantors:     intoTeacherSlice(from.Guarantors),
		Note:           from.Note,
		BlueprintYears: []int64(from.BlueprintYears),
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
