package recommend

import (
	"github.com/jmoiron/sqlx"
)

type NewCourses struct {
	DB *sqlx.DB
}

func (m NewCourses) Recommend(userID string) ([]string, error) {
	var courses []string
	// TODO:
	query := `--sql
		SELECT code
		FROM courses
		WHERE department->>'id' IN ('32-KSI', '32-UFAL', '32-KSVI', '32-KAM', '32-KTIML', '32-KDSS')
		AND lang = 'cs'
		AND taught_state = 'V'
		AND code NOT LIKE '%#%'
		AND code NOT LIKE '%$%'
		ORDER BY valid_from DESC
		LIMIT 10;
	`
	err := m.DB.Select(&courses, query)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	return courses, nil
}
