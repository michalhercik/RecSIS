package sqlquery

const DegreePlan = `
WITH student AS (
	SELECT degree_plan_code, start_year FROM bla_studies
	WHERE user_id = (SELECT user_id FROM sessions WHERE id=$1)
	ORDER BY start_year DESC
	LIMIT 1
)

SELECT blocs FROM degree_plans
WHERE code = (SELECT degree_plan_code FROM student)
AND plan_year = (SELECT start_year FROM student)
AND lang=$2
`
