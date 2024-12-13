package sqlquery

const ShiftCourses = `
UPDATE blueprint_semesters
SET position = position + 1
WHERE blueprint_year = (
	SELECT id FROM blueprint_years 
	WHERE student=$1 AND position=$2
	)
AND semester = $3
AND position >= $4;
`

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
);
`

const MoveCourse = `
UPDATE blueprint_semesters
SET semester=$1, position=$2, blueprint_year=(
	SELECT id FROM blueprint_years 
	WHERE student=$3 
	AND position=$4
	)
WHERE id=$5;
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
		SELECT MAX(position) + 1
		FROM blueprint_semesters
		WHERE blueprint_year=(SELECT id FROM year)
		AND semester=$3
	) 
WHERE id=$4;
`

const DeleteCourse = `
DELETE FROM blueprint_semesters 
WHERE blueprint_year IN (SELECT id FROM blueprint_years WHERE student=$1)
AND id = $2
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
	c.start_semester,
	c.semester_count,
	c.lecture_range1,
	COALESCE(c.lecture_range2, -1),
	c.seminar_range1,
	COALESCE(c.seminar_range2, -1),
	c.exam_type,
	c.credits,
	COALESCE(t1.sis_id, -1), 
	COALESCE(t1.first_name, ''),
	COALESCE(t1.last_name, ''),
	COALESCE(t1.title_before, ''),
	COALESCE(t1.title_after,''),
	COALESCE(t2.sis_id, -1), 
	COALESCE(t2.first_name, ''),
	COALESCE(t2.last_name, ''),
	COALESCE(t2.title_before, ''),
	COALESCE(t2.title_after, '')
FROM blueprint_semesters AS s
LEFT JOIN blueprint_years AS y on y.id=s.blueprint_year
LEFT JOIN courses AS c ON s.course=c.id
LEFT JOIN teachers AS t1 ON t1.sis_id = c.teacher1
LEFT JOIN teachers AS t2 ON t2.sis_id = c.teacher2
WHERE student = $1;
`
