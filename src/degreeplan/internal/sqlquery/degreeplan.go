package sqlquery

const DegreePlan = `
WITH user_blueprint_courses AS (
	SELECT course_code, ARRAY_AGG(y.academic_year) AS academic_years FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = y.id
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
	WHERE y.user_id = $1
	GROUP BY course_code
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
	COALESCE(c.lecture_range1, -1) lecture_range1,
	COALESCE(c.lecture_range2, -1) lecture_range2,
	COALESCE(c.seminar_range1, -1) seminar_range1,
	COALESCE(c.seminar_range2, -1) seminar_range2,
	-- c.semester_count,
	c.exam_type,
	c.guarantors,
	-- ubc.course_code IS NOT NULL in_blueprint
	ubc.academic_years
FROM bla_studies bs
LEFT JOIN degree_plans dp ON bs.degree_plan_code = dp.plan_code AND bs.start_year = dp.plan_year
LEFT JOIN courses c ON dp.course_code = c.code
LEFT JOIN user_blueprint_courses ubc ON ubc.course_code = c.code
WHERE bs.user_id = $1
AND dp.lang = $2
AND c.lang = $2
AND interchangeability IS NULL
-- TODO: pick a user selected study or max
AND bs.start_year = ( SELECT MIN(bs.start_year) FROM bla_studies bs WHERE bs.user_id = $1 )
ORDER BY dp.bloc_type, dp.seq
;
`

const Course = `
WITH user_blueprint_courses AS (
	SELECT course_code FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = y.id
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
	WHERE y.user_id = $1
)
SELECT
	c.code,
	c.title,
	c.credits,
	c.start_semester,
	COALESCE(c.lecture_range1, -1) lecture_range1,
	COALESCE(c.lecture_range2, -1) lecture_range2,
	COALESCE(c.seminar_range1, -1) seminar_range1,
	COALESCE(c.seminar_range2, -1) seminar_range2,
	c.semester_count,
	c.exam_type,
	ubc.course_code IS NOT NULL in_blueprint
FROM courses c
LEFT JOIN user_blueprint_courses ubc ON ubc.course_code = c.code
WHERE c.code = $2
AND c.lang = $3
;
`
