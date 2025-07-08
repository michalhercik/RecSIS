package coursedetail

import (
	"encoding/json"
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type Search struct {
	Client meilisearch.ServiceManager
	Survey meilisearch.IndexConfig
}

// TODO: write own meilisearch client
func (s Search) comments(r request) (response, error) {
	t := texts[r.lang]
	var result response
	searchReq := makeMultiSearchRequest(r, s.Survey)
	response, err := s.Client.MultiSearch(searchReq)
	if err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			t.errCannotSearchForSurvey,
		)
	}
	rawResponse, err := response.MarshalJSON()
	if err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			t.errCannotSearchForSurvey,
		)
	}
	multi := multiResponse{}
	if err = json.Unmarshal(rawResponse, &multi); err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			t.errCannotSearchForSurvey,
		)
	}
	result = multi.Results[0]
	for _, res := range multi.Results[1:] {
		for param, distribution := range res.FacetDistribution {
			result.FacetDistribution[param] = distribution
		}
	}
	return result, nil
}

type expression interface {
	String() string
	Except() func(func(string, string) bool)
	ConditionsCount() int
	Append(param string, values ...string)
}

type request struct {
	userID string
	query  string
	offset int
	limit  int
	lang   language.Language
	filter expression
	facets []string
	sort   string
}

type response struct {
	EstimatedTotalHits int                       `json:"estimatedTotalHits"`
	Survey             []survey                  `json:"Hits"`
	FacetDistribution  map[string]map[string]int `json:"FacetDistribution"`
}

type multiResponse struct {
	Results []response `json:"results"`
}

func makeMultiSearchRequest(r request, index meilisearch.IndexConfig) *meilisearch.MultiSearchRequest {
	numOfReq := 1 + r.filter.ConditionsCount()
	result := &meilisearch.MultiSearchRequest{
		Queries: make([]*meilisearch.SearchRequest, 0, numOfReq),
	}
	result.Queries = append(result.Queries, &meilisearch.SearchRequest{
		IndexUID:             index.Uid,
		Query:                r.query,
		Limit:                int64(r.limit),
		Offset:               int64(r.offset),
		AttributesToRetrieve: attributesToRetrieve(r.lang),
		Filter:               r.filter.String(),
		Facets:               r.facets,
		Sort:                 []string{r.sort},
	})
	for param, filter := range r.filter.Except() {
		_ = param
		_ = filter
		result.Queries = append(result.Queries, &meilisearch.SearchRequest{
			IndexUID:             index.Uid,
			Query:                r.query,
			Limit:                0,          // TODO: not working, probably bug in meilisearch-go -> write own client...
			AttributesToRetrieve: []string{}, // TODO: not working, probably bug in meilisearch-go -> write own client...
			Filter:               filter,
			Facets:               []string{param},
		})
	}
	return result
}

// TODO: remove params package
func attributesToRetrieve(lang language.Language) []string {
	var (
		studyTypeName  string
		studyTypeAbbr  string
		studyFieldName string
	)
	if lang == language.CS {
		studyTypeName = "study_type.name.cs"
		studyTypeAbbr = "study_type.abbr.cs"
		studyFieldName = "study_field.name.cs"
	} else {
		studyTypeName = "study_type.name.en"
		studyTypeAbbr = "study_type.abbr.en"
		studyFieldName = "study_field.name.en"
	}
	attrs := []string{
		"content",
		"course_code",
		"study_year",
		"academic_year",
		"teacher",
		"target_type",
		"study_field.id",
		studyFieldName,
		studyTypeName,
		studyTypeAbbr,
	}
	return attrs
}
