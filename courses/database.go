package courses

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/michalhercik/RecSIS/courses/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) Blueprint(user int, courses []string) (map[string][]Assignment, error) {
	rows, err := m.DB.Query(sqlquery.Blueprint, user, pq.Array(courses))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blueprint: %w", err)
	}
	defer rows.Close()
	res := map[string][]Assignment{}
	for rows.Next() {
		b := Assignment{}
		var code string
		if err := rows.Scan(&code, &b.year, &b.semester); err != nil {
			return nil, fmt.Errorf("failed to scan blueprint: %w", err)
		}
		courses, ok := res[code]
		if ok {
			res[code] = append(courses, b)
		} else {
			res[code] = []Assignment{b}
		}
	}
	return res, nil
}
