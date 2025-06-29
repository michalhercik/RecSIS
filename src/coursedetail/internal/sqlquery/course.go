package sqlquery

const Course = `--sql
WITH user_course_overall_ratings AS (
	SELECT cor.course_code, cor.rating AS rating
	FROM course_overall_ratings cor
	WHERE cor.user_id = $1
	AND cor.course_code = $2
),
avg_course_overall_ratings AS (
	SELECT cor.course_code, AVG(rating) AS avg_rating, COUNT(rating) AS rating_count
	FROM course_overall_ratings cor
	WHERE course_code = $2
	GROUP BY course_code
),
degree_plan AS (
	SELECT dp.course_code
	FROM bla_studies bs
	LEFT JOIN degree_plans dp
		ON dp.plan_code = bs.degree_plan_code
		AND dp.plan_year = bs.start_year
	WHERE bs.user_id = $1
		AND dp.course_code = $2
		AND dp.lang = 'cs'
),
user_blueprint_semesters AS (
	SELECT array_agg(course_code IS NOT NULL) AS semesters
	FROM (
		SELECT
			by.user_id,
			bc.course_code
		FROM blueprint_years by
		LEFT JOIN blueprint_semesters bs
			ON by.id = bs.blueprint_year_id
		LEFT JOIN blueprint_courses bc
			ON bs.id = bc.blueprint_semester_id
			AND bc.course_code = $2
		WHERE by.user_id = $1
		ORDER BY by.academic_year, bs.semester
	)
)
SELECT
	c.code,
	c.title,
	c.faculty,
	c.guarantor,
	c.taught,
	c.start_semester,
	c.taught_lang,
	c.lecture_range1,
	c.seminar_range1,
	c.lecture_range2,
	c.seminar_range2,
	c.range_unit,
	c.exam_type,
	c.credits,
	c.guarantors,
	c.teachers,
	c.min_number,
	c.capacity,
	c.annotation,
	c.aim,
	c.requirements_for_assesment,
	c.syllabus,
	c.literature,
	c.entry_requirements,
	c.terms_of_passing,
	ucor.rating,
	acor.avg_rating,
	acor.rating_count,
	c.preqrequisities,
	c.corequisities,
	c.incompatibilities,
	c.interchangebilities,
	c.classes,
	c.classifications,
	ubs.semesters,
	dp.course_code IS NOT NULL AS in_degree_plan
FROM courses c
LEFT JOIN user_course_overall_ratings ucor
	ON c.code = ucor.course_code
LEFT JOIN avg_course_overall_ratings acor
	ON c.code = acor.course_code
LEFT JOIN degree_plan dp
	ON c.code = dp.course_code
LEFT JOIN user_blueprint_semesters ubs
	ON TRUE
WHERE code = $2
	AND c.lang = $3;
`

const Rating = `--sql
WITH user_ratings AS (
	SELECT cr.category_code, cr.course_code, cr.rating
	FROM course_ratings cr
	WHERE cr.user_id = $1
	AND cr.course_code = $2
),
avg_course_rating AS (
	SELECT cr.category_code, cr.course_code, AVG(cr.rating) AS avg_rating, COUNT(cr.rating) AS rating_count
	FROM course_ratings cr
	WHERE cr.course_code = $2
	GROUP BY cr.course_code, cr.category_code
)
SELECT
	crc.code AS category_code,
	crc.title AS rating_title,
	ur.rating,
	avg_cr.avg_rating,
	avg_cr.rating_count
FROM course_rating_categories crc
LEFT JOIN user_ratings ur
	ON ur.category_code = crc.code
LEFT JOIN avg_course_rating avg_cr
	ON avg_cr.category_code = crc.code
WHERE crc.lang = $3;
`

const RateCategory = `
INSERT INTO course_ratings (user_id, course_code, category_code, rating)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, category_code, course_code) DO
UPDATE SET rating = $4;
`

const DeleteCategoryRating = `--sql
DELETE FROM course_ratings
WHERE user_id = $1
AND course_code = $2
AND category_code = $3;
`

// TODO: valid_from use the date he finished the course???
const Rate = `--sql
INSERT INTO course_overall_ratings (user_id, course_code, rating)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, course_code) DO
UPDATE SET rating = $3;
`

const DeleteRating = `--sql
DELETE FROM course_overall_ratings
WHERE user_id = $1
AND course_code = $2;
`

const CourseOverallRating = `--sql
SELECT
	(
		SELECT rating
		FROM course_overall_ratings
		WHERE user_id = $1 AND course_code = $2
		LIMIT 1
	) AS rating,
	(
		SELECT AVG(rating)
		FROM course_overall_ratings
		WHERE course_code = $2
	) AS avg_rating,
	(
		SELECT COUNT(rating)
		FROM course_overall_ratings
		WHERE course_code = $2
	) AS rating_count
;
`
