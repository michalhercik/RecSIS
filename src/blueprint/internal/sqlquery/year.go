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
SELECT academic_year FROM blueprint_years b
WHERE b.user_id=$1
ORDER BY academic_year;
`

const InsertYear = `--sql
WITH max_year AS (
	SELECT MAX(academic_year) AS max_academic_year FROM blueprint_years b
	WHERE b.user_id = $1
)
INSERT INTO blueprint_years (user_id, academic_year)
VALUES (
	$1,
	(SELECT max_academic_year + 1 FROM max_year)
)
RETURNING id;
;
`

// TODO: Delete only if last year
const DeleteYear = `--sql
WITH max_academic_year AS (
        SELECT MAX(academic_year) AS academic_year FROM blueprint_years b
        WHERE b.user_id = $1
)
DELETE FROM blueprint_years b
USING max_academic_year m
WHERE b.user_id = $1
AND b.academic_year = m.academic_year
;
`

const FoldSemester = `--sql
UPDATE blueprint_semesters bs
SET folded = $4
FROM blueprint_years BY
WHERE bs.blueprint_year_id = by.id
AND by.user_id = $1
AND by.academic_year = $2
AND bs.semester = $3
`
