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
SELECT
    code,
    title,
    faculty,
    guarantor,
    taught,
    semester_description,
    semester_count,
    taught_lang,
    lecture_range1,
    seminar_range1,
    COALESCE(lecture_range2, -1) AS lecture_range2,
    COALESCE(seminar_range2, -1) AS seminar_range2,
    exam_type,
    credits,
    guarantors,
    teachers,
    COALESCE(min_number, -1) AS min_number,
    COALESCE(capacity, '-1') AS capacity,
    COALESCE(annotation, '{}') AS annotation,
    COALESCE(aim, '{}') AS aim,
    COALESCE(requirements, '{}') AS requirements,
    COALESCE(requirements, '{}') AS syllabus
FROM bla_courses c
LEFT JOIN start_semester_to_desc sd ON c.start_semester = sd.id AND c.lang = sd.lang
WHERE code = $1 AND c.lang = $2;
`
