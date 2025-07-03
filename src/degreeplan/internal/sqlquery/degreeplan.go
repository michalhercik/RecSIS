package sqlquery

const UserDegreePlan = `--sql
WITH user_start_year AS (
	-- TODO: pick a user selected study plan or max
	SELECT MIN(start_year) AS year FROM studies WHERE user_id = $1
),
user_blueprint_semesters AS (
	SELECT
		dp.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM studies s
	INNER JOIN user_start_year my
		ON s.start_year = my.year
	LEFT JOIN degree_plans dp
		ON s.degree_plan_code = dp.plan_code
		AND s.start_year = dp.plan_year
	LEFT JOIN blueprint_years by
		ON by.user_id = s.user_id
	LEFT JOIN blueprint_semesters bs
		ON by.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
		AND bc.course_code = dp.course_code
	WHERE s.user_id = $1
		AND dp.lang = $2
	GROUP BY dp.course_code, dp.bloc_subject_code
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
	c.lecture_range_winter,
	c.lecture_range_summer,
	c.seminar_range_winter,
	c.seminar_range_summer,
	c.exam,
	c.guarantors,
	ubs.semesters,
	CASE
		WHEN dp.bloc_type = 'A' THEN TRUE
		WHEN dp.bloc_type = 'B' THEN FALSE
	END AS is_compulsory
FROM studies s
INNER JOIN user_start_year my
	ON s.start_year = my.year
LEFT JOIN degree_plans dp
	ON s.degree_plan_code = dp.plan_code
	AND s.start_year = dp.plan_year
LEFT JOIN courses c
	ON dp.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON dp.course_code = ubs.course_code
WHERE s.user_id = $1
	AND dp.lang = $2
	AND c.lang = $2
	AND interchangeability IS NULL
	AND s.start_year = my.year
ORDER BY dp.bloc_type, dp.seq;
`

const DegreePlan = `
WITH user_start_year AS (
	-- TODO: pick a user selected study plan or max
	SELECT MIN(start_year) AS year FROM studies WHERE user_id = $1
),
user_blueprint_semesters AS (
	SELECT
		dp.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM studies s
	INNER JOIN user_start_year my
		ON s.start_year = my.year
	LEFT JOIN degree_plans dp
		ON s.degree_plan_code = dp.plan_code
		AND s.start_year = dp.plan_year
	LEFT JOIN blueprint_years by
		ON by.user_id = s.user_id
	LEFT JOIN blueprint_semesters bs
		ON by.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
		AND bc.course_code = dp.course_code
	WHERE s.user_id = $1
		AND dp.lang = $2
	GROUP BY dp.course_code, dp.bloc_subject_code
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
	c.lecture_range_winter,
	c.lecture_range_summer,
	c.seminar_range_winter,
	c.seminar_range_summer,
	c.exam,
	c.guarantors,
	ubs.semesters,
	CASE
		WHEN dp.bloc_type = 'A' THEN TRUE
		WHEN dp.bloc_type = 'B' THEN FALSE
	END AS is_compulsory
FROM degree_plans dp
LEFT JOIN courses c
	ON dp.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON dp.course_code = ubs.course_code
WHERE dp.plan_code = $2
	AND dp.plan_year = $3
	AND dp.lang = $4
	AND c.lang = $4
	AND interchangeability IS NULL
-- TODO: pick a user selected study or max
ORDER BY dp.bloc_type, dp.seq;
`
