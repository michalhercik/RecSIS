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
	SELECT DISTINCT(dpc.course_code)
	FROM studies bs
	LEFT JOIN degree_plan_courses dpc
		ON dpc.plan_code = bs.degree_plan_code
		AND dpc.lang = $3
	WHERE bs.user_id = $1
		AND dpc.interchangeability IS NULL
)
SELECT
	c.code,
	c.title,
	c.annotation,
	COALESCE(c.start_semester, '') start_semester,
	c.lecture_range_winter,
	c.seminar_range_winter,
	c.lecture_range_summer,
	c.seminar_range_summer,
	COALESCE(c.exam, '') exam,
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
