package sqlquery

const DegreePlan = `--sql
SELECT
	dp.plan_code AS degree_plan_code,
	dp.title AS degree_plan_title,
	-- dp.valid_from AS degree_plan_valid_from,
	-- dp.valid_to AS degree_plan_valid_to,
	dpc.bloc_subject_code,
	dpc.bloc_name,
	COALESCE(dpc.bloc_limit, -1) AS bloc_limit,
	dpc.is_required,
	dpc.is_elective,
	c.code,
	c.title,
	COALESCE(c.credits, 0) AS credits,
	c.credits IS NOT NULL AS course_is_supported
FROM degree_plans dp
LEFT JOIN degree_plan_courses dpc
	ON dp.plan_code = dpc.plan_code
	AND dpc.lang = $2
LEFT JOIN courses c
	ON dpc.course_code = c.code
	AND c.lang = $2
WHERE dp.plan_code = $1
	AND dp.lang = $2
	AND interchangeability IS NULL
ORDER BY dpc.is_required DESC, (dpc.is_elective = false) DESC, dpc.seq, c.code;
`
