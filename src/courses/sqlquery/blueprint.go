package sqlquery

const Blueprint = `
	SELECT DISTINCT c.code, bs.semester, y.academic_year
	FROM sessions s
	LEFT JOIN blueprint_years y ON y.user_id = s.user_id
	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = y.id
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
	LEFT JOIN courses c ON c.code = bc.course_code AND c.valid_from = bc.course_valid_from
	WHERE s.id = $1
	AND c.code = ANY($2)
`
