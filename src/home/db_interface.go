package home

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/michalhercik/RecSIS/language"
)

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

// page content
type HomePage struct {
	RecommendedCourses []Course
	NewCourses         []Course
}

// course data structure
type Course struct {
	Code               string           `json:"code"`
	Title              string           `json:"title"`
	Semester           TeachingSemester `json:"semester"`
	LectureRangeWinter sql.NullInt64    `json:"lecture_range_winter"`
	SeminarRangeWinter sql.NullInt64    `json:"seminar_range_winter"`
	LectureRangeSummer sql.NullInt64    `json:"lecture_range_summer"`
	SeminarRangeSummer sql.NullInt64    `json:"seminar_range_summer"`
	ExamType           string           `json:"exam"`
	Credits            int              `json:"credits"`
	Guarantors         TeacherSlice     `json:"guarantors"`
	//BlueprintAssignments AssignmentSlice
	//Annotation           NullDescription
}

// semester types and methods
type TeachingSemester int

const (
	teachingWinterOnly TeachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts *TeachingSemester) String(lang language.Language) string {
	t := texts[lang]
	semester := ""
	switch *ts {
	case teachingWinterOnly:
		semester = t.Winter
	case teachingSummerOnly:
		semester = t.Summer
	case teachingBoth:
		semester = t.Both
	default:
		semester = "unsupported"
	}
	return semester
}

// teacher types and methods
type Teacher struct {
	SisID       int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TitleBefore string `json:"title_before"`
	TitleAfter  string `json:"title_after"`
}

func (t Teacher) String() string {
	firstRune, _ := utf8.DecodeRuneInString(t.FirstName)
	return fmt.Sprintf("%c. %s", firstRune, t.LastName)
}

type TeacherSlice []Teacher

func (ts TeacherSlice) string(t text) string {
	names := []string{}
	for _, teacher := range ts {
		names = append(names, teacher.String())
	}
	if len(names) == 0 {
		return t.NoGuarantors
	}
	return strings.Join(names, ", ")
}

func (ts *TeacherSlice) Scan(val interface{}) error {
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
