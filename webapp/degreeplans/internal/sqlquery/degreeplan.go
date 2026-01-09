package sqlquery

const UserDegreePlanCode = `--sql
SELECT degree_plan_code
FROM studies
WHERE user_id = $1
`

const DegreePlanMetadataForSearch = `--sql
SELECT
	plan_code,
	title,
	study_type,
	valid_from,
	valid_to
FROM degree_plan_metadata
WHERE plan_code = ANY($1)
	AND lang = $2
ORDER BY study_type ASC, valid_to DESC, title ASC;
`
