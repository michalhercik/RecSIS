package filter

import (
	"fmt"
	"net/url"
	"strings"

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

type Filters struct {
	categories   []FilterCategory
	Facets       []string
	idToCategory map[string]FilterCategory
	idToValue    map[string]FilterValue
	prefix       string
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
		Facets:       facets,
		idToCategory: idToCategory,
		idToValue:    idToValue,
		prefix:       "par",
	}
}

//===================================================================================
// Methods
//====================================================================================

func (f Filters) ParseURLQuery(query url.Values) (Expression, error) {
	var result expression
	conditions := make([]condition, 0, len(query))
	for k, v := range query {
		if strings.HasPrefix(k, f.prefix) {
			cond, err := f.parseParams(k, v)
			if err != nil {
				return result, err
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
	result.param = category
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
