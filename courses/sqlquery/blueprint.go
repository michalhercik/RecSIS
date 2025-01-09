package sqlquery

const Blueprint = `
	SELECT c.code, s.semester, y.position
	FROM blueprint_years AS y
	LEFT JOIN blueprint_semesters AS s ON s.blueprint_year = y.id
	LEFT JOIN courses AS c ON s.course = c.id
	WHERE y.student = $1
	AND c.code = ANY($2)
`
