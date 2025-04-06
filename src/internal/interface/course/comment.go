package course

import (
	"encoding/json"
	"strconv"

	"github.com/michalhercik/RecSIS/internal/interface/teacher"
)

type Comment struct {
	Student
	CommentTarget
	AcademicYear int    `json:"academic_year"`
	Semester     string `json:"semester"`
	Content      string `json:"content"`
}

func (c Comment) AcademicYearString() string {
	return strconv.Itoa(c.AcademicYear)
}

type CommentTarget struct {
	Type          string          `json:"target_type"` // Lecture or Seminar
	CourseCode    string          `json:"course_code"`
	TargetTeacher teacher.Teacher `json:"teacher"`
}

type Student struct {
	StudyYear  int       `json:"study_year"`
	StudyField string    `json:"study_field"`
	Study      StudyType `json:"study_type"`
}

func (c Comment) StudiesYearString() string {
	return strconv.Itoa(c.StudyYear)
}

type StudyType struct {
	Code string `json:"code"`
	Abbr string `json:"abbr"`
	Name string `json:"name"`
}

func (st *StudyType) UnmarshalJSON(val []byte) error {
	var tmp struct {
		Code   string `json:"code"`
		Abbr   string `json:"abbr"`
		NameCs string `json:"name_cs"`
		NameEn string `json:"name_en"`
	}
	if err := json.Unmarshal(val, &tmp); err != nil {
		return err
	}
	st.Code = tmp.Code
	st.Abbr = tmp.Abbr
	st.Name = tmp.NameCs
	if len(st.Name) == 0 {
		st.Name = tmp.NameEn
	}
	return nil
}
