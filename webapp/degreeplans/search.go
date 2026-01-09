package degreeplans

import (
	"encoding/json"
	"maps"
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type expression interface {
	String() string
	Except() func(func(string, string) bool)
	ConditionsCount() int
}

type request struct {
	userID   string
	query    string
	indexUID string
	lang     language.Language
	filter   expression
	facets   []string
}

type response struct {
	TotalHits         int
	TotalPages        int
	DegreePlanCodes   []string
	FacetDistribution map[string]map[string]int
}

func (r *response) UnmarshalJSON(data []byte) error {
	var hit struct {
		TotalHits  int64 `json:"totalHits"`
		TotalPages int64 `json:"totalPages"`
		Hits       []struct {
			Code string `json:"code"`
		} `json:"Hits"`
		FacetDistribution map[string]map[string]int `json:"FacetDistribution"`
	}
	if err := json.Unmarshal(data, &hit); err != nil {
		return err
	}
	r.TotalHits = int(hit.TotalHits)
	r.TotalPages = int(hit.TotalPages)
	r.DegreePlanCodes = make([]string, len(hit.Hits))
	for i, hit := range hit.Hits {
		r.DegreePlanCodes[i] = hit.Code
	}
	r.FacetDistribution = hit.FacetDistribution
	return nil
}

type multiResponse struct {
	Results []response `json:"results"`
}

type searchEngine interface {
	Search(r request) (response, error)
}

type MeiliSearch struct {
	Client      meilisearch.ServiceManager
	DegreePlans meilisearch.IndexConfig
}

func newRequestWithDisjunctiveFaceting(r request, index meilisearch.IndexConfig) *meilisearch.MultiSearchRequest {
	numOfReq := 1 + r.filter.ConditionsCount()
	result := &meilisearch.MultiSearchRequest{
		Queries: make([]*meilisearch.SearchRequest, 0, numOfReq),
	}
	result.Queries = append(result.Queries, &meilisearch.SearchRequest{
		IndexUID:             index.Uid,
		Query:                r.query,
		Page:                 1,
		HitsPerPage:          200,
		AttributesToRetrieve: []string{"code"}, // TODO might get all data from MeiliSearch and avoid extra DB query
		Filter:               r.filter.String(),
		Facets:               r.facets,
	})
	for param, filter := range r.filter.Except() {
		result.Queries = append(result.Queries, &meilisearch.SearchRequest{
			IndexUID:             index.Uid,
			Query:                r.query,
			Limit:                0,          // not working returns more than zero, probably bug in meilisearch-go -> write own client...
			AttributesToRetrieve: []string{}, // not working returns more than zero, probably bug in meilisearch-go -> write own client...
			Filter:               filter,
			Facets:               []string{param},
		})
	}
	return result
}

// TODO: write own meilisearch client
func (s MeiliSearch) Search(r request) (response, error) {
	t := texts[r.lang]
	var result response
	searchReq := newRequestWithDisjunctiveFaceting(r, s.DegreePlans)
	response, err := s.Client.MultiSearch(searchReq)
	if err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("index", r.indexUID), errorx.P("query", r.query)),
			http.StatusInternalServerError,
			t.errFailedDPSearch,
		)
	}
	rawResponse, err := response.MarshalJSON()
	if err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("index", r.indexUID), errorx.P("query", r.query)),
			http.StatusInternalServerError,
			t.errFailedDPSearch,
		)
	}
	multi := multiResponse{}
	if err = json.Unmarshal(rawResponse, &multi); err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("index", r.indexUID), errorx.P("query", r.query)),
			http.StatusInternalServerError,
			t.errFailedDPSearch,
		)
	}
	result = multi.Results[0]
	for _, res := range multi.Results[1:] {
		maps.Copy(result.FacetDistribution, res.FacetDistribution)
	}
	return result, nil
}
