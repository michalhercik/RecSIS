package sqlquery

const RateCategory = `
INSERT INTO course_ratings (user_id, category_code, course_code, rating)
VALUES (
    (SELECT user_id FROM sessions WHERE id=$1),
    $2, $3, $4
)
ON CONFLICT (user_id, category_code, course_code) DO
UPDATE SET rating=$4
;
`

const DeleteCategoryRating = `--sql
DELETE FROM course_ratings
WHERE user_id=(SELECT user_id FROM sessions WHERE id=$1)
AND course_code=$2
AND category_code=$3
;
`

// TODO: valid_from use the date he finished the course???
const Rate = `
INSERT INTO course_ratings (user_id, course_code, overall_rating)
VALUES (
    (SELECT user_id FROM sessions WHERE id=$1),
    $2, $3
)
ON CONFLICT (user_id, course_code) DO
UPDATE SET overall_rating=$3
;
`

const DeleteRating = `--sql
DELETE FROM course_ratings
WHERE user_id=(SELECT user_id FROM sessions WHERE id=$1)
AND course_code=$2
;
`

const Course = `
WITH user_session AS (
	SELECT DISTINCT user_id FROM sessions WHERE id=$1
),
user_ratings AS (
    SELECT cr.category_code, cr.course_code, cr.rating, crc.title, crc.lang FROM user_session
    LEFT jOIN course_ratings cr ON user_session.user_id = cr.user_id
    LEFT JOIN course_rating_categories crc ON cr.category_code = crc.code
)
SELECT
    c.code,
    c.title,
    c.faculty,
    c.guarantor,
    c.taught,
    sd.semester_description,
    c.semester_count,
    c.taught_lang,
    c.lecture_range1,
    c.seminar_range1,
    COALESCE(c.lecture_range2, -1) AS lecture_range2,
    COALESCE(c.seminar_range2, -1) AS seminar_range2,
    c.exam_type,
    c.credits,
    c.guarantors,
    c.teachers,
    COALESCE(c.min_number, -1) AS min_number,
    COALESCE(c.capacity, '-1') AS capacity,
    COALESCE(c.annotation, '{}') AS annotation,
    COALESCE(c.aim, '{}') AS aim,
    COALESCE(c.requirements, '{}') AS requirements,
    COALESCE(c.syllabus, '{}') AS syllabus,
    cor.rating AS overall_rating,
    ur.category_code,
    ur.title AS rating_title,
    ur.rating
FROM courses c
LEFT JOIN start_semester_to_desc sd ON c.start_semester = sd.id AND c.lang = sd.lang
LEFT JOIN course_overall_ratings cor ON c.code = cor.course_code AND cor.user_id = (SELECT user_id FROM user_session)
LEFT JOIN user_ratings ur ON c.code = ur.course_code AND ur.lang = c.lang
WHERE code = $2 AND c.lang = $3;
`
