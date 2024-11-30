package sqlquery

const Course = `
WITH teacher_faculties AS (
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
)

SELECT
    c.id,
    c.code,
    c.name_cs,
    c.name_en,
    c.valid_from,
    c.valid_to,
    f.id,
    f.sis_id,
    f.name_cs,
    f.name_en,
    f.abbr,
    c.guarantor,
    c.taught,
    c.start_semester,
    c.semester_count,
    c.taught_lang,
    c.lecture_range1,
    c.seminar_range1,
    COALESCE(c.lecture_range2, -1),
    COALESCE(c.seminar_range2, -1),
    c.exam_type,
    c.credits,
    t1.t_id,
    t1.t_sis_id,
    t1.department,
    t1.f_id,
    t1.f_sis_id,
    t1.name_cs,
    t1.name_en,
    t1.abbr,
    t1.first_name,
    t1.last_name,
    t1.title_before,
    t1.title_after,
    t2.t_id,
    t2.t_sis_id,
    t2.department,
    t2.f_id,
    t2.f_sis_id,
    t2.name_cs,
    t2.name_en,
    t2.abbr,
    t2.first_name,
    t2.last_name,
    t2.title_before,
    t2.title_after,
    c.min_number,
    COALESCE(c.capacity, -1),
    c.annotation_cs,
    c.annotation_en,
    c.sylabus_cs,
    c.sylabus_en
FROM
    courses as c
    LEFT JOIN faculties as f on c.faculty = f.sis_id
    LEFT JOIN teacher_faculties AS t1 on c.teacher1 = t1.t_sis_id
    LEFT JOIN teacher_faculties AS t2 on c.teacher2 = t2.t_sis_id
WHERE
    code = $1
    AND valid_to > date_part('year', CURRENT_DATE)
`
