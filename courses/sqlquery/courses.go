package sqlquery

const CountCourses = `
	SELECT COUNT(*) FROM courses
`

const GetCourses = `
SELECT
	c.code,
	c.name_cs,
	c.name_en,
	c.start_semester,
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
FROM courses AS c
LEFT JOIN teachers AS t1 ON t1.id = c.teacher1
LEFT JOIN teachers AS t2 ON t2.id = c.teacher2
OFFSET $1 LIMIT $2
`
