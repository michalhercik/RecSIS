package sqlquery

const CourseTeachers = `
SELECT 
    COALESCE(sis_id, -1),
    COALESCE(first_name, ''),
    COALESCE(last_name, ''),
    COALESCE(title_before, ''),
    COALESCE(title_after, '')
FROM course_teachers 
LEFT JOIN teachers ON teacher = sis_id
WHERE course=$1
ORDER BY last_name, first_name;
`

const Course = `
WITH 
teacher_faculties AS (
    SELECT
    t.id AS t_id,
    t.sis_id AS t_sis_id,
    t.department,
    f.id AS f_id,
    f.sis_id AS f_sis_id,
    f.name_cs,
    f.name_en,
    f.abbr,
    t.first_name,
    t.last_name,
    t.title_before,
    t.title_after
    FROM
        teachers AS t
        LEFT JOIN faculties AS f ON t.faculty = f.sis_id
), 
texts AS (
    SELECT * FROM course_texts 
    WHERE course=$1 AND lang=$2 AND audience='ALL'
)

SELECT
    c.code,
    CASE WHEN $2='CZE' THEN c.name_cs ELSE c.name_en END,
    f.sis_id,
    CASE WHEN $2='CZE' THEN f.name_cs ELSE f.name_en END,
    f.abbr,
    c.guarantor,
    c.taught,
    c.start_semester,
    c.semester_count,
    CASE WHEN c.taught_lang = 'CZE' THEN 'Czech' ELSE 'English' END,
    c.lecture_range1,
    c.seminar_range1,
    COALESCE(c.lecture_range2, -1),
    COALESCE(c.seminar_range2, -1),
    c.exam_type,
    c.credits,
    COALESCE(t1.t_sis_id, -1),
    COALESCE(t1.first_name, ''),
    COALESCE(t1.last_name, ''),
    COALESCE(t1.title_before, ''),
    COALESCE(t1.title_after, ''),
    COALESCE(t2.t_sis_id, -1),
    COALESCE(t2.first_name, ''),
    COALESCE(t2.last_name, ''),
    COALESCE(t2.title_before, ''),
    COALESCE(t2.title_after, ''),
    COALESCE(t3.t_sis_id, -1),
    COALESCE(t3.first_name, ''),
    COALESCE(t3.last_name, ''),
    COALESCE(t3.title_before, ''),
    COALESCE(t3.title_after, ''),
    COALESCE(c.min_number, -1),
    COALESCE(c.capacity, -1),
    COALESCE(annotation.title, ''),
    COALESCE(annotation.content, ''),
    COALESCE(completition.title, ''),
    COALESCE(completition.content, ''),
    COALESCE(exam.title, ''),
    COALESCE(exam.content, ''),
    COALESCE(sylabus.title, ''),
    COALESCE(sylabus.content, '')
FROM
    courses as c
    LEFT JOIN faculties as f on c.faculty = f.sis_id
    LEFT JOIN teacher_faculties AS t1 on c.teacher1 = t1.t_sis_id
    LEFT JOIN teacher_faculties AS t2 on c.teacher2 = t2.t_sis_id
    LEFT JOIN teacher_faculties AS t3 on c.teacher3 = t3.t_sis_id
    LEFT JOIN (
        SELECT * FROM texts WHERE text_type='A' 
        ) AS annotation ON c.code = annotation.course
    LEFT JOIN (
        SELECT * FROM texts WHERE text_type='C' 
        ) AS completition ON c.code = completition.course
    LEFT JOIN (
        SELECT * FROM texts WHERE text_type='P' 
        ) AS exam ON c.code = exam.course
    LEFT JOIN (
        SELECT * FROM texts WHERE text_type='S' 
        ) AS sylabus ON c.code = sylabus.course
WHERE code = $1;
`
