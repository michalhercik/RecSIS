/*
Package sqlquery defines SQL queries for blueprint page.
*/
package sqlquery

/*
Select all positions of a specified student from blueprint_years table.

Params:

	$1 student
*/
const SelectYears = `--sql
SELECT academic_year
FROM sessions s
LEFT JOIN blueprint_years b ON s.user_id=b.user_id
WHERE s.id=$1
ORDER BY academic_year;
`

const InsertYear = `
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
max_year AS (
	SELECT MAX(academic_year) AS max_academic_year FROM blueprint_years b
	RIGHT JOIN user_session u ON u.user_id=b.user_id
)
INSERT INTO blueprint_years (user_id, academic_year)
VALUES (
	(SELECT user_id FROM user_session),
	(SELECT max_academic_year + 1 FROM max_year)
)
RETURNING id;
;
`

// TODO: Delete only if last year
const DeleteYear = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
max_academic_year AS (
        SELECT MAX(academic_year) AS academic_year FROM blueprint_years b
        RIGHT JOIN user_session s ON b.user_id = s.user_id
)
DELETE FROM blueprint_years b
USING user_session u, max_academic_year m
WHERE b.user_id = u.user_id
AND b.academic_year = m.academic_year
;
`
