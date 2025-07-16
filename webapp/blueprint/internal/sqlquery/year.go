/*
Package sqlquery defines SQL queries for blueprint page.
*/
package sqlquery

/*
Select all positions of a specified student from blueprint_years table.

Params:

	$1 student
*/

const InsertYear = `--sql
WITH max_year AS (
	SELECT MAX(academic_year) AS max_academic_year
	FROM blueprint_years b
	WHERE b.user_id = $1
)
INSERT INTO blueprint_years (user_id, academic_year)
VALUES (
	$1,
	(SELECT max_academic_year + 1 FROM max_year)
)
RETURNING id;
`

const DeleteLastYear = `--sql
WITH max_academic_year AS (
	SELECT MAX(academic_year) AS academic_year FROM blueprint_years b
	WHERE b.user_id = $1
)
DELETE FROM blueprint_years b
USING max_academic_year m
WHERE b.user_id = $1
	AND m.academic_year != 0
	AND b.academic_year = m.academic_year;
`

const UnassignLastYear = `--sql
WITH max_academic_year AS (
	SELECT MAX(academic_year) AS academic_year FROM blueprint_years b
	WHERE b.user_id = $1
),
origin AS (
	SELECT bs.id, bs.semester
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	WHERE y.user_id = $1
		AND y.academic_year = (SELECT academic_year FROM max_academic_year)
),
unassigned AS (
	SELECT bs.id, COALESCE(position, 0) AS max_position
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
	WHERE y.user_id = $1
		AND y.academic_year = 0
		AND bs.semester = 0
	ORDER BY bc.position DESC
	LIMIT 1
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = u.id,
	position = u.max_position + (o.semester * bc.position)
FROM unassigned u, origin o
WHERE o.id = bc.blueprint_semester_id;
`
