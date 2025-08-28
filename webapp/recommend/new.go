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
		AND code NOT IN (
			SELECT bc.course_code 
			FROM blueprint_years by, blueprint_semesters bs, blueprint_courses bc 
			WHERE by.id = bs.blueprint_year_id
			AND bs.id = bc.blueprint_semester_id
			AND by.user_id = $1
		)
		ORDER BY valid_from DESC
		LIMIT 30;
	`
	err := m.DB.Select(&courses, query, userID)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	selected := chooseRandom(courses, 10)
	return selected, nil
}
