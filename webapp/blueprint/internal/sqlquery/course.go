package sqlquery

/*
TODO: OUT of date -> upgrade to named param
Param order:
	1. student
	2. blueprint_years.position
	3. blueprint_semesters.semester
	4. blueprint_semesters.position
	5. blueprint.id
	6. course.code
*/

const SelectCourses = `--sql
SELECT
	y.academic_year,
	bs.semester,
	bs.folded,
	bc.id,
	c.code,
	c.title,
	c.start_semester,
	c.lecture_range_winter,
	c.lecture_range_summer,
	c.seminar_range_winter,
	c.seminar_range_summer,
	c.exam,
	c.credits,
	c.guarantors
FROM blueprint_years y
INNER JOIN blueprint_semesters bs
	ON y.id=bs.blueprint_year_id
INNER JOIN blueprint_courses bc
	ON bs.id=bc.blueprint_semester_id
LEFT JOIN courses c
	ON bc.course_code=c.code
WHERE y.user_id = $1
	AND (c.lang = $2 OR c.lang IS NULL)
ORDER BY y.academic_year, bs.semester, bc.position;
`

const SelectRequisites = `--sql
SELECT
	parent_course,
	child_course,
	req_type,
	group_type
FROM requisites
WHERE target_course = $1;
`

const SelectSemestersInfo = `--sql
SELECT
	y.academic_year,
	bs.semester,
	bs.folded
FROM blueprint_years y
INNER JOIN blueprint_semesters bs
	ON y.id=bs.blueprint_year_id
WHERE y.user_id = $1
ORDER BY y.academic_year, bs.semester;
`

const MoveCourses = `--sql
WITH user_semesters AS (
	SELECT bs.id, bs.semester, y.academic_year
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON bs.blueprint_year_id = y.id
	WHERE y.user_id = $1
),
origin_courses AS (
	SELECT bc.id
	FROM user_semesters us
	LEFT JOIN blueprint_courses bc
		ON bc.blueprint_semester_id = us.id
	WHERE bc.id = ANY($2)
),
target_semester_id AS (
	SELECT us.id
	FROM user_semesters us
	WHERE us.academic_year = $3
	AND us.semester = $4
),
condition AS (
	SELECT count(*) as c
	FROM target_semester_id ts
	LEFT JOIN blueprint_courses bc
		ON bc.blueprint_semester_id = ts.id
	INNER JOIN origin_courses oc
		ON bc.id = oc.id
	WHERE bc.position < $5
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = ts.id,
	position = $5,
	secondary_position =
		CASE WHEN (SELECT c FROM condition) > 0
			THEN 2 + array_position($2, bc.id)
			ELSE 1 - array_length($2, 1) + array_position($2, bc.id)
		END
FROM target_semester_id ts
WHERE bc.id IN ( SELECT id FROM origin_courses );
`

const AppendCourses = `--sql
WITH origin AS (
	SELECT bc.id
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
	WHERE y.user_id = $1
	AND bc.id = any($4)
),
target_semester_position AS (
	SELECT bs.id AS blueprint_semester_id, COALESCE(bc.position, 0) AS max_position
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
	WHERE y.user_id = $1
		AND y.academic_year = $2
		AND bs.semester = $3
	ORDER BY bc.position DESC
	LIMIT 1
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = t.blueprint_semester_id,
	position = t.max_position + array_position($4, id)
FROM target_semester_position t
WHERE bc.id IN ( SELECT id FROM origin );
`

const UnassignCoursesBySemester = `--sql
WITH origin_semester_id AS (
	SELECT bs.id
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	WHERE y.user_id = $1
	AND y.academic_year = $2
	AND bs.semester = $3
),
unassigned AS (
	SELECT bs.id, COALESCE(bc.position, 0) AS max_position
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
	WHERE y.user_id = $1
		AND y.academic_year = 0
		AND bs.semester = 0
	ORDER BY bc.position DESC
	LIMIT 1
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = u.id,
	position = u.max_position + bc.position
FROM unassigned u, origin_semester_id os
WHERE bc.blueprint_semester_id = os.id;
`

const RemoveCoursesByID = `--sql
WITH target_semesters_id AS (
	SELECT bc.id
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc
		ON bs.id = bc.blueprint_semester_id
	WHERE y.user_id = $1
	AND bc.id = any($2)
)
DELETE FROM blueprint_courses bc
USING target_semesters_id tsi
WHERE bc.id = tsi.id;
`

const RemoveCoursesBySemester = `--sql
WITH target_semester_id AS (
	SELECT bs.id
	FROM blueprint_years y
	LEFT JOIN blueprint_semesters bs
		ON y.id = bs.blueprint_year_id
	WHERE y.user_id = $1
	AND y.academic_year = $2
	AND bs.semester = $3
)
DELETE FROM blueprint_courses bc
USING target_semester_id ts
WHERE bc.blueprint_semester_id = ts.id;
`
