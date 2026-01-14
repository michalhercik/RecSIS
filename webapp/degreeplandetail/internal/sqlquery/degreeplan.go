package sqlquery

const UserDegreePlanCode = `--sql
SELECT degree_plan_code
FROM studies
WHERE user_id = $1
`

const UserDegreePlan = `--sql
WITH user_blueprint_semesters AS (
	SELECT DISTINCT
		dpc.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM studies s
	LEFT JOIN degree_plan_courses dpc
		ON s.degree_plan_code = dpc.plan_code
	LEFT JOIN blueprint_years by
		ON by.user_id = s.user_id
	LEFT JOIN blueprint_semesters bs
		ON by.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
		AND bc.course_code = dpc.course_code
	WHERE s.user_id = $1
		AND dpc.lang = $2
	GROUP BY dpc.course_code, dpc.bloc_subject_code
)
SELECT
	s.degree_plan_code,
	dpc.bloc_subject_code,
	COALESCE(dpc.bloc_limit, -1) bloc_limit,
	COALESCE(dpc.bloc_name, '') bloc_name,
	dpc.bloc_type,
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
LEFT JOIN degree_plan_courses dpc
	ON s.degree_plan_code = dpc.plan_code
LEFT JOIN courses c
	ON dpc.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON dpc.course_code = ubs.course_code
WHERE s.user_id = $1
	AND dpc.lang = $2
	AND c.lang = $2
	AND interchangeability IS NULL
ORDER BY dpc.bloc_type, dpc.seq;
`

const DegreePlan = `--sql
WITH user_blueprint_semesters AS (
	SELECT DISTINCT
		dpc.course_code,
		array_agg(bc.course_code IS NOT NULL ORDER BY by.academic_year, bs.semester) AS semesters
	FROM degree_plan_courses dpc
	LEFT JOIN blueprint_years by
		ON by.user_id = $1
	LEFT JOIN blueprint_semesters bs
		ON by.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
		AND bc.course_code = dpc.course_code
	WHERE dpc.plan_code = $2
		AND dpc.lang = $3
	GROUP BY dpc.course_code, dpc.bloc_subject_code
)
SELECT
	dpc.plan_code AS degree_plan_code,
	dpc.bloc_subject_code,
	COALESCE(dpc.bloc_limit, -1) bloc_limit,
	COALESCE(dpc.bloc_name, '') bloc_name,
	dpc.bloc_type,
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
FROM degree_plan_courses dpc
LEFT JOIN courses c
	ON dpc.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs
	ON dpc.course_code = ubs.course_code
WHERE dpc.plan_code = $2
	AND dpc.lang = $3
	AND c.lang = $3
	AND interchangeability IS NULL
ORDER BY dpc.bloc_type, dpc.seq;
`

const SaveDegreePlan = `--sql
UPDATE studies
SET degree_plan_code = $2
WHERE user_id = $1
`

const DeleteSavedDegreePlan = `--sql
WITH old_plan AS (
	SELECT degree_plan_code
	FROM studies
	WHERE user_id = $1
)
UPDATE studies
SET degree_plan_code = NULL
WHERE user_id = $1
RETURNING (SELECT degree_plan_code FROM old_plan)
`
