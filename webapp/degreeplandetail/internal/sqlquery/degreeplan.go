package sqlquery

const UserDegreePlanCode = `--sql
SELECT degree_plan_code
FROM studies
WHERE user_id = $1
`

const UserDegreePlan = `--sql
WITH user_blueprint_semesters AS (
	SELECT DISTINCT
		dp.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM studies s
	LEFT JOIN degree_plans dp
		ON s.degree_plan_code = dp.plan_code
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
	s.degree_plan_code,
	dp.bloc_subject_code,
	COALESCE(dp.bloc_limit, -1) bloc_limit,
	COALESCE(dp.bloc_name, '') bloc_name,
	dp.bloc_type,
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
	ubs.semesters
FROM studies s
LEFT JOIN degree_plans dp
	ON s.degree_plan_code = dp.plan_code
LEFT JOIN courses c
	ON dp.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON dp.course_code = ubs.course_code
WHERE s.user_id = $1
	AND dp.lang = $2
	AND c.lang = $2
	AND interchangeability IS NULL
ORDER BY dp.bloc_type, dp.seq;
`

const DegreePlan = `--sql
WITH user_blueprint_semesters AS (
	SELECT DISTINCT
		dp.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM degree_plans dp
	LEFT JOIN blueprint_years by
		ON by.user_id = $1
	LEFT JOIN blueprint_semesters bs
		ON by.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
		AND bc.course_code = dp.course_code
	WHERE dp.plan_code = $2
		AND dp.lang = $3
	GROUP BY dp.course_code, dp.bloc_subject_code
)
SELECT
	dp.plan_code AS degree_plan_code,
	dp.bloc_subject_code,
	COALESCE(dp.bloc_limit, -1) bloc_limit,
	COALESCE(dp.bloc_name, '') bloc_name,
	dp.bloc_type,
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
	ubs.semesters
FROM degree_plans dp
LEFT JOIN courses c
	ON dp.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON dp.course_code = ubs.course_code
WHERE dp.plan_code = $2
	AND dp.lang = $3
	AND c.lang = $3
	AND interchangeability IS NULL
ORDER BY dp.bloc_type, dp.seq;
`

const SaveDegreePlan = `--sql
UPDATE studies
SET degree_plan_code = $2
WHERE user_id = $1
`

const DeleteSavedDegreePlan = `--sql
UPDATE studies
SET degree_plan_code = NULL
WHERE user_id = $1
`
