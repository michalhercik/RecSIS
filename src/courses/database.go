package courses

import (
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

func (m DBManager) Courses(sessionID string, courseCodes []string, lang language.Language) ([]Course, error) {
	result := []Course{}
	if err := m.DB.Select(&result, sqlquery.Courses, sessionID, pq.Array(courseCodes), lang); err != nil {
		return nil, fmt.Errorf("failed to fetch courses: %w", err)
	}
	return result, nil
}

func (m DBManager) ParamLabels(lang language.Language) (map[string][]filter.ParamValue, error) {
	var result map[string][]filter.ParamValue = make(map[string][]filter.ParamValue)
	var rows []struct {
		Param string `db:"param_name"`
		filter.ParamValue
	}
	if err := m.DB.Select(&rows, sqlquery.ParamLabels, lang); err != nil {
		return nil, fmt.Errorf("failed to fetch param labels: %w", err)
	}
	for _, row := range rows {
		result[row.Param] = append(result[row.Param], row.ParamValue)
	}
	return result, nil
}
