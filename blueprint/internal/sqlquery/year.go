/*
Package sqlquery defines SQL queries for blueprint page.
*/
package sqlquery

/*
Select all positions of a specified student from blueprint_years table.

Params:

	$1 student
*/
const SelectYears string = `--sql
SELECT position
FROM blueprint_years
WHERE student=$1
ORDER BY position;
`

const InsertYear string = `
INSERT INTO blueprint_years (student, position)
VALUES ($1, (
	SELECT MAX(position) + 1
	FROM blueprint_years
	WHERE student=42
	)
);`

// TODO: Delete only if last year
const DeleteYearCourses = `--sql
DELETE FROM blueprint_semesters
WHERE blueprint_year = (
	SELECT id
	FROM blueprint_years
	WHERE student=$1
	AND position=(
		SELECT MAX(position)
		FROM blueprint_years
		WHERE student=$1
	)
);`

// TODO: Delete only if last year
const DeleteYear = `--sql
DELETE FROM blueprint_years
WHERE student=$1
AND position=(
	SELECT MAX(position)
	FROM blueprint_years
	WHERE student=$1
);`
