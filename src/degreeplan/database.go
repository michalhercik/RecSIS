package degreeplan

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/degreeplan/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) DegreePlan(uid string, lang language.Language) (*DegreePlan, error) {
	result, err := m.degreePlan(uid, lang)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m DBManager) degreePlan(uid string, lang language.Language) (*DegreePlan, error) {
	// Query Database
	fail := func(err error) (*DegreePlan, error) {
		return nil, fmt.Errorf("degreePlan: %v", err)
	}
	var records []DegreePlanRecord
	if err := m.DB.Select(&records, sqlquery.DegreePlan, uid, lang); err != nil {
		return fail(err)
	}
	var dp DegreePlan
	for _, record := range records {
		dp.add(record)
	}
	return &dp, nil
}
