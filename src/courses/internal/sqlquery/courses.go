package sqlquery

const Courses = `--sql
WITH user_blueprint_semesters AS (
	SELECT
		t.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM unnest($2::text[]) t(course_code)
	LEFT JOIN blueprint_years by
		ON by.user_id = $1
	LEFT JOIN blueprint_semesters bs
		ON by.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
		AND bc.course_code = t.course_code
	GROUP BY t.course_code
),
degree_plan AS (
	SELECT
		dp.course_code
	FROM bla_studies bs
	LEFT JOIN degree_plans dp
		ON dp.plan_code = bs.degree_plan_code
		AND dp.plan_year = bs.start_year
	WHERE bs.user_id = $1
		AND bs.start_year = (SELECT MIN(start_year) FROM bla_studies WHERE user_id = $1)
		AND dp.interchangeability IS NULL
		AND dp.lang = $3
)
SELECT
	c.code,
	c.title,
	c.annotation,
	COALESCE(c.start_semester, -1) start_semester,
	c.lecture_range1,
	c.seminar_range1,
	c.lecture_range2,
	c.seminar_range2,
	COALESCE(c.exam_type, '') exam_type,
	COALESCE(c.credits, -1) credits,
	c.guarantors,
	ubs.semesters,
	dp.course_code IS NOT NULL AS in_degree_plan
FROM courses c
JOIN unnest($2::text[]) WITH ORDINALITY t(id, ord)
	ON t.id = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON ubs.course_code = c.code
LEFT JOIN degree_plan dp
	ON c.code = dp.course_code
WHERE c.lang = $3
	AND c.valid_to = 9999
ORDER BY t.ord;
`
