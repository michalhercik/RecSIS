package courses

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/courses/internal/filter"
	"github.com/michalhercik/RecSIS/courses/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) Courses(userID string, courseCodes []string, lang language.Language) ([]Course, error) {
	result := []Course{}
	if err := m.DB.Select(&result, sqlquery.Courses, userID, pq.Array(courseCodes), lang); err != nil {
		return nil, fmt.Errorf("failed to fetch courses: %w", err)
	}
	return result, nil
}

// func (m DBManager) ParamLabels(lang language.Language) (map[string][]filter.ParamValue, error) {
// 	var result map[string][]filter.ParamValue = make(map[string][]filter.ParamValue)
// 	var rows []struct {
// 		Param string `db:"param_name"`
// 		filter.ParamValue
// 	}
// 	if err := m.DB.Select(&rows, sqlquery.ParamLabels, lang); err != nil {
// 		return nil, fmt.Errorf("failed to fetch param labels: %w", err)
// 	}
// 	for _, row := range rows {
// 		result[row.Param] = append(result[row.Param], row.ParamValue)
// 	}
// 	return result, nil
// }

func (m DBManager) Filters() (filter.Filters, error) {
	// Retrieve
	tmpResult := []struct {
		CategoryID                  string         `db:"category_id"`
		CategoryFacetID             string         `db:"category_facet_id"`
		CategoryTitleCS             string         `db:"category_title_cs"`
		CategoryTitleEN             string         `db:"category_title_en"`
		CategoryDescCS              sql.NullString `db:"category_description_cs"`
		CategoryDescEN              sql.NullString `db:"category_description_en"`
		CategoryDisplayedValueLimit int            `db:"category_displayed_value_limit"`
		ValueID                     sql.NullString `db:"value_id"`
		ValueFacetID                sql.NullString `db:"value_facet_id"`
		ValueTitleCS                sql.NullString `db:"value_title_cs"`
		ValueTitleEN                sql.NullString `db:"value_title_en"`
		ValueDescCS                 sql.NullString `db:"value_description_cs"`
		ValueDescEN                 sql.NullString `db:"value_description_en"`
	}{}
	if err := m.DB.Select(&tmpResult, sqlquery.Filters); err != nil {
		return filter.Filters{}, fmt.Errorf("failed to fetch filters: %w", err)
	}
	// Parse
	fb := filter.FilterBuilder{}
	for _, row := range tmpResult {
		if fb.IsLastCategory(row.CategoryID) {
			fb.Category(filter.MakeFilterIdentity(
				row.CategoryID,
				row.CategoryFacetID,
				language.MakeLangString(row.CategoryTitleCS, row.CategoryTitleEN),
				language.MakeLangString(row.CategoryDescCS.String, row.CategoryDescEN.String),
			), row.CategoryDisplayedValueLimit)
		}
		if row.ValueID.Valid {
			fb.Value(filter.MakeFilterIdentity(
				row.ValueID.String,
				row.ValueFacetID.String,
				language.MakeLangString(row.ValueTitleCS.String, row.ValueTitleEN.String),
				language.MakeLangString(row.ValueDescCS.String, row.ValueDescEN.String),
			))
		}
	}
	return fb.Build(), nil
}
