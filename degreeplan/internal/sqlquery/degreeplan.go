package sqlquery

const DegreePlan = `
WITH user_session AS (
	SELECT user_id FROM sessions WHERE id=$1
),
user_blueprint_courses AS (
	SELECT course_code FROM user_session s
	LEFT JOIN blueprint_years y ON y.user_id = s.user_id
	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = y.id
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
)
SELECT
	dp.bloc_subject_code,
	COALESCE(dp.bloc_limit, -1) bloc_limit,
	COALESCE(dp.bloc_name, '') bloc_name,
	COALESCE(dp.bloc_note, '') bloc_note,
	COALESCE(dp.note, '') note,
	c.code,
	c.title,
	c.credits,
	c.start_semester,
	c.semester_count,
	c.exam_type,
	ubc.course_code IS NOT NULL in_blueprint
FROM user_session s
LEFT JOIN bla_studies bs ON s.user_id = bs.user_id
LEFT JOIN degree_plans dp ON bs.degree_plan_code = dp.plan_code AND bs.start_year = dp.plan_year
LEFT JOIN courses c ON dp.course_code = c.code
LEFT JOIN user_blueprint_courses ubc ON ubc.course_code = c.code
WHERE dp.lang = $2
AND c.lang = $2
ORDER BY dp.seq
;
`
