package sqlquery

const DegreePlan = `
WITH user_blueprint_semesters AS (
        SELECT dp.course_code, array_agg(bc.course_code IS NOT NULL) AS semesters
        FROM bla_studies s
        LEFT JOIN degree_plans dp ON s.degree_plan_code = dp.plan_code AND s.start_year = dp.plan_year
        LEFT JOIN blueprint_years by ON by.user_id=s.user_id
        LEFT JOIN blueprint_semesters bs ON by.id = bs.blueprint_year_id
        LEFT JOIN blueprint_courses bc ON bs.id = bc.blueprint_semester_id AND bc.course_code=dp.course_code
        WHERE s.user_id=$1
        AND dp.lang=$2
        AND s.start_year = ( SELECT MIN(s.start_year) FROM bla_studies s WHERE s.user_id = $1 )
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
	COALESCE(c.lecture_range1, -1) lecture_range1,
	COALESCE(c.lecture_range2, -1) lecture_range2,
	COALESCE(c.seminar_range1, -1) seminar_range1,
	COALESCE(c.seminar_range2, -1) seminar_range2,
	c.exam_type,
	c.guarantors,
        ubs.semesters,
	CASE WHEN dp.bloc_type = 'A' THEN TRUE WHEN dp.bloc_type = 'B' THEN FALSE END AS is_compulsory
FROM bla_studies s
LEFT JOIN degree_plans dp ON s.degree_plan_code = dp.plan_code AND s.start_year = dp.plan_year
LEFT JOIN courses c ON dp.course_code = c.code
LEFT JOIN user_blueprint_semesters ubs ON dp.course_code = ubs.course_code

WHERE s.user_id = $1
AND dp.lang = $2
AND c.lang = $2
AND interchangeability IS NULL
-- TODO: pick a user selected study or max
AND s.start_year = ( SELECT MIN(s.start_year) FROM bla_studies s WHERE s.user_id = $1 )
ORDER BY dp.bloc_type, dp.seq
;
-- WITH user_blueprint_courses AS (
-- 	SELECT course_code, ARRAY_AGG(y.academic_year) AS academic_years FROM blueprint_years y
-- 	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = y.id
-- 	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id = bs.id
-- 	WHERE y.user_id = $1
-- 	GROUP BY course_code
-- ),
-- user_blueprint_semesters AS (
-- 	SELECT
-- 		t.course_code,
-- 		array_agg(bc.course_code IS NOT NULL) AS semesters
-- 	FROM unnest($2::text[]) t(course_code)
-- 	LEFT JOIN blueprint_years by ON by.user_id=$1
-- 	LEFT JOIN blueprint_semesters bs ON by.id = bs.blueprint_year_id
-- 	LEFT JOIN blueprint_courses bc ON bs.id = bc.blueprint_semester_id AND bc.course_code=t.course_code
-- 	GROUP BY t.course_code
-- )
-- SELECT
-- 	dp.bloc_subject_code,
-- 	COALESCE(dp.bloc_limit, -1) bloc_limit,
-- 	COALESCE(dp.bloc_name, '') bloc_name,
-- 	COALESCE(dp.bloc_note, '') bloc_note,
-- 	COALESCE(dp.note, '') note,
-- 	c.code,
-- 	c.title,
-- 	c.credits,
-- 	c.start_semester,
-- 	COALESCE(c.lecture_range1, -1) lecture_range1,
-- 	COALESCE(c.lecture_range2, -1) lecture_range2,
-- 	COALESCE(c.seminar_range1, -1) seminar_range1,
-- 	COALESCE(c.seminar_range2, -1) seminar_range2,
-- 	c.exam_type,
-- 	c.guarantors,
-- 	ubc.academic_years,
-- 	CASE WHEN dp.bloc_type = 'A' THEN TRUE WHEN dp.bloc_type = 'B' THEN FALSE END AS is_compulsory
-- FROM bla_studies bs
-- LEFT JOIN degree_plans dp ON bs.degree_plan_code = dp.plan_code AND bs.start_year = dp.plan_year
-- LEFT JOIN courses c ON dp.course_code = c.code
-- LEFT JOIN user_blueprint_courses ubc ON ubc.course_code = c.code
-- WHERE bs.user_id = $1
-- AND dp.lang = $2
-- AND c.lang = $2
-- AND interchangeability IS NULL
-- -- TODO: pick a user selected study or max
-- AND bs.start_year = ( SELECT MIN(bs.start_year) FROM bla_studies bs WHERE bs.user_id = $1 )
-- ORDER BY dp.bloc_type, dp.seq
-- ;
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
