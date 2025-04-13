package sqlquery

const RateCategory = `
INSERT INTO course_ratings (user_id, course_code, category_code, rating)
VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (user_id, category_code, course_code) DO
UPDATE SET rating=$4
;
`

const DeleteCategoryRating = `--sql
DELETE FROM course_ratings
WHERE user_id=$1
AND course_code=$2
AND category_code=$3
;
`

// TODO: valid_from use the date he finished the course???
const Rate = `--sql
INSERT INTO course_overall_ratings (user_id, course_code, rating)
VALUES (
    $1, $2, $3
)
ON CONFLICT (user_id, course_code) DO
UPDATE SET rating=$3
;
`
const CourseOverallRating = `--sql
WITH avg_course_overall_ratings AS (
    SELECT cor.course_code, AVG(rating) AS avg_rating, COUNT(rating) AS rating_count FROM course_overall_ratings cor
    WHERE course_code=$2
    GROUP BY course_code
)
SELECT rating, avg_rating, rating_count FROM course_overall_ratings cor
LEFT JOIN avg_course_overall_ratings acor ON cor.course_code=acor.course_code
WHERE cor.user_id=$1
AND cor.course_code=$2
`

const DeleteRating = `--sql
DELETE FROM course_overall_ratings
WHERE user_id=$1
AND course_code=$2
;
`

const Course = `--sql
WITH user_course_overall_ratings AS (
    SELECT cor.course_code, cor.rating AS rating FROM course_overall_ratings cor
    WHERE cor.user_id=$1
    AND cor.course_code=$2
),
avg_course_overall_ratings AS (
    SELECT cor.course_code, AVG(rating) AS avg_rating, COUNT(rating) AS rating_count FROM course_overall_ratings cor
    WHERE course_code=$2
    GROUP BY course_code
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
    ucor.rating,
    avg_cor.avg_rating,
    avg_cor.rating_count,
    c.preqrequisities,
    c.corequisities,
    c.incompatibilities,
    c.interchangebilities,
    c.classes,
    c.classifications
FROM courses c
LEFT JOIN user_course_overall_ratings ucor ON c.code = ucor.course_code
LEFT JOIN avg_course_overall_ratings avg_cor ON c.code = avg_cor.course_code
WHERE code = $2 AND c.lang = $3;
`

const Rating = `--sql
WITH user_ratings AS (
    SELECT cr.category_code, cr.course_code, cr.rating FROM course_ratings cr
    WHERE cr.user_id=$1
    AND cr.course_code=$2
),
avg_course_rating AS (
    SELECT cr.category_code, cr.course_code, AVG(cr.rating) AS avg_rating, COUNT(cr.rating) AS rating_count FROM course_ratings cr
    WHERE cr.course_code=$2
    GROUP BY cr.course_code, cr.category_code
)
SELECT
    crc.code AS category_code,
    crc.title AS rating_title,
    ur.rating,
    avg_cr.avg_rating,
    avg_cr.rating_count
FROM course_rating_categories crc
LEFT JOIN user_ratings ur ON ur.category_code = crc.code
LEFT JOIN avg_course_rating avg_cr ON avg_cr.category_code = crc.code
WHERE crc.lang = $3
`
