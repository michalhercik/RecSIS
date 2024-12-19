package sqlquery

const InsertCourse = `
WITH year AS (
	SELECT id FROM blueprint_years
	WHERE student=$1 AND position=$2
)

INSERT INTO blueprint_semesters (semester, course, blueprint_year, position)
VALUES(
	$3,
	(SELECT id FROM courses WHERE code=$4),
	(SELECT id FROM year),
	(
		SELECT MAX(position) + 1
		FROM blueprint_semesters
		WHERE blueprint_year=(SELECT id FROM year)
		AND semester=$3
	)
)
RETURNING id;
`

const MoveCourse = `
WITH course AS (
	SELECT * FROM blueprint_semesters
	WHERE id=$5
), new_year AS (
	SELECT * FROM blueprint_years 
	WHERE student=$3 
	AND position=$4
)

UPDATE blueprint_semesters
SET 
	semester = $1, 
	position = $2, 
	secondary_position = 2 + array_position($5, id) 
	blueprint_year = (SELECT id FROM new_year)
WHERE id=any($5);
`

const AppendCourse = `
WITH year AS (
	SELECT id FROM blueprint_years
	WHERE student=$1 AND position=$2
)

UPDATE blueprint_semesters
SET 
	semester=$3, 
	blueprint_year=(SELECT id FROM year),
	position= (
		COALESCE(
			(SELECT MAX(position) + 1
			FROM blueprint_semesters
			WHERE blueprint_year=(SELECT id FROM year)
			AND semester=$3) + array_position($4, id), 1
		)
	) 
WHERE id = any($4);
`

const UnassignYear = `
WITH unassigned_year AS (
	SELECT id FROM blueprint_years
	WHERE student=$1 AND position=0
), max_position AS (
	SELECT MAX(position) AS pos FROM blueprint_semesters
	WHERE blueprint_year=(SELECT id from unassigned_year)
), selected_year AS (
	SELECT id FROM blueprint_years
	WHERE student=$1 AND position=$2
)

UPDATE blueprint_semesters
SET 
	blueprint_year = (SELECT id from unassigned_year),
	position = (select pos from max_position) + position,
	semester = 0
WHERE blueprint_year = (SELECT id from selected_year);
`

const UnassignSemester = `
WITH unassigned_year AS (
	SELECT id FROM blueprint_years
	WHERE student=$1 AND position=0
), max_position AS (
	SELECT COALESCE(MAX(position), 0) AS pos FROM blueprint_semesters
	WHERE blueprint_year=(SELECT id from unassigned_year)
), selected_year AS (
	SELECT id FROM blueprint_years
	WHERE student=$1 AND position=$2
)

UPDATE blueprint_semesters
SET 
	blueprint_year = (SELECT id from unassigned_year),
	position = (select pos from max_position) + position,
	semester = 0
WHERE blueprint_year = (SELECT id from selected_year) 
AND semester=$3;
`

const DeleteCourse = `
DELETE FROM blueprint_semesters 
WHERE blueprint_year IN (SELECT id FROM blueprint_years WHERE student=$1)
AND id = any($2);
`

const DeleteCoursesBySemester = `
DELETE FROM blueprint_semesters
WHERE semester=$3
AND blueprint_year IN (
	SELECT id FROM blueprint_years 
		WHERE student=$1
		AND position=$2
);
`
const DeleteCoursesByYear = `
DELETE FROM blueprint_semesters
WHERE blueprint_year IN (
	SELECT id FROM blueprint_years 
		WHERE student=$1
		AND position=$2
);
`

const SelectCourses = `
SELECT
	y.position,
	s.id,
	s.position,
	s.semester,
	c.code,
	c.name_cs,
	c.name_en,
	COALESCE(c.start_semester, -1),
	COALESCE(c.semester_count, -1),
	COALESCE(c.lecture_range1, -1),
	COALESCE(c.lecture_range2, -1),
	COALESCE(c.seminar_range1, -1),
	COALESCE(c.seminar_range2, -1),
	COALESCE(c.exam_type, ''),
	COALESCE(c.credits, -1),
	COALESCE(t1.sis_id, -1), 
	COALESCE(t1.first_name, ''),
	COALESCE(t1.last_name, ''),
	COALESCE(t1.title_before, ''),
	COALESCE(t1.title_after,''),
	COALESCE(t2.sis_id, -1), 
	COALESCE(t2.first_name, ''),
	COALESCE(t2.last_name, ''),
	COALESCE(t2.title_before, ''),
	COALESCE(t2.title_after, ''),
	COALESCE(t3.sis_id, -1), 
	COALESCE(t3.first_name, ''),
	COALESCE(t3.last_name, ''),
	COALESCE(t3.title_before, ''),
	COALESCE(t3.title_after, '')
FROM blueprint_semesters AS s
LEFT JOIN blueprint_years AS y on y.id=s.blueprint_year
LEFT JOIN courses AS c ON s.course=c.id
LEFT JOIN teachers AS t1 ON t1.sis_id = c.teacher1
LEFT JOIN teachers AS t2 ON t2.sis_id = c.teacher2
LEFT JOIN teachers AS t3 ON t3.sis_id = c.teacher3
WHERE student = $1
ORDER BY s.position;
`
