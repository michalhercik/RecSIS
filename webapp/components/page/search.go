package page

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type MeiliSearch struct {
	Client meilisearch.ServiceManager
	Index  string
	Limit  int64
}

func (m MeiliSearch) quickSearchResult(query string, lang language.Language) ([]quickCourse, error) {
	var result quickResponse
	t := texts[lang]
	index := m.Client.Index(m.Index)
	searchReq, err := m.buildQuickSearchRequest(lang)
	if err != nil {
		return nil, errorx.AddContext(err)
	}
	rawResponse, err := index.SearchRaw(query, searchReq)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("query", query), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			t.errQuickSearchFailed,
		)
	}
	if err = json.Unmarshal(*rawResponse, &result); err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("query", query), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			t.errQuickSearchFailed,
		)
	}
	return result.courses, nil
}

func (m MeiliSearch) buildQuickSearchRequest(lang language.Language) (*meilisearch.SearchRequest, error) {
	result := &meilisearch.SearchRequest{
		Limit: m.Limit,
	}
	switch lang {
	case language.CS:
		result.AttributesToRetrieve = []string{"code", "title.cs"}
	case language.EN:
		result.AttributesToRetrieve = []string{"code", "title.en"}
	default:
		return result, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unsupported language: %s", lang)),
			http.StatusBadRequest,
			texts[language.EN].errUnsupportedLanguage,
		)
	}
	return result, nil
}

type quickResponse struct {
	approxHits int
	courses    []quickCourse
}

func (r *quickResponse) UnmarshalJSON(data []byte) error {
	var hit struct {
		ApproxHits int64 `json:"approxHits"`
		Hits       []struct {
			Code string `json:"code"`
			Name struct {
				CS string `json:"cs"`
				EN string `json:"en"`
			} `json:"title"`
		} `json:"Hits"`
	}
	if err := json.Unmarshal(data, &hit); err != nil {
		return err
	}
	r.approxHits = int(hit.ApproxHits)
	r.courses = make([]quickCourse, len(hit.Hits))
	for i, hit := range hit.Hits {
		r.courses[i].code = hit.Code
		r.courses[i].name = hit.Name.CS
		if r.courses[i].name == "" {
			r.courses[i].name = hit.Name.EN
		}
	}
	return nil
}

type quickCourse struct {
	code string
	name string
}
