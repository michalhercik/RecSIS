package sqlquery

// TODO: valid_from use the date he finished the course???
const OverallRating = `
INSERT INTO course_ratings (user_id, course_code, overall_rating)
VALUES (
    (SELECT user_id FROM sessions WHERE id=$1),
    $2, $3
)
ON CONFLICT (user_id, course_code) DO
UPDATE SET overall_rating=$3
;
`

const Course = `
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
)
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
    COALESCE(requirements, '{}') AS syllabus,
    cr.overall_rating ,
    cr.difficulty_rating,
    cr.workload_rating,
    cr.usefulness_rating,
    cr.fun_rating
FROM courses c
LEFT JOIN start_semester_to_desc sd ON c.start_semester = sd.id AND c.lang = sd.lang
LEFT JOIN course_ratings cr ON c.code = cr.course_code AND cr.user_id = (SELECT user_id FROM user_session)
WHERE code = $2 AND c.lang = $3;
`
