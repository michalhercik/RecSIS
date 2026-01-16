package blueprint

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/blueprint/internal/sqlquery"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

const (
	uniqueViolationCode       = "23505"
	duplicateCoursesViolation = "blueprint_courses_blueprint_semester_id_course_code_key"
)

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
	t := texts[lang]
	var records []dbBlueprintRecord
	var semestersInfo []blueprintRecordPosition

	tx, err := m.DB.Beginx()
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("DB.Beginx: %w", err)),
			http.StatusInternalServerError,
			t.errCannotGetBlueprint,
		)
	}
	defer tx.Rollback()
	if err := tx.Select(&records, sqlquery.SelectCourses, userID, lang); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.SelectCourses: %w", err)),
			http.StatusInternalServerError,
			t.errCannotGetBlueprint,
		)
	}
	if err := tx.Select(&semestersInfo, sqlquery.SelectSemestersInfo, userID); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.SelectSemestersInfo: %w", err)),
			http.StatusInternalServerError,
			t.errCannotGetBlueprint,
		)
	}
	if err := tx.Commit(); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("tx.Commit: %w", err)),
			http.StatusInternalServerError,
			t.errCannotGetBlueprint,
		)
	}

	var bp blueprintPage
	if err := makeBlueprintPage(&bp, semestersInfo, lang); err != nil {
		return nil, errorx.AddContext(err)
	}
	for _, record := range records {
		if err := add(&bp, record, lang); err != nil {
			return nil, errorx.AddContext(err)
		}
	}
	return &bp, nil
}

func makeBlueprintPage(bp *blueprintPage, semestersInfo []blueprintRecordPosition, lang language.Language) error {
	if len(semestersInfo) == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("no semesters found in the database for user")),
			http.StatusInternalServerError,
			texts[lang].errNoSemestersFound,
		)
	}
	bp.unassigned = semester{
		folded: semestersInfo[0].Folded,
	}
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

func add(bp *blueprintPage, record dbBlueprintRecord, lang language.Language) error {
	if record.AcademicYear < 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid academic year %d", record.AcademicYear)),
			http.StatusInternalServerError,
			texts[lang].errInvalidYearInDB,
		)
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
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unknown semester assignment %d", record.Semester)),
			http.StatusInternalServerError,
			texts[lang].errInvalidSemesterInDB,
		)
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
		prerequisites:      intoRequisites(from.Prerequisites),
		corequisites:       intoRequisites(from.Corequisites),
		incompatibles:      intoRequisites(from.Incompatibilities),
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

func intoRequisites(from dbds.RequisiteSlice) requisiteSlice {
	result := make(requisiteSlice, len(from))
	for i, r := range from {
		result[i] = requisite{
			courseCode: r.CourseCode,
			children:   intoRequisites(r.Children),
			group:      r.Group,
		}
	}
	return result
}

func (m DBManager) moveCourses(userID string, lang language.Language, year int, semester semesterAssignment, position int, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.MoveCourses, userID, pq.Array(courses), year, int(semester), position)
	if err != nil {
		// Handle unique violation for blueprint_semester_id, course_code
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode && pqErr.Constraint == duplicateCoursesViolation {
			userErrMsg := texts[lang].errDuplicateCourseInBP
			if len(courses) > 1 {
				userErrMsg = texts[lang].errDuplicateCoursesInBP
			}
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.MoveCourses: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("position", position), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
				http.StatusConflict,
				userErrMsg,
			)
		}
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.MoveCourses: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("position", position), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
			http.StatusInternalServerError,
			texts[lang].errCannotMoveCourses,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.MoveCourses: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("position", position), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
			http.StatusBadRequest,
			texts[lang].errCannotMoveCourses,
		)
	}
	return nil
}

func (m DBManager) appendCourses(userID string, lang language.Language, year int, semester semesterAssignment, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.AppendCourses, userID, year, int(semester), pq.Array(courses))
	if err != nil {
		// Handle unique violation for blueprint_semester_id, course_code
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode && pqErr.Constraint == duplicateCoursesViolation {
			userErrMsg := texts[lang].errDuplicateCourseInBP
			if len(courses) > 1 {
				userErrMsg = texts[lang].errDuplicateCoursesInBP
			}
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.AppendCourses: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
				http.StatusConflict,
				userErrMsg,
			)
		}
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.AppendCourses: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
			http.StatusInternalServerError,
			texts[lang].errCannotAppendCourses,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.MoveCourses: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
			http.StatusBadRequest,
			texts[lang].errCannotAppendCourses,
		)
	}
	return nil
}

