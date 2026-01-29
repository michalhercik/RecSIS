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
	dp.title AS degree_plan_title,
	dp.field_code,
	dp.field_title,
	dp.valid_from AS degree_plan_valid_from,
	dp.valid_to AS degree_plan_valid_to,
	dp.requisite_graph_data,
	dpc.bloc_subject_code,
	COALESCE(dpc.bloc_limit, -1) bloc_limit,
	dpc.bloc_name,
	dpc.bloc_type,
	c.code,
	c.title,
	COALESCE(c.credits, 0) AS credits, -- there exists N# courses without credits in DPs, TODO: what to do with them? 
	COALESCE(c.start_semester, '3') AS start_semester, -- there exist N# courses without start_semester in DPs, TODO: what to do with them?
	c.lecture_range_winter,
	c.lecture_range_summer,
	c.seminar_range_winter,
	c.seminar_range_summer,
	COALESCE(c.exam, '') AS exam, -- there exist N# courses without exam type in DPs, TODO: what to do with them?
	c.guarantors,
	dpc.recommended_year_from,
	dpc.recommended_year_to,
	dpc.recommended_semester,
	c.credits IS NOT NULL as course_is_supported,
	ubs.semesters
FROM studies s
LEFT JOIN degree_plans dp
	ON s.degree_plan_code = dp.plan_code
	AND dp.lang = $2
LEFT JOIN degree_plan_courses dpc
	ON s.degree_plan_code = dpc.plan_code
	AND dpc.lang = $2
LEFT JOIN courses c
	ON dpc.course_code = c.code
	AND c.lang = $2
LEFT JOIN user_blueprint_semesters ubs
	ON dpc.course_code = ubs.course_code
WHERE s.user_id = $1
	AND interchangeability IS NULL
ORDER BY dpc.bloc_type, dpc.seq, c.code;
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
	dp.plan_code AS degree_plan_code,
	dp.title AS degree_plan_title,
	dp.field_code,
	dp.field_title,
	dp.valid_from AS degree_plan_valid_from,
	dp.valid_to AS degree_plan_valid_to,
	dp.requisite_graph_data,
	dpc.bloc_subject_code,
	COALESCE(dpc.bloc_limit, -1) bloc_limit,
	dpc.bloc_name,
	dpc.bloc_type,
	c.code,
	c.title,
	COALESCE(c.credits, 0) AS credits, -- there exist N# courses without credits in DPs, TODO: what to do with them? 
	COALESCE(c.start_semester, '3') AS start_semester, -- there exist N# courses without start_semester in DPs, TODO: what to do with them?
	c.lecture_range_winter,
	c.lecture_range_summer,
	c.seminar_range_winter,
	c.seminar_range_summer,
	COALESCE(c.exam, '') AS exam, -- there exist N# courses without exam type in DPs, TODO: what to do with them?
	c.guarantors,
	dpc.recommended_year_from,
	dpc.recommended_year_to,
	dpc.recommended_semester,
	c.credits IS NOT NULL as course_is_supported,
	ubs.semesters
FROM degree_plans dp
LEFT JOIN degree_plan_courses dpc
	ON dp.plan_code = dpc.plan_code
	AND dpc.lang = $3
LEFT JOIN courses c
	ON dpc.course_code = c.code
	AND c.lang = $3
LEFT JOIN user_blueprint_semesters ubs
	ON dpc.course_code = ubs.course_code
LEFT JOIN studies s
	ON s.user_id = $1
WHERE dp.plan_code = $2
	AND dp.lang = $3
	AND interchangeability IS NULL
ORDER BY dpc.bloc_type, dpc.seq, c.code;
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

/* === Blueprint operations === */

const ClearBlueprintCourses = `--sql
DELETE FROM blueprint_courses
USING blueprint_semesters bs, blueprint_years by
WHERE blueprint_courses.blueprint_semester_id = bs.id
AND bs.blueprint_year_id = by.id
AND by.user_id = $1
`

const CountBlueprintYears = `--sql
SELECT COALESCE(MAX(academic_year), 0) FROM blueprint_years
WHERE user_id = $1
`

const InsertMissingBlueprintYears = `--sql
INSERT INTO blueprint_years (user_id, academic_year)
SELECT $1 AS user_id, generate_series($2::int, $3::int) AS academic_year
`

