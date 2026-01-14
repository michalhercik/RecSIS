package sqlquery

const SelectCategoriesAndValues = `--sql
	SELECT
		fc.id AS category_id,
		fc.facet_id AS category_facet_id,
		fc.title_cs AS category_title_cs,
		fc.description_cs AS category_description_cs,
		fc.title_en AS category_title_en,
		fc.description_en AS category_description_en,
		fc.condition AS category_condition,
		fc.displayed_value_limit AS category_displayed_value_limit,
		fv.id AS value_id,
		fv.facet_id AS value_facet_id,
		fv.title_cs AS value_title_cs,
		fv.description_cs AS value_description_cs,
		fv.title_en AS value_title_en,
		fv.description_en AS value_description_en
	FROM filter_categories fc
	LEFT JOIN filter_values fv ON fc.id = fv.category_id
	WHERE fc.filter_id = $1
	ORDER BY fc.position, fv.position
`
