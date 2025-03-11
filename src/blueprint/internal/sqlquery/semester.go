package sqlquery

const InsertSemestersByYear = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
target_year_id AS (
        SELECT id FROM blueprint_years y
        RIGHT JOIN user_session s ON y.user_id = s.user_id
		WHERE y.id = $2
)
INSERT INTO blueprint_semesters(blueprint_year_id, semester)
VALUES
	((SELECT id from target_year_id), 1),
	((SELECT id from target_year_id), 2)
;
`
