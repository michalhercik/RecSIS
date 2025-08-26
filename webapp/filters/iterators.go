package filters

import (
	"iter"
	"net/url"
	"slices"

	"github.com/michalhercik/RecSIS/language"
)

type Facets map[string]map[string]int

func (f Filters) IterFiltersWithFacets(facets Facets, query url.Values, lang language.Language) iter.Seq[FacetIterator] {
	return func(yield func(FacetIterator) bool) {
		for _, c := range f.categories {
			f := facets[c.facetID]
			checked := query[prefix+c.id]
			result := FacetIterator{
				title:   c.title.String(lang),
				desc:    c.desc.String(lang),
				count:   len(f),
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

type FacetValue struct {
	ID      string
	Title   string
	Desc    string
	Count   int
	Checked bool
	Prefix  string
}

func SkipEmptyFacet(iter iter.Seq2[int, FacetValue]) iter.Seq2[int, FacetValue] {
	return func(yield func(int, FacetValue) bool) {
		for i, v := range iter {
			if v.Count <= 0 {
				continue
			}
			if !yield(i, v) {
				return
			}
		}
	}
}

type FacetIterator struct {
	title   string
	desc    string
	count   int
	filter  category
	facets  facetCategory
	lang    language.Language
	checked []string
}

type facetCategory map[string]int

func CategoryWithAtLeast(n int, iter iter.Seq[FacetIterator]) iter.Seq[FacetIterator] {
	return func(yield func(FacetIterator) bool) {
		for c := range iter {
			if c.count < n {
				continue
			}
			if !yield(c) {
				return
			}
		}
	}
}

func (ci FacetIterator) IterWithFacets() iter.Seq2[int, FacetValue] {
	return func(yield func(int, FacetValue) bool) {
		for i, v := range ci.filter.values {
			count := ci.facets[v.facetID]
			checked := slices.Contains(ci.checked, v.id)
			result := FacetValue{
				ID:      v.id,
				Title:   v.title.String(ci.lang),
				Desc:    v.desc.String(ci.lang),
				Count:   count,
				Checked: checked,
				Prefix:  prefix,
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

func (fi FacetIterator) Count() int {
	return fi.count
}

func (fi FacetIterator) ID() string {
	return fi.filter.id
}

func (fi FacetIterator) Title() string {
	return fi.title
}

func (fi FacetIterator) Desc() string {
	return fi.desc
}

func (fi FacetIterator) DisplayedValueLimit() int {
	return fi.filter.displayedValueLimit
}

func (fi FacetIterator) Active() bool {
	active := false
	for _, v := range fi.checked {
		if v != "" {
			active = true
			break
		}
	}
	return active
}
