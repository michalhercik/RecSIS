package sqlquery

const InsertSemestersByYear = `--sql
WITH target_year_id AS (
        SELECT id FROM blueprint_years y
		WHERE y.user_id = $1
		AND y.id = $2
)
INSERT INTO blueprint_semesters(blueprint_year_id, semester)
VALUES
	((SELECT id from target_year_id), 1),
	((SELECT id from target_year_id), 2)
;
`
