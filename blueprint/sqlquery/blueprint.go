package blueprint

const blueprint = `
SELECT
	y.position 
	s.position 
	c.code,
	c.name_cs,
	c.name_en,
	c.start_semester,
	c.lecture_range1,
	c.lecture_range2,
	c.seminar_range1,
	c.seminar_range2,
	c.exam_type,
	c.credits,
	t1.sis_id, 
	t1.first_name, 
	t1.last_name, 
	t1.title_before, 
	t1.title_after,
	t2.sis_id, 
	t2.first_name, 
	t2.last_name, 
	t2.title_before, 
	t2.title_after
FROM blueprint_years AS y 
LEFT JOIN blueprint_semesters AS s on y.id=s.blueprint_year
LEFT JOIN courses AS c ON s.course=c.id
LEFT JOIN teachers AS t1 ON t1.id = c.teacher1
LEFT JOIN teachers AS t2 ON t2.id = c.teacher2
WHERE student = $1;
`
