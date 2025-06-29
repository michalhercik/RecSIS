package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Course struct {
	Code               string        `db:"code" json:"code"`
	Title              string        `db:"title" json:"title"`
	Semester           int           `db:"start_semester" json:"semester"`
	LectureRangeWinter sql.NullInt64 `db:"lecture_range_winter" json:"lecture_range_winter"`
	SeminarRangeWinter sql.NullInt64 `db:"seminar_range_winter" json:"seminar_range_winter"`
	LectureRangeSummer sql.NullInt64 `db:"lecture_range_summer" json:"lecture_range_summer"`
	SeminarRangeSummer sql.NullInt64 `db:"seminar_range_summer" json:"seminar_range_summer"`
	ExamType           string        `db:"exam" json:"exam"`
	Credits            int           `db:"credits" json:"credits"`
	Guarantors         TeacherSlice  `db:"guarantors" json:"guarantors"`
}

type TeacherSlice []Teacher

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

type Teacher struct {
	SISID       string `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	TitleBefore string `json:"title_before"`
	TitleAfter  string `json:"title_after"`
}

func getCourses(w http.ResponseWriter, r *http.Request, query string) {
	var courses []Course
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "cs" // Default language
	}
	err := db.Select(&courses, query, lang)
	if err != nil {
		log.Printf("DB error: %v", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

const sql_recommended = `
SELECT
    c.code,
    c.title,
    c.start_semester,
    c.lecture_range_winter,
    c.seminar_range_winter,
    c.lecture_range_summer,
    c.seminar_range_summer,
    c.exam,
    c.credits,
    c.guarantors
FROM courses c
WHERE c.lang = $1
AND c.start_semester IS NOT NULL
LIMIT 20;
`

const sql_newest = `
SELECT
    c.code,
    c.title,
    c.start_semester,
    c.lecture_range_winter,
    c.seminar_range_winter,
    c.lecture_range_summer,
    c.seminar_range_summer,
    c.exam,
    c.credits,
    c.guarantors
FROM courses c
WHERE c.lang = $1
AND c.valid_from IS NOT NULL
AND c.start_semester IS NOT NULL
ORDER BY c.valid_from DESC
LIMIT 20;
`
