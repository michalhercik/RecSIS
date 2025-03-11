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

const InsertCourse = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
target_position AS (
	SELECT bs.id AS blueprint_semester_id, COALESCE(bc.position, 0) + 1 AS position FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
	WHERE y.academic_year=$2
	AND bs.semester=$3
	ORDER BY bc.position DESC
	LIMIT 1
),
target_course AS (
	SELECT code, valid_from FROM courses WHERE code=$4 ORDER BY valid_from DESC LIMIT 1
)
INSERT INTO blueprint_courses(blueprint_semester_id, course_code, course_valid_from, position)
VALUES (
	(SELECT blueprint_semester_id FROM target_position),
	(SELECT code FROM target_course),
	(SELECT valid_from FROM target_course),
	(SELECT position FROM target_position)
)
RETURNING id
;
`

const MoveCourses = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
user_semesters AS (
	SELECT bs.id, bs.semester, y.academic_year FROM sessions s
	LEFT JOIN blueprint_years y ON s.user_id = y.user_id
	LEFT JOIN blueprint_semesters bs ON bs.blueprint_year_id=y.id
),
origin_courses AS (
	SELECT bc.id FROM user_semesters us
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id=us.id
	WHERE bc.id = ANY($2)
),
target_semester_id AS (
	SELECT us.id FROM user_semesters us
	WHERE us.academic_year = $3
	AND us.semester = $4
),
condition AS (
	SELECT count(*) as c FROM target_semester_id ts
	LEFT JOIN blueprint_courses bc ON bc.blueprint_semester_id=ts.id
	INNER JOIN origin_courses oc ON bc.id=oc.id
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
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
origin AS (
	SELECT bc.id FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
	WHERE bc.id = any($4)
),
target_semester_position AS (
	SELECT bs.id AS blueprint_semester_id, COALESCE(bc.position, 0) AS max_position FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
	WHERE y.academic_year=$2
	AND bs.semester=$3
	ORDER BY bc.position DESC
	LIMIT 1
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = t.blueprint_semester_id,
	position = t.max_position + array_position($4, id)
FROM target_semester_position t
WHERE bc.id IN ( SELECT id FROM origin );
`

const UnassignYear = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
origin AS (
	SELECT bs.id, bs.semester FROM sessions s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	WHERE y.academic_year=$2
),
unassigned AS (
	SELECT bs.id, COALESCE(position, 0) AS max_position FROM sessions s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
	WHERE y.academic_year=0
	AND bs.semester = 0
	ORDER BY bc.position DESC
	LIMIT 1
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = u.id,
	position = u.max_position + (o.semester * bc.position)
FROM unassigned u, origin o
WHERE o.id = bc.blueprint_semester_id
;
`

const UnassignSemester = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
origin_semester_id AS (
	SELECT bs.id FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	WHERE y.academic_year=$2
	AND bs.semester=$3
),
unassigned AS (
	SELECT bs.id, COALESCE(bc.position, 0) AS max_position FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
	WHERE y.academic_year=0
	AND bs.semester = 0
	ORDER BY bc.position DESC
	LIMIT 1
)
UPDATE blueprint_courses bc
SET blueprint_semester_id = u.id,
	position = u.max_position + bc.position
FROM unassigned u, origin_semester_id os
WHERE bc.blueprint_semester_id = os.id
;
`

const DeleteCoursesByID = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
target_semesters_id AS (
	SELECT bc.id FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
	WHERE bc.id = any($2)
)
DELETE FROM blueprint_courses
WHERE id IN ( SELECT id FROM target_semesters_id )
;
`

const DeleteCoursesBySemester = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
target_semester_id AS (
	SELECT bs.id FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	WHERE y.academic_year=$2
	AND bs.semester=$3
)
DELETE FROM blueprint_courses bc
USING target_semester_id ts
WHERE bc.blueprint_semester_id = ts.id
;
`
const DeleteCoursesByYear = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
target_semesters_id AS (
	SELECT bs.id FROM user_session s
	LEFT JOIN blueprint_years y ON s.user_id=y.user_id
	LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
	WHERE y.academic_year=$2
)
DELETE FROM blueprint_courses
WHERE blueprint_semester_id IN
	( SELECT id FROM target_semesters_id )
;
`

const SelectCourses = `--sql
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
)
SELECT
	y.academic_year,
	bc.id,
	-- bc.position,
	bs.semester,
	c.code,
	c.title,
	COALESCE(c.start_semester, -1) start_semester,
	COALESCE(c.semester_count, -1) semester_count,
	COALESCE(c.lecture_range1, -1) lecture_range1,
	COALESCE(c.lecture_range2, -1) lecture_range2,
	COALESCE(c.seminar_range1, -1) seminar_range1,
	COALESCE(c.seminar_range2, -1) seminar_range2,
	COALESCE(c.exam_type, '') exam_type,
	COALESCE(c.credits, -1) credits,
	c.guarantors
FROM user_session s
LEFT JOIN blueprint_years y ON s.user_id=y.user_id
LEFT JOIN blueprint_semesters bs ON y.id=bs.blueprint_year_id
LEFT JOIN blueprint_courses bc ON bs.id=bc.blueprint_semester_id
LEFT JOIN courses c ON bc.course_code=c.code AND bc.course_valid_from = c.valid_from
WHERE c.lang = $2
ORDER BY y.academic_year, bs.semester, bc.position
;
`
