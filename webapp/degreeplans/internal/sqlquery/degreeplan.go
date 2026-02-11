package sqlquery

const DegreePlanMetadataForSearch = `--sql
SELECT
	plan_code as code,
	title,
	valid_from,
	valid_to,
	study_type
FROM degree_plans
WHERE plan_code = ANY($1)
	AND lang = $2
ORDER BY study_type ASC, valid_to DESC, title ASC;
`
