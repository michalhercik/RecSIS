package sqlquery

const SelectYears = `
SELECT position
FROM blueprint_years
WHERE student=$1
ORDER BY position;
`

const InsertYear = `
INSERT INTO blueprint_years (student, position)
VALUES ($1, (
	SELECT MAX(position) + 1 
	FROM blueprint_years 
	WHERE student=42
	)
);
`

// TODO: Delete only if last year
const DeleteYearCourses = `
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
const DeleteYear = `
DELETE FROM blueprint_years 
WHERE student=$1
AND position=(
	SELECT MAX(position)
	FROM blueprint_years
	WHERE student=$1
);`
