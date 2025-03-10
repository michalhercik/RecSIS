package coursedetail

import (
	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (reader DBManager) Course(code string, lang DBLang) (*Course, error) {
	var course Course
	if err := reader.DB.Get(&course, sqlquery.Course, code, lang); err != nil {
		return nil, err
	}

	return &course, nil
}