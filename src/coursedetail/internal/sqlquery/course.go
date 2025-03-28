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
),
avg_course_overall_ratings AS (
    SELECT cor.course_code, AVG(rating) AS avg_overall_rating FROM course_overall_ratings cor
    WHERE course_code=$2
    GROUP BY course_code
),
avg_course_rating AS (
    SELECT cr.category_code, cr.course_code, AVG(cr.rating) AS avg_rating, crc.title, crc.lang FROM course_ratings cr
    LEFT JOIN course_rating_categories crc ON cr.category_code = crc.code
    WHERE course_code=$2
    GROUP BY cr.course_code, cr.category_code, crc.title, crc.lang
)
SELECT
    c.code,
    c.title,
    c.faculty,
    c.guarantor,
    c.taught,
    c.start_semester,
    c.semester_count,
    c.taught_lang,
    COALESCE(c.lecture_range1, -1) AS lecture_range1,
    COALESCE(c.seminar_range1, -1) AS seminar_range1,
    COALESCE(c.lecture_range2, -1) AS lecture_range2,
    COALESCE(c.seminar_range2, -1) AS seminar_range2,
    c.exam_type,
    c.credits,
    c.guarantors,
    c.teachers,
    COALESCE(c.min_number, -1) AS min_number,
    COALESCE(c.capacity, '-1') AS capacity,
    c.annotation,
    c.aim,
    c.requirements_for_assesment,
    c.syllabus,
    c.literature,
    c.entry_requirements,
    c.terms_of_passing,
    cor.rating AS overall_rating,
    ur.category_code,
    ur.title AS rating_title,
    ur.rating,
    avg_cr.avg_rating,
    avg_cor.avg_overall_rating,
    c.comments,
    c.preqrequisities,
    c.corequisities,
    c.incompatibilities,
    c.interchangebilities,
    c.classes,
    c.classifications
FROM courses c
LEFT JOIN course_overall_ratings cor ON c.code = cor.course_code AND cor.user_id = (SELECT user_id FROM user_session)
LEFT JOIN user_ratings ur ON c.code = ur.course_code AND ur.lang = c.lang
LEFT JOIN avg_course_overall_ratings avg_cor ON c.code = avg_cor.course_code
LEFT JOIN avg_course_rating avg_cr ON c.code = avg_cr.course_code AND avg_cr.lang = c.lang AND avg_cr.category_code = ur.category_code
WHERE code = $2 AND c.lang = $3;
`
