package degreeplan

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/degreeplan/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) DegreePlan(uid string, lang DBLang) (*DegreePlan, error) {
	result, err := m.degreePlan(uid, lang)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m DBManager) degreePlan(uid string, lang DBLang) (*DegreePlan, error) {
	// Query Database
	row := m.DB.QueryRow(sqlquery.DegreePlan, uid, lang.String())
	var jsonData string
	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}
	// Parse result
	result := DegreePlan{}
	err = json.Unmarshal([]byte(jsonData), &result.blocs)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
