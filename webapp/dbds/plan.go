package dbds

import (
	"database/sql"
)

type DegreePlan struct {
	Code                    string               `db:"code"`
	Title                   string               `db:"title"`
	ValidFrom               int                  `db:"valid_from"`
	ValidTo                 int                  `db:"valid_to"`
	FieldCode               string               `db:"field_code"`
	FieldTitle              string               `db:"field_title"`
	StudyType               string               `db:"study_type"`
	RequiredCredits         int                  `db:"required_credits"`
	RequiredElectiveCredits int                  `db:"required_elective_credits"`
	TotalCredits            int                  `db:"total_credits"`
	Studying                JSONArray[Studying]  `db:"studying"`
	Graduates               JSONArray[Graduates] `db:"graduates"`
	RequisiteGraphData      sql.NullString       `db:"requisite_graph_data"`
}

type Studying struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

type Graduates struct {
	Year     int     `json:"year"`
	Count    int     `json:"count"`
	AvgYears float64 `json:"avg_years"`
}
