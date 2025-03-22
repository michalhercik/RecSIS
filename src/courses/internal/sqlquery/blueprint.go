package sqlquery

const ParamLabels = `--sql
SELECT fp.param_name, fl.id, fl.label FROM filter_params fp
LEFT JOIN filter_labels fl ON fp.value_id = fl.id
WHERE fl.lang = $1
ORDER BY fl.label
`

const Courses = `--sql
WITH user_session AS (
	SELECT user_id FROM sessions WHERE id=$1
),
user_blueprint_courses AS (
	SELECT
		bc.course_code,
		json_agg(json_build_object('year', y.academic_year, 'semester', bs.semester)) AS assignment
	FROM user_session s
	LEFT JOIN blueprint_years y ON y.user_id = s.user_id
	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = y.id
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
	GROUP BY bc.course_code
)
SELECT
	c.code,
	c.title,
	c.annotation,
	COALESCE(c.start_semester, -1) start_semester,
	COALESCE(c.semester_count, -1) semester_count,
	COALESCE(c.lecture_range1, -1) lecture_range1,
	COALESCE(c.seminar_range1, -1) seminar_range1,
	COALESCE(c.lecture_range2, -1) lecture_range2,
	COALESCE(c.seminar_range2, -1) seminar_range2,
	COALESCE(c.exam_type, '') exam_type,
	COALESCE(c.credits, -1) credits,
	c.guarantors,
	ubc.assignment
FROM courses c
JOIN unnest($2::text[]) WITH ORDINALITY t(id, ord) ON t.id = c.code
LEFT JOIN user_blueprint_courses ubc ON ubc.course_code = c.code
WHERE c.lang = $3
AND c.valid_to = 9999
;
`

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
