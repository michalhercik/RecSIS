package courses

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
