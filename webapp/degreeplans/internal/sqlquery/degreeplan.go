package sqlquery

const DegreePlanMetadataForSearch = `--sql
SELECT
	plan_code,
	title,
	study_type,
	valid_from,
	valid_to
FROM degree_plans
WHERE plan_code = ANY($1)
	AND lang = $2
ORDER BY study_type ASC, valid_to DESC, title ASC;
`