func (m DBManager) unassignSemester(userID string, lang language.Language, year int, semester semesterAssignment) error {
	res, err := m.DB.Exec(sqlquery.UnassignCoursesBySemester, userID, year, int(semester))
	if err != nil {
		// Handle unique violation for blueprint_semester_id, course_code
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode && pqErr.Constraint == duplicateCoursesViolation {
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.UnassignSemester: %w", err), errorx.P("year", year), errorx.P("semester", semester)),
				http.StatusConflict,
				texts[lang].errDuplicateCoursesInBPUnassigned,
			)
		}
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.UnassignSemester: %w", err), errorx.P("year", year), errorx.P("semester", semester)),
			http.StatusInternalServerError,
			texts[lang].errCannotUnassignSemester,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.UnassignSemester: %w", err), errorx.P("year", year), errorx.P("semester", semester)),
			http.StatusBadRequest,
			texts[lang].errCannotUnassignSemester,
		)
	}
	return nil
}

func (m DBManager) removeCourses(userID string, lang language.Language, courses ...int) error {
	res, err := m.DB.Exec(sqlquery.RemoveCoursesByID, userID, pq.Array(courses))
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.RemoveCoursesByID: %w", err), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
			http.StatusInternalServerError,
			texts[lang].errCannotRemoveCourses,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.RemoveCoursesByID: %w", err), errorx.P("courses", strings.Join(itoaSlice(courses), ","))),
			http.StatusBadRequest,
			texts[lang].errCannotRemoveCourses,
		)
	}
	return nil
}

func (m DBManager) removeCoursesBySemester(userID string, lang language.Language, year int, semester semesterAssignment) error {
	res, err := m.DB.Exec(sqlquery.RemoveCoursesBySemester, userID, year, int(semester))
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.RemoveCoursesBySemester: %w", err), errorx.P("year", year), errorx.P("semester", semester)),
			http.StatusInternalServerError,
			texts[lang].errCannotRemoveCourses,
		)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.RemoveCoursesBySemester: %w", err), errorx.P("year", year), errorx.P("semester", semester)),
			http.StatusBadRequest,
			texts[lang].errCannotRemoveCourses,
		)
	}
	return nil
}

func (m DBManager) addYear(userID string, lang language.Language) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("DB.Beginx: %w", err)),
			http.StatusInternalServerError,
			texts[lang].errCannotAddYear,
		)
	}
	defer tx.Rollback()
	var newYearID int
	err = tx.Get(&newYearID, sqlquery.InsertYear, userID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.InsertYear: %w", err)),
			http.StatusInternalServerError,
			texts[lang].errCannotAddYear,
		)
	}
	_, err = tx.Exec(sqlquery.InsertSemestersByYear, userID, newYearID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.InsertSemestersByYear: %w", err), errorx.P("year", newYearID)),
			http.StatusInternalServerError,
			texts[lang].errCannotAddYear,
		)
	}
	if err = tx.Commit(); err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("tx.Commit: %w", err)),
			http.StatusInternalServerError,
			texts[lang].errCannotAddYear,
		)
	}
	return nil
}

func (m DBManager) removeYear(userID string, lang language.Language, shouldUnassign bool) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("DB.Beginx: %w", err)),
			http.StatusInternalServerError,
			texts[lang].errCannotRemoveYear,
		)
	}
	defer tx.Rollback()
	if shouldUnassign {
		_, err := tx.Exec(sqlquery.UnassignLastYear, userID)
		if err != nil {
			// Handle unique violation for blueprint_semester_id, course_code
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode && pqErr.Constraint == duplicateCoursesViolation {
				return errorx.NewHTTPErr(
					errorx.AddContext(fmt.Errorf("sqlquery.UnassignLastYear: %w", err)),
					http.StatusConflict,
					texts[lang].errDuplicateCoursesInBPUnassigned,
				)
			}
			return errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("sqlquery.UnassignLastYear: %w", err)),
				http.StatusInternalServerError,
				texts[lang].errCannotUnassignYear,
			)
		}
	}
	_, err = tx.Exec(sqlquery.DeleteLastYear, userID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DeleteLastYear: %w", err)),
			http.StatusInternalServerError,
			texts[lang].errCannotRemoveYear,
		)
	}
	if err = tx.Commit(); err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("tx.Commit: %w", err)),
			http.StatusInternalServerError,
			texts[lang].errCannotRemoveYear,
		)
	}
	return nil
}

func (m DBManager) foldSemester(userID string, lang language.Language, year int, semester semesterAssignment, folded bool) error {
	_, err := m.DB.Exec(sqlquery.FoldSemester, userID, year, int(semester), folded)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.FoldSemester: %w", err), errorx.P("year", year), errorx.P("semester", semester), errorx.P("folded", folded)),
			http.StatusInternalServerError,
			texts[lang].errCannotUnFoldSemester,
		)
	}
	return nil
}

// transform int slice to string slice
func itoaSlice(ints []int) []string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strs
}
