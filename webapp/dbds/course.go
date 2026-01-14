package dbds

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Course struct {
	ID                     int               `db:"id"`
	Code                   string            `db:"code"`
	Title                  string            `db:"title"`
	Start                  int               `db:"start_semester"`
	URL                    sql.NullString    `db:"course_url"`
	LectureRangeWinter     sql.NullInt64     `db:"lecture_range_winter"`
	SeminarRangeWinter     sql.NullInt64     `db:"seminar_range_winter"`
	LectureRangeSummer     sql.NullInt64     `db:"lecture_range_summer"`
	SeminarRangeSummer     sql.NullInt64     `db:"seminar_range_summer"`
	RangeUnit              NullRangeUnit     `db:"range_unit"`
	ExamType               string            `db:"exam"`
	Credits                int               `db:"credits"`
	Guarantors             TeacherSlice      `db:"guarantors"`
	Faculty                Faculty           `db:"faculty"`
	Department             Department        `db:"department"`
	Teachers               TeacherSlice      `db:"teachers"`
	State                  string            `db:"taught_state_title"`
	Language               sql.NullString    `db:"taught_lang"`
	MinOccupancy           sql.NullString    `db:"min_occupancy"`
	MaxOccupancy           sql.NullString    `db:"capacity"`
	Classes                JSONArray[string] `db:"classes"`
	Classifications        JSONArray[string] `db:"classifications"`
	Annotation             NullDescription   `db:"annotation"`
	Syllabus               NullDescription   `db:"syllabus"`
	PassingTerms           NullDescription   `db:"terms_of_passing"`
	Literature             NullDescription   `db:"literature"`
	AssessmentRequirements NullDescription   `db:"requirements_of_assesment"`
	EntryRequirements      NullDescription   `db:"entry_requirements"`
	Aim                    NullDescription   `db:"aim"`
}

type NullRangeUnit struct {
	RangeUnit
	Valid bool
}

func (r *NullRangeUnit) Scan(val any) error {
	if val == nil {
		r.Valid = false
		return nil
	}
	if err := r.RangeUnit.Scan(val); err != nil {
		return err
	}
	r.Valid = true
	return nil
}

type RangeUnit struct {
	Abbr string `json:"abbr"`
	Name string `json:"name"`
}

func (r *RangeUnit) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &r)
		return nil
	case string:
		json.Unmarshal([]byte(v), &r)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type Faculty struct {
	Abbr string `json:"abbr"`
	Name string `json:"name"`
}

func (f *Faculty) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &f)
		return nil
	case string:
		json.Unmarshal([]byte(v), &f)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (d *Department) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &d)
		if err != nil {
			return fmt.Errorf("error unmarshalling Department: %w", err)
		}
		return nil
	case string:
		err := json.Unmarshal([]byte(v), &d)
		if err != nil {
			return fmt.Errorf("error unmarshalling Department: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type JSONArray[T any] []T

func (jsa *JSONArray[T]) Scan(val any) error {
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

type Description struct {
	Title   string `json:"title"`
	Content string `json:"content"`
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

type Requisite struct {
	Parent string         `db:"parent_course"`
	Child  string         `db:"child_course"`
	Type   string         `db:"req_type"`
	Group  sql.NullString `db:"group_type"`
}
