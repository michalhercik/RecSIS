package filters

/** PACKAGE DESCRIPTION

The filters package provides a flexible system for managing and applying search filters in the application. Its main purpose is to fetch filter categories and values from the database, build filter structures, and parse user-selected filter parameters from HTTP requests. This enables users to refine search results using facets like course type, department, or other attributes, with support for localization and error handling.

Typical usage involves creating a Filters instance using MakeFilters and injecting it into the server, initializing it with Init() to load filter data, and then calling ParseURLQuery to convert URL query parameters into filter expressions for search. The package automatically maps filter categories and values, handles invalid or missing filter selections gracefully, and integrates with the application's error reporting and localization.

*/

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/filters/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

const (
	prefix = "par"
)

type Filters struct {
	filters
}

type filters struct {
	source       *sqlx.DB
	id           string
	categories   []category
	facets       []string
	idToCategory map[string]category
	idToValue    map[string]value
}

func MakeFilters(source *sqlx.DB, id string) Filters {
	return Filters{
		filters: filters{
			source: source,
			id:     id,
		},
	}
}

func (f *filters) Init() error {
	var rec []record
	if err := f.source.Select(&rec, sqlquery.SelectCategoriesAndValues, f.id); err != nil {
		return fmt.Errorf("failed to fetch filters: %w", err)
	}
	f.categories = fromRecords(rec)
	f.idToCategory = initIDToCategory(f.categories)
	f.idToValue = initIDToValue(f.categories)
	f.facets = initFacets(f.categories)
	return nil
}

func (f filters) Facets() []string {
	return f.facets
}

func (f filters) ParseURLQuery(query url.Values, lang language.Language) (expression, error) {
	var result expression
	conditions := make([]condition, 0, len(query))
	for k, v := range query {
		if strings.HasPrefix(k, prefix) {
			cond, err := f.parseParams(k, v, lang)
			if err != nil {
				return nil, errorx.AddContext(err)
			}
			conditions = append(conditions, cond)
		}
	}
	result = expression(conditions)
	return result, nil
}

func (f filters) parseParams(k string, v []string, lang language.Language) (condition, error) {
	t := texts[lang]
	var result condition
	categoryID := k[len(prefix):]
	category, ok := f.idToCategory[categoryID]
	if !ok {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(
				fmt.Errorf("category %s not found", categoryID),
				errorx.P("k", k),
				errorx.P("v", strings.Join(v, ",")),
				errorx.P("categoryID", categoryID),
			),
			http.StatusBadRequest,
			t.errCategoryNotFound,
		)
	}
	result.param = category.facetID
	result.values = make([]string, len(v))
	for i, value := range v {
		valueObj, ok := f.idToValue[value]
		if !ok {
			return result, errorx.NewHTTPErr(
				errorx.AddContext(
					fmt.Errorf("value %s not found in category %s", value, categoryID),
					errorx.P("k", k),
					errorx.P("v", strings.Join(v, ",")),
					errorx.P("categoryID", categoryID),
					errorx.P("value", value),
				),
				http.StatusBadRequest,
				t.errValueNotFound,
			)
		}
		result.values[i] = valueObj.facetID
	}
	return result, nil
}

func fromRecords(rec []record) []category {
	fb := categoryBuilder{}
	for _, row := range rec {
		if fb.isLastCategory(row.CategoryID) {
			fb.category(makeFilterIdentity(
				row.CategoryID,
				row.CategoryFacetID,
				language.MakeLangString(row.CategoryTitleCS, row.CategoryTitleEN),
				language.MakeLangString(row.CategoryDescCS.String, row.CategoryDescEN.String),
			), row.CategoryDisplayedValueLimit)
		}
		if row.ValueID.Valid {
			fb.value(makeFilterIdentity(
				row.ValueID.String,
				row.ValueFacetID.String,
				language.MakeLangString(row.ValueTitleCS.String, row.ValueTitleEN.String),
				language.MakeLangString(row.ValueDescCS.String, row.ValueDescEN.String),
			))
		}
	}
	return fb.build()
}

func initIDToCategory(categories []category) map[string]category {
	idToCategory := make(map[string]category)
	for _, category := range categories {
		idToCategory[category.id] = category
	}
	return idToCategory
}

func initIDToValue(categories []category) map[string]value {
	idToValue := make(map[string]value)
	for _, category := range categories {
		for _, value := range category.values {
			idToValue[value.id] = value
		}
	}
	return idToValue
}

func initFacets(categories []category) []string {
	facets := make([]string, len(categories))
	for i, category := range categories {
		facets[i] = category.facetID
	}
	return facets
}

type record struct {
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
}

type category struct {
	identity
	displayedValueLimit int
	values              []value
}

func (fc category) Values() []value {
	return fc.values
}

func (fc category) DisplayedValueLimit() int {
	return fc.displayedValueLimit
}

type categoryBuilder struct {
	categ []category
}

func (fb categoryBuilder) isLastCategory(categoryID string) bool {
	if len(fb.categ) == 0 {
		return true
	}
	lastID := fb.categ[len(fb.categ)-1].id
	return lastID != categoryID
}

func (fb *categoryBuilder) category(identity identity, displayedValueLimit int) {
	fb.categ = append(fb.categ, category{
		identity:            identity,
		displayedValueLimit: displayedValueLimit,
		values:              []value{},
	})
}

func (fb *categoryBuilder) value(identity identity) {
	value := makeFilterValue(identity)
	category := fb.categ[len(fb.categ)-1]
	category.values = append(category.values, value)
	fb.categ[len(fb.categ)-1] = category
}

func (fb *categoryBuilder) build() []category {
	return fb.categ
}

type value struct {
	identity
}

func makeFilterValue(identity identity) value {
	return value{
		identity: identity,
	}
}

type identity struct {
	id      string
	facetID string
	title   language.LangString
	desc    language.LangString
}

func (fi identity) Title(lang language.Language) string {
	return fi.title.String(lang)
}

func (fi identity) Desc(lang language.Language) string {
	return fi.desc.String(lang)
}

func (fi identity) ID() string {
	return fi.id
}

func makeFilterIdentity(id, facetID string, title, desc language.LangString) identity {
	return identity{
		id:      id,
		facetID: facetID,
		title:   title,
		desc:    desc,
	}
}
