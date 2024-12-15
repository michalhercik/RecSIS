package sqlquery

const CountCourses = `
	SELECT COUNT(*) FROM courses
`

const GetCourses = `
SELECT
	c.code,
	COALESCE(c.name_cs, ''),
	COALESCE(c.name_en, ''),
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
FROM courses AS c
LEFT JOIN teachers AS t1 ON t1.sis_id = c.teacher1
LEFT JOIN teachers AS t2 ON t2.sis_id = c.teacher2
LEFT JOIN teachers AS t3 ON t3.sis_id = c.teacher3
OFFSET $1 LIMIT $2
`
