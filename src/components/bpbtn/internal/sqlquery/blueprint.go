package sqlquery

const InsertCourse = `--sql
	WITH target_position AS (
		SELECT
			bs.id AS blueprint_semester_id,
			COALESCE(bc.position, 0) + 1 AS last_position
		FROM blueprint_years y
		LEFT JOIN blueprint_semesters bs
			ON y.id = bs.blueprint_year_id
		LEFT JOIN blueprint_courses bc
			ON bs.id = bc.blueprint_semester_id
		WHERE y.user_id = $1
			AND y.academic_year = $2
			AND bs.semester = $3
		ORDER BY bc.position DESC
		LIMIT 1
	)
	INSERT INTO blueprint_courses(blueprint_semester_id, course_code, course_valid_from, position)
	SELECT
		tp.blueprint_semester_id,
		c.code,
		c.valid_from,
		tp.last_position + ROW_NUMBER() OVER (ORDER BY c.code)
	FROM
		target_position tp,
		UNNEST($4::text[]) AS course_code
		JOIN LATERAL (
			SELECT code, valid_from FROM courses WHERE code = course_code ORDER BY valid_from DESC LIMIT 1
		) c ON TRUE
	RETURNING id;
	`
