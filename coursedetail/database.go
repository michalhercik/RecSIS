package coursedetail

import (
	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (reader DBManager) Course(sessionID string, code string, lang DBLang) (*Course, error) {
	var course Course
	if err := reader.DB.Get(&course, sqlquery.Course, sessionID, code, lang); err != nil {
		return nil, err
	}
	return &course, nil
}

func (db DBManager) OverallRating(sessionID string, code string, value int) error {
	_, err := db.DB.Exec(sqlquery.OverallRating, sessionID, code, value)
	return err
}
