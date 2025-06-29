package home

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

// page content
type homePage struct {
	recommendedCourses []course
	newCourses         []course
}

// course data structure
type course struct {
	Code               string           `json:"code"`
	Title              string           `json:"title"`
	Semester           teachingSemester `json:"semester"`
	LectureRangeWinter sql.NullInt64    `json:"lecture_range_winter"`
	SeminarRangeWinter sql.NullInt64    `json:"seminar_range_winter"`
	LectureRangeSummer sql.NullInt64    `json:"lecture_range_summer"`
	SeminarRangeSummer sql.NullInt64    `json:"seminar_range_summer"`
	ExamType           string           `json:"exam"`
	Credits            int              `json:"credits"`
	Guarantors         teacherSlice     `json:"guarantors"`
	//BlueprintAssignments AssignmentSlice
	//Annotation           NullDescription
}

// semester types and methods
type teachingSemester int

const (
	teachingWinterOnly teachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts *teachingSemester) string(t text) string {
	switch *ts {
	case teachingWinterOnly:
		return t.winter
	case teachingSummerOnly:
		return t.summer
	case teachingBoth:
		return t.both
	default:
		return "unsupported"
	}
}

// teacher types and methods
type teacher struct {
	SisID       int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TitleBefore string `json:"title_before"`
	TitleAfter  string `json:"title_after"`
}

func (t teacher) string() string {
	firstRune, _ := utf8.DecodeRuneInString(t.FirstName)
	return fmt.Sprintf("%c. %s", firstRune, t.LastName)
}

type teacherSlice []teacher

func (ts teacherSlice) string(t text) string {
	names := []string{}
	for _, teacher := range ts {
		names = append(names, teacher.string())
	}
	if len(names) == 0 {
		return t.noGuarantors
	}
	return strings.Join(names, ", ")
}

func (ts *teacherSlice) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &ts)
		return nil
	case string:
		json.Unmarshal([]byte(v), &ts)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

// ================================================================================
// Helper Functions
// ================================================================================

func hoursString(course *course, t text) string {
	result := ""
	winter := course.LectureRangeWinter.Valid && course.SeminarRangeWinter.Valid
	summer := course.LectureRangeSummer.Valid && course.SeminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", course.LectureRangeWinter.Int64, course.SeminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", course.LectureRangeSummer.Int64, course.SeminarRangeSummer.Int64)
	}
	return result
}
