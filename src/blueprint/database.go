package blueprint

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/blueprint/internal/sqlquery"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/language"
)

// public interface for the blueprint database manager
type DBManager struct {
	DB *sqlx.DB
}

type dbBlueprintRecord struct {
	dbds.Course
	blueprintRecordPosition
}

type blueprintRecordPosition struct {
	AcademicYear int                `db:"academic_year"`
	Semester     semesterAssignment `db:"semester"`
	Folded       bool               `db:"folded"`
}

func (m DBManager) blueprint(userID string, lang language.Language) (*blueprintPage, error) {
	var records []dbBlueprintRecord
	var semestersInfo []blueprintRecordPosition

	tx, err := m.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := tx.Select(&records, sqlquery.SelectCourses, userID, lang); err != nil {
		return nil, err
	}
	if err := tx.Select(&semestersInfo, sqlquery.SelectSemestersInfo, userID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var bp blueprintPage
	if err := makeBlueprintPage(&bp, semestersInfo); err != nil {
		return nil, fmt.Errorf("make blueprint page: %w", err)
	}
	for _, record := range records {
		if err := add(&bp, record); err != nil {
			return nil, err
		}
	}
	return &bp, nil
}

func makeBlueprintPage(bp *blueprintPage, semestersInfo []blueprintRecordPosition) error {
	if len(semestersInfo) == 0 {
		return fmt.Errorf("no years found for user")
	}
	// make unnassigned semester
	bp.unassigned = semester{
		folded: semestersInfo[0].Folded,
	}
	// make years and semesters
	for i := 1; i < len(semestersInfo); i += 2 {
		bp.years = append(bp.years, academicYear{
			position: semestersInfo[i].AcademicYear,
			winter: semester{
				folded: semestersInfo[i].Folded,
			},
			summer: semester{
				folded: semestersInfo[i+1].Folded,
			},
		})
	}

	return nil
}

func add(bp *blueprintPage, record dbBlueprintRecord) error {
	if record.AcademicYear < 0 {
		return fmt.Errorf("year must be non-negative %d", record.AcademicYear)
	}
	var semester *semester
	switch record.Semester {
	case assignmentWinter:
		semester = &bp.years[record.AcademicYear-1].winter
	case assignmentSummer:
		semester = &bp.years[record.AcademicYear-1].summer
	case assignmentNone:
		semester = &bp.unassigned
	default:
		return fmt.Errorf("unknown semester assignment %d", record.Semester)
	}
	semester.courses = append(semester.courses, intoCourse(&record))
	return nil
}

func intoCourse(from *dbBlueprintRecord) course {
	return course{
		id:                 from.ID,
		code:               from.Code,
		title:              from.Title,
		semester:           teachingSemester(from.Start),
		lectureRangeWinter: from.LectureRangeWinter,
		seminarRangeWinter: from.SeminarRangeWinter,
		lectureRangeSummer: from.LectureRangeSummer,
		seminarRangeSummer: from.SeminarRangeSummer,
		examType:           from.ExamType,
		credits:            from.Credits,
		guarantors:         intoTeacherSlice(from.Guarantors),
	}
}

func intoTeacherSlice(from []dbds.Teacher) []teacher {
	teachers := make([]teacher, len(from))
	for i, t := range from {
		teachers[i] = teacher{
			sisID:       t.SISID,
			lastName:    t.LastName,
			firstName:   t.FirstName,
			titleBefore: t.TitleBefore,
			titleAfter:  t.TitleAfter,
		}
	}
	return teachers
}

// func (m DBManager) newCourse(userID string, course string, year int, semester semesterAssignment) (int, error) {
// 	row := m.DB.QueryRow(sqlquery.InsertCourse, userID, year, int(semester), course)
// 	var courseID int
// 	err := row.Scan(&courseID)
// 	return courseID, err
// }

func (m DBManager) insertCourses(userID string, year int, semester semesterAssignment, position int, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.MoveCourses, userID, pq.Array(courses), year, int(semester), position)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", count)
	}
	return nil
}

func (m DBManager) appendCourses(userID string, year int, semester semesterAssignment, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.AppendCourses, userID, year, int(semester), pq.Array(courses))
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) unassignYear(userID string, year int) error {
	_, err := m.DB.Exec(sqlquery.UnassignYear, userID, year)
	return err
}

func (m DBManager) unassignSemester(userID string, year int, semester semesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.UnassignSemester, userID, year, int(semester))
	return err
}

func (m DBManager) removeCourses(userID string, courses ...int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesByID, userID, pq.Array(courses))
	return err
}

func (m DBManager) removeCoursesBySemester(userID string, year int, semester semesterAssignment) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesBySemester, userID, year, int(semester))
	return err
}

func (m DBManager) removeCoursesByYear(userID string, year int) error {
	_, err := m.DB.Exec(sqlquery.DeleteCoursesByYear, userID, year)
	return err
}

func (m DBManager) addYear(userID string) error {
	fail := func(err error) error {
		return fmt.Errorf("AddYear: %v", err)
	}
	tx, err := m.DB.Beginx()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	var newYearID int
	err = tx.Get(&newYearID, sqlquery.InsertYear, userID)
	if err != nil {
		return fail(err)
	}
	_, err = tx.Exec(sqlquery.InsertSemestersByYear, userID, newYearID)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

func (m DBManager) removeYear(userID string, year int, shouldUnassign bool) error {
	fail := func(err error) error {
		return fmt.Errorf("RemoveYear: %v", err)
	}
	tx, err := m.DB.Beginx()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()
	if shouldUnassign {
		_, err := tx.Exec(sqlquery.UnassignYear, userID, year)
		if err != nil {
			return fail(err)
		}
	}
	_, err = tx.Exec(sqlquery.DeleteYear, userID)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

func (m DBManager) foldSemester(userID string, year int, semester semesterAssignment, folded bool) error {
	_, err := m.DB.Exec(sqlquery.FoldSemester, userID, year, int(semester), folded)
	if err != nil {
		return err
	}
	return nil
}
