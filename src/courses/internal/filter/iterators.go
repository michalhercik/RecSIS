package filter

import (
	"iter"
	"net/url"
	"slices"

	"github.com/michalhercik/RecSIS/language"
)

func IterFiltersWithFacets(filters Filters, facets Facets, query url.Values, lang language.Language) iter.Seq[FacetIterator] {
	return func(yield func(FacetIterator) bool) {
		for _, c := range filters.categories {
			f := facets[c.facetID]
			checked := query["par"+c.id]
			result := FacetIterator{
				ID:      c.id,
				Title:   c.Title(lang),
				Desc:    c.Desc(lang),
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

func (ci FacetIterator) IterWithFacets() iter.Seq2[int, FacetValue] {
	return func(yield func(int, FacetValue) bool) {
		for i, v := range ci.filter.values {
			count, ok := ci.facets[v.facetID]
			if !ok {
				count = 0
			}
			checked := slices.Contains(ci.checked, v.id)
			result := FacetValue{
				ID:      v.id,
				Title:   v.Title(ci.lang),
				Desc:    v.Desc(ci.lang),
				Count:   count,
				Checked: checked,
			}
			if !yield(i, result) {
				return
			}
		}
	}
}

func (ci FacetIterator) Size() int {
	return len(ci.filter.values)
}

type Facets map[string]map[string]int
type FacetCategory map[string]int

type FacetIterator struct {
	ID      string
	Title   string
	Desc    string
	filter  FilterCategory
	facets  FacetCategory
	lang    language.Language
	checked []string
}

type FacetValue struct {
	ID      string
	Title   string
	Desc    string
	Count   int
	Checked bool
}