const InsertMissingBlueprintSemesters = `--sql
INSERT INTO blueprint_semesters (blueprint_year_id, semester)
SELECT by.id, s.semester
FROM blueprint_years by
CROSS JOIN (VALUES (1), (2)) AS s(semester)
WHERE by.user_id = $1
AND by.academic_year BETWEEN $2::int AND $3::int
`

const InsertRecommendedPlanCourses = `--sql
INSERT INTO blueprint_courses (blueprint_semester_id, course_code, course_valid_from, position)
WITH recommended_courses AS (
    SELECT DISTINCT ON (dpc.course_code)
        dpc.course_code,
        dpc.recommended_year_from,
        dpc.recommended_year_to,
        CASE
            WHEN c.start_semester = '3' THEN COALESCE(dpc.recommended_semester, 1)
            WHEN c.start_semester = '1' THEN 1
            WHEN c.start_semester = '2' THEN 2
            ELSE COALESCE(dpc.recommended_semester, 1)
        END AS semester,
        c.valid_from
    FROM degree_plan_courses dpc
    LEFT JOIN courses c ON dpc.course_code = c.code AND c.lang = $3
    WHERE dpc.plan_code = $2 
        AND dpc.lang = $3
        AND dpc.recommended_year_from IS NOT NULL
		AND dpc.interchangeability IS NULL
	ORDER BY dpc.course_code, (dpc.recommended_semester IS NULL) ASC
),
years_and_semesters AS (
    SELECT 
        rc.course_code,
        y,
        rc.semester,
        bs.id AS blueprint_semester_id,
        rc.valid_from
    FROM recommended_courses rc
    CROSS JOIN LATERAL generate_series(rc.recommended_year_from, rc.recommended_year_to) AS y
    LEFT JOIN blueprint_years by ON by.user_id = $1 AND by.academic_year = y
    LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = by.id AND bs.semester = rc.semester
)
SELECT 
	blueprint_semester_id,
	course_code,
	valid_from,
	ROW_NUMBER() OVER (PARTITION BY blueprint_semester_id ORDER BY course_code) AS position
FROM years_and_semesters
`

const MergeRecommendedPlanCourses = `--sql
INSERT INTO blueprint_courses (blueprint_semester_id, course_code, course_valid_from, position)
WITH recommended_courses AS (
    SELECT DISTINCT ON (dpc.course_code)
        dpc.course_code,
        dpc.recommended_year_from,
        dpc.recommended_year_to,
        CASE
            WHEN c.start_semester = '3' THEN COALESCE(dpc.recommended_semester, 1)
            WHEN c.start_semester = '1' THEN 1
            WHEN c.start_semester = '2' THEN 2
            ELSE COALESCE(dpc.recommended_semester, 1)
        END AS semester,
        c.valid_from
    FROM degree_plan_courses dpc
    LEFT JOIN courses c ON dpc.course_code = c.code AND c.lang = $3
    WHERE dpc.plan_code = $2 
        AND dpc.lang = $3
        AND dpc.recommended_year_from IS NOT NULL
        AND dpc.interchangeability IS NULL
    ORDER BY dpc.course_code, (dpc.recommended_semester IS NULL) ASC
),
years_and_semesters AS (
    SELECT 
        rc.course_code,
        y,
        rc.semester,
        bs.id AS blueprint_semester_id,
        rc.valid_from
    FROM recommended_courses rc
    CROSS JOIN LATERAL generate_series(rc.recommended_year_from, rc.recommended_year_to) AS y
    LEFT JOIN blueprint_years by ON by.user_id = $1 AND by.academic_year = y
    LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id = by.id AND bs.semester = rc.semester
),
max_positions AS (
    SELECT 
        blueprint_semester_id,
        COALESCE(MAX(position), 0) AS max_pos
    FROM blueprint_courses
    GROUP BY blueprint_semester_id
),
with_positions AS (
    SELECT 
        yas.blueprint_semester_id,
        yas.course_code,
        yas.valid_from,
        COALESCE(mp.max_pos, 0) + ROW_NUMBER() OVER (PARTITION BY yas.blueprint_semester_id ORDER BY yas.course_code) AS position
    FROM years_and_semesters yas
    LEFT JOIN max_positions mp ON yas.blueprint_semester_id = mp.blueprint_semester_id
)
SELECT 
    blueprint_semester_id,
    course_code,
    valid_from,
    position
FROM with_positions
ON CONFLICT (blueprint_semester_id, course_code) DO NOTHING
`
