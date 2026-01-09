package sqlquery

const Courses = `--sql
	WITH degree_plan AS (
		SELECT
			dp.course_code
		FROM studies bs
		LEFT JOIN degree_plans dp
			ON dp.plan_code = bs.degree_plan_code
		WHERE bs.user_id = $1
			AND dp.interchangeability IS NULL
			AND dp.lang = $3
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
		dp.course_code IS NOT NULL AS in_degree_plan
	FROM courses c
	JOIN unnest($2::text[]) WITH ORDINALITY t(id, ord)
		ON t.id = c.code
	LEFT JOIN degree_plan dp
		ON c.code = dp.course_code
	WHERE c.lang = $3
		AND c.valid_to = 9999
	ORDER BY t.ord;
`
