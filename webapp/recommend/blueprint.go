package recommend

import (
	"fmt"
	"net/http"
	"github.com/michalhercik/RecSIS/errorx"

	"github.com/jmoiron/sqlx"
)

type BlueprintFetcher struct {
	DB *sqlx.DB
}

func (bf BlueprintFetcher) fetch(userID string) (string, error) {
	query := `--sql
		WITH course_per_semester AS (
			-- 1️⃣ Gather courses for each academic year + semester
			SELECT
				by.academic_year,
				bs.semester,
				jsonb_agg(bc.course_code) AS courses
			FROM blueprint_years by
			LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = by.id
			LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
			WHERE by.user_id = $1
			GROUP BY by.academic_year, bs.semester
		),
		semester_pivot AS (
			-- 2️⃣ Pivot semesters into JSON object keys
			SELECT
				academic_year,
				jsonb_object_agg(
					CASE
						WHEN semester = 0 THEN 'unassigned'
						WHEN semester = 1 THEN 'winter'
						WHEN semester = 2 THEN 'summer'
						ELSE semester::text
					END,
					COALESCE(
						NULLIF(courses, '[null]'::jsonb),  -- treat [null] as empty
						'[]'::jsonb
					)
				) AS semester_data
			FROM course_per_semester
			GROUP BY academic_year
		)

		-- 3️⃣ Final JSON array with {year, <semester>: [...]}
		SELECT jsonb_agg(
				jsonb_build_object('year', academic_year) || semester_data
				ORDER BY academic_year
			) AS result
		FROM semester_pivot;
	`
	var encodedBlueprint string
	err := bf.DB.QueryRow(query, userID).Scan(&encodedBlueprint)
	if err != nil {
		return "", errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("Cannot fetch blueprint: %w", err)),
			http.StatusInternalServerError,
			"",
		)
	}
	return encodedBlueprint, nil
}
