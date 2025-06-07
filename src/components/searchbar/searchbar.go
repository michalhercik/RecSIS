package searchbar

import (
	"encoding/json"
	"fmt"

	"github.com/a-h/templ"
	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/language"
)

type QuickCourse struct {
	code string
	name string
}

type MeiliSearch struct {
	Client                meilisearch.ServiceManager
	Index                 string
	Limit                 int64
	SearchBarView         func(searchBarModel) templ.Component
	SearchResultsView     func(quickResultsModel) templ.Component
	ResultsDetailEndpoint func(code string) string
	SearchEndpoint        string
	FiltersSelector       string
	QuickEndpoint         string
	Param                 string
}

func (m MeiliSearch) SearchParam() string {
	return m.Param
}

func (m MeiliSearch) QuickSearchEndpoint() string {
	return m.QuickEndpoint
}

func (m MeiliSearch) View(searchInput string, lang language.Language, includeFilters bool) templ.Component {
	model := searchBarModel{
		t:                   texts[lang],
		lang:                lang,
		searchInput:         searchInput,
		searchParam:         m.Param,
		searchEndpoint:      m.SearchEndpoint,  // "/courses/search",
		filtersSelector:     m.FiltersSelector, // "#filter-form",
		quickSearchEndpoint: m.QuickEndpoint,   // "page/quicksearch",
		includeFilters:      includeFilters,
	}
	return m.SearchBarView(model)
}

func (m MeiliSearch) QuickSearchResult(query string, lang language.Language) (templ.Component, error) {
	var result QuickResponse
	index := m.Client.Index(m.Index)
	searchReq, err := m.buildQuickSearchRequest(lang)
	if err != nil {
		return nil, err
	}
	rawResponse, err := index.SearchRaw(query, searchReq)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(*rawResponse, &result); err != nil {
		return nil, err
	}
	t := texts[lang]
	model := quickResultsModel{
		t:                    t,
		lang:                 lang,
		courses:              result.courses,
		resultDetailEndpoint: m.ResultsDetailEndpoint,
	}
	return m.SearchResultsView(model), nil
}

func (m MeiliSearch) buildQuickSearchRequest(lang language.Language) (*meilisearch.SearchRequest, error) {
	result := &meilisearch.SearchRequest{
		Limit: m.Limit,
	}
	switch lang {
	case language.CS:
		result.AttributesToRetrieve = []string{"code", "cs.NAME"}
	case language.EN:
		result.AttributesToRetrieve = []string{"code", "en.NAME"}
	default:
		return result, fmt.Errorf("SearchRequest: unsupported language: %v", lang)
	}
	return result, nil
}

type QuickResponse struct {
	approxHits int
	courses    []QuickCourse
}

func (r *QuickResponse) UnmarshalJSON(data []byte) error {
	var hit struct {
		ApproxHits int64 `json:"approxHits"`
		Hits       []struct {
			Code string `json:"code"`
			Cs   struct {
				Name string `json:"NAME"`
			} `json:"cs"`
			En struct {
				Name string `json:"NAME"`
			} `json:"en"`
		} `json:"Hits"`
	}
	if err := json.Unmarshal(data, &hit); err != nil {
		return err
	}
	r.approxHits = int(hit.ApproxHits)
	r.courses = make([]QuickCourse, len(hit.Hits))
	for i, hit := range hit.Hits {
		r.courses[i].code = hit.Code
		if hit.Cs.Name != "" {
			r.courses[i].name = hit.Cs.Name
		} else {
			r.courses[i].name = hit.En.Name
		}
	}
	return nil
}

type searchBarModel struct {
	t                   text
	lang                language.Language
	searchInput         string
	searchParam         string
	searchEndpoint      string
	filtersSelector     string
	quickSearchEndpoint string
	includeFilters      bool
}

type quickResultsModel struct {
	t                    text
	lang                 language.Language
	courses              []QuickCourse
	resultDetailEndpoint func(code string) string
}
