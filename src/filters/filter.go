package filters

import (
	"database/sql"
	"fmt"
	"iter"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/language"
)

type FilterBuilder struct {
	categ []FilterCategory
}

func (fb FilterBuilder) IsLastCategory(categoryID string) bool {
	if len(fb.categ) == 0 {
		return true
	}
	lastID := fb.categ[len(fb.categ)-1].id
	return lastID != categoryID
}

func (fb *FilterBuilder) Category(identity FilterIdentity, displayedValueLimit int) {
	fb.categ = append(fb.categ, FilterCategory{
		FilterIdentity:      identity,
		displayedValueLimit: displayedValueLimit,
		values:              []FilterValue{},
	})
}

func (fb *FilterBuilder) Value(identity FilterIdentity) {
	value := MakeFilterValue(identity)
	category := fb.categ[len(fb.categ)-1]
	category.values = append(category.values, value)
	fb.categ[len(fb.categ)-1] = category
}

func (fb *FilterBuilder) Build() Filters {
	return MakeFilters(fb.categ)
}

func (fb *FilterBuilder) FilterCategories() []FilterCategory {
	return fb.categ
}

type Filters struct {
	DB           *sqlx.DB
	Filter       string
	categories   []FilterCategory
	facets       []string
	idToCategory map[string]FilterCategory
	idToValue    map[string]FilterValue
	prefix       string
}

func (f *Filters) Init() error {
	const query = `--sql
	SELECT
		fc.id AS category_id,
		fc.facet_id AS category_facet_id,
		fc.title_cs AS category_title_cs,
		fc.description_cs AS category_description_cs,
		fc.title_en AS category_title_en,
		fc.description_en AS category_description_en,
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
	// Retrieve
	tmpResult := []struct {
		CategoryID                  string         `db:"category_id"`
		CategoryFacetID             string         `db:"category_facet_id"`
		CategoryTitleCS             string         `db:"category_title_cs"`
		CategoryTitleEN             string         `db:"category_title_en"`
		CategoryDescCS              sql.NullString `db:"category_description_cs"`
		CategoryDescEN              sql.NullString `db:"category_description_en"`
		CategoryDisplayedValueLimit int            `db:"category_displayed_value_limit"`
		ValueID                     sql.NullString `db:"value_id"`
		ValueFacetID                sql.NullString `db:"value_facet_id"`
		ValueTitleCS                sql.NullString `db:"value_title_cs"`
		ValueTitleEN                sql.NullString `db:"value_title_en"`
		ValueDescCS                 sql.NullString `db:"value_description_cs"`
		ValueDescEN                 sql.NullString `db:"value_description_en"`
	}{}
	if err := f.DB.Select(&tmpResult, query, f.Filter); err != nil {
		return fmt.Errorf("failed to fetch filters: %w", err)
	}
	// Parse
	fb := FilterBuilder{}
	for _, row := range tmpResult {
		if fb.IsLastCategory(row.CategoryID) {
			fb.Category(MakeFilterIdentity(
				row.CategoryID,
				row.CategoryFacetID,
				language.MakeLangString(row.CategoryTitleCS, row.CategoryTitleEN),
				language.MakeLangString(row.CategoryDescCS.String, row.CategoryDescEN.String),
			), row.CategoryDisplayedValueLimit)
		}
		if row.ValueID.Valid {
			fb.Value(MakeFilterIdentity(
				row.ValueID.String,
				row.ValueFacetID.String,
				language.MakeLangString(row.ValueTitleCS.String, row.ValueTitleEN.String),
				language.MakeLangString(row.ValueDescCS.String, row.ValueDescEN.String),
			))
		}
	}
	f.categories = fb.FilterCategories()
	f.idToCategory = make(map[string]FilterCategory)
	f.idToValue = make(map[string]FilterValue)
	f.facets = make([]string, len(f.categories))

	for i, category := range f.categories {
		f.facets[i] = category.facetID
		f.idToCategory[category.id] = category
		for _, value := range category.values {
			f.idToValue[value.id] = value
		}
	}
	f.prefix = "par"
	return nil
}

func (f Filters) Facets() []string {
	return f.facets
}

func (f Filters) IterFacets(facets Facets, query url.Values, lang language.Language) iter.Seq[FacetIterator] {
	return func(yield func(FacetIterator) bool) {
		for _, c := range f.categories {
			f := facets[c.facetID]
			checked := query["par"+c.id]
			result := FacetIterator{
				title:   c.Title(lang),
				desc:    c.Desc(lang),
				filter:  c,
				facets:  f,
				lang:    lang,
				checked: checked,
			}
			if !yield(result) {
				return
			}
		}
	}
}

type FilterCategory struct {
	FilterIdentity
	displayedValueLimit int
	values              []FilterValue
}

func MakeFilterCategory(identity FilterIdentity, values []FilterValue) FilterCategory {
	return FilterCategory{
		FilterIdentity: identity,
		values:         values,
	}
}

func (fc FilterCategory) Values() []FilterValue {
	return fc.values
}

func (fc FilterCategory) DisplayedValueLimit() int {
	return fc.displayedValueLimit
}

type FilterValue struct {
	FilterIdentity
}

func MakeFilterValue(identity FilterIdentity) FilterValue {
	return FilterValue{
		FilterIdentity: identity,
	}
}

type FilterIdentity struct {
	id      string
	facetID string
	title   language.LangString
	desc    language.LangString
}

func MakeFilterIdentity(id, facetID string, title, desc language.LangString) FilterIdentity {
	return FilterIdentity{
		id:      id,
		facetID: facetID,
		title:   title,
		desc:    desc,
	}
}

func (fi FilterIdentity) Title(lang language.Language) string {
	return fi.title.String(lang)
}

func (fi FilterIdentity) Desc(lang language.Language) string {
	return fi.desc.String(lang)
}

func (fi FilterIdentity) ID() string {
	return fi.id
}

func MakeFilters(categories []FilterCategory) Filters {
	idToCategory := make(map[string]FilterCategory)
	idToValue := make(map[string]FilterValue)
	facets := make([]string, len(categories))

	for i, category := range categories {
		facets[i] = category.facetID
		idToCategory[category.id] = category
		for _, value := range category.values {
			idToValue[value.id] = value
		}
	}
	return Filters{
		categories:   categories,
		facets:       facets,
		idToCategory: idToCategory,
		idToValue:    idToValue,
		prefix:       "par",
	}
}

//===================================================================================
// Methods
//====================================================================================

func (f Filters) ParseURLQuery(query url.Values) (expression, error) {
	var result expression
	conditions := make([]condition, 0, len(query))
	for k, v := range query {
		if strings.HasPrefix(k, f.prefix) {
			cond, err := f.parseParams(k, v)
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, cond)
		}
	}
	result = expression(conditions)
	return result, nil
}

func (f Filters) parseParams(k string, v []string) (condition, error) {
	var (
		result condition
	)
	categoryID := k[len(f.prefix):]
	category, ok := f.idToCategory[categoryID]
	if !ok {
		return result, fmt.Errorf("category %s not found", categoryID)
	}
	result.param = category.facetID
	result.values = make([]string, len(v))
	for i, value := range v {
		valueObj, ok := f.idToValue[value]
		if !ok {
			return result, fmt.Errorf("value %s not found", value)
		}
		result.values[i] = valueObj.facetID
	}
	return result, nil
}
