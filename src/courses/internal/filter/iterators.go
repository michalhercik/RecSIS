package filter

import (
	"iter"

	"github.com/michalhercik/RecSIS/language"
)

func IterFiltersWithFacets(filters Filters, facets Facets, lang language.Language) iter.Seq[FacetIterator] {
	return func(yield func(FacetIterator) bool) {
		for _, c := range filters.categories {
			f := facets[c.facetID]
			result := FacetIterator{
				ID:     c.ID(),
				Title:  c.Title(lang),
				Desc:   c.Desc(lang),
				filter: c,
				facets: f,
				lang:   lang,
			}
			if !yield(result) {
				return
			}
		}
	}
}

func (ci FacetIterator) IterWithFacets() iter.Seq[FacetValue] {
	return func(yield func(FacetValue) bool) {
		for _, v := range ci.filter.values {
			count, ok := ci.facets[v.facetID]
			if !ok {
				count = 0
			}
			result := FacetValue{
				ID:    v.id,
				Title: v.Title(ci.lang),
				Desc:  v.Desc(ci.lang),
				Count: count,
			}
			if !yield(result) {
				return
			}
		}
	}
}

type Facets map[string]map[string]int
type FacetCategory map[string]int

type FacetIterator struct {
	ID     string
	Title  string
	Desc   string
	filter FilterCategory
	facets FacetCategory
	lang   language.Language
}

type FacetValue struct {
	ID    string
	Title string
	Desc  string
	Count int
}
