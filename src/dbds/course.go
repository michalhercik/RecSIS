package dbds

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Faculty struct {
	SisID int
	Name  string
	Abbr  string
}

type Course struct {
	// SemesterCount         int             `db:"semester_count"`
	Code                  string          `db:"code"`
	Title                 string          `db:"title"`
	Start                 int             `db:"start_semester"`
	LectureRangeWinter    sql.NullInt64   `db:"lecture_range1"`
	SeminarRangeWinter    sql.NullInt64   `db:"seminar_range1"`
	LectureRangeSummer    sql.NullInt64   `db:"lecture_range2"`
	SeminarRangeSummer    sql.NullInt64   `db:"seminar_range2"`
	ExamType              string          `db:"exam_type"`
	Credits               int             `db:"credits"`
	Guarantors            TeacherSlice    `db:"guarantors"`
	Faculty               string          `db:"faculty"`
	Department            string          `db:"guarantor"`
	Teachers              TeacherSlice    `db:"teachers"`
	State                 string          `db:"taught"`
	Language              sql.NullString  `db:"taught_lang"`
	MinOccupancy          sql.NullInt64   `db:"min_number"`
	MaxOccupancy          sql.NullString  `db:"capacity"`
	Prereq                JSONStringArray `db:"preqrequisities"`
	Coreq                 JSONStringArray `db:"corequisities"`
	Incompa               JSONStringArray `db:"incompatibilities"`
	Interchange           JSONStringArray `db:"interchangebilities"`
	Classes               ClassSlice      `db:"classes"`
	Classifications       ClassSlice      `db:"classifications"`
	Annotation            NullDescription `db:"annotation"`
	Syllabus              NullDescription `db:"syllabus"`
	PassingTerms          NullDescription `db:"terms_of_passing"`
	Literature            NullDescription `db:"literature"`
	AssesmentRequirements NullDescription `db:"requirements_for_assesment"`
	EntryRequirements     NullDescription `db:"entry_requirements"`
	Aim                   NullDescription `db:"aim"`
}

type OverallRating struct {
	UserRating sql.NullInt64   `db:"rating"`
	AvgRating  sql.NullFloat64 `db:"avg_rating"`
	Count      sql.NullInt64   `db:"rating_count"`
}

type CourseCategoryRating struct {
	Code        int             `db:"category_code"`
	Title       string          `db:"rating_title"`
	UserRating  sql.NullInt64   `db:"rating"`
	AvgRating   sql.NullFloat64 `db:"avg_rating"`
	RatingCount sql.NullInt64   `db:"rating_count"`
}

type SemesterAssignment int

func (sa *SemesterAssignment) Scan(val any) error {
	switch v := val.(type) {
	case int64:
		*sa = SemesterAssignment(v)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func SemesterAssignmentFromString(s string) (SemesterAssignment, error) {
	switch s {
	case winterStr:
		return assignmentWinter, nil
	case summerStr:
		return assignmentSummer, nil
	case unassignedStr:
		return assignmentNone, nil
	default:
		return 0, fmt.Errorf("unknown semester assignment %s", s)
	}
}

const (
	assignmentNone SemesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

const (
	winterStr     = "winter"
	summerStr     = "summer"
	unassignedStr = "unassigned"
)

type Description struct {
	Title   string `json:"TITLE"`
	Content string `json:"MEMO"`
}

func (d *Description) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &d)
		return nil
	case string:
		json.Unmarshal([]byte(v), &d)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (d Description) Value() (any, error) {
	return json.Marshal(d)
}

type NullDescription struct {
	Description
	Valid bool
}

func (d *NullDescription) Scan(val any) error {
	if val == nil {
		d.Valid = false
		return nil
	}
	if err := d.Description.Scan(val); err != nil {
		return err
	}
	d.Valid = true
	return nil
}

type JSONStringArray []string

func (jsa *JSONStringArray) Scan(val any) error {
	switch v := val.(type) {
	case nil:
		jsa = nil
		return nil
	case []byte:
		*jsa = nil
		err := json.Unmarshal(v, &jsa)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &jsa)
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type ClassSlice []Class

func (cs *ClassSlice) Scan(val any) error {
	switch v := val.(type) {
	case nil:
		*cs = nil
		return nil
	case []byte:
		*cs = nil
		err := json.Unmarshal(v, &cs)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &cs)
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type Class struct {
	Code string `json:"KOD"`
	Name string `json:"NAZEV"`
}
