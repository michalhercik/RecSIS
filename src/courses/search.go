package courses

import (
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/language"
)

type expression interface {
	String() string
	Except() func(func(string, string) bool)
	ConditionsCount() int
}

type request struct {
	userID      string
	query       string
	indexUID    string
	page        int
	hitsPerPage int
	lang        language.Language
	filter      expression
	facets      []string
}

type response struct {
	TotalHits         int
	TotalPages        int
	Courses           []string
	FacetDistribution map[string]map[string]int
}

type multiResponse struct {
	Results []response `json:"results"`
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
	r.Courses = make([]string, len(hit.Hits))
	for i, hit := range hit.Hits {
		r.Courses[i] = hit.Code
	}
	r.FacetDistribution = hit.FacetDistribution
	return nil
}

type quickRequest struct {
	query    string
	indexUID string
	limit    int64
	offset   int64
	lang     language.Language
}

type quickCourse struct {
	code string
	name string
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
	r.courses = make([]quickCourse, len(hit.Hits))
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

type searchEngine interface {
	Search(r request) (response, error)
	QuickSearch(r quickRequest) (quickResponse, error)
	FacetDistribution() (map[string]map[int]int, error)
}

type MeiliSearch struct {
	Client  meilisearch.ServiceManager
	Courses meilisearch.IndexConfig
}

func (s MeiliSearch) FacetDistribution() (map[string]map[int]int, error) {
	searchReq := &meilisearch.SearchRequest{
		Limit:  0,          // TODO: not working, probably bug in meilisearch-go -> write own client...
		Facets: []string{}, //filter.SliceOfParamStr(),
	}
	response, err := s.Client.Index(s.Courses.Uid).Search("", searchReq)
	if err != nil {
		return nil, err
	}
	var result map[string]map[int]int
	marshalRes, err := json.Marshal(response.FacetDistribution)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(marshalRes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func makeMultiSearchRequest(r request, index meilisearch.IndexConfig) *meilisearch.MultiSearchRequest {
	numOfReq := 1 + r.filter.ConditionsCount()
	result := &meilisearch.MultiSearchRequest{
		Queries: make([]*meilisearch.SearchRequest, 0, numOfReq),
	}
	result.Queries = append(result.Queries, &meilisearch.SearchRequest{
		IndexUID:             index.Uid,
		Query:                r.query,
		Page:                 int64(r.page),
		HitsPerPage:          int64(r.hitsPerPage),
		AttributesToRetrieve: []string{"code"},
		Filter:               r.filter.String(),
		Facets:               r.facets,
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

// TODO: write own meilisearch client
func (s MeiliSearch) Search(r request) (response, error) {
	var result response
	searchReq := makeMultiSearchRequest(r, s.Courses)
	response, err := s.Client.MultiSearch(searchReq)
	if err != nil {
		return result, err
	}
	rawResponse, err := response.MarshalJSON()
	if err != nil {
		return result, err
	}
	multi := multiResponse{}
	if err = json.Unmarshal(rawResponse, &multi); err != nil {
		return result, err
	}
	result = multi.Results[0]
	for _, res := range multi.Results[1:] {
		for param, distribution := range res.FacetDistribution {
			result.FacetDistribution[param] = distribution
		}
	}
	return result, nil
}

func (s MeiliSearch) QuickSearch(r quickRequest) (quickResponse, error) {
	var result quickResponse
	index := s.Client.Index(r.indexUID)
	searchReq, err := buildQuickSearchRequest(r)
	if err != nil {
		return result, err
	}
	rawResponse, err := index.SearchRaw(r.query, searchReq)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(*rawResponse, &result); err != nil {
		return result, err
	}
	return result, nil
}

func buildQuickSearchRequest(r quickRequest) (*meilisearch.SearchRequest, error) {
	result := &meilisearch.SearchRequest{
		Limit:  r.limit,
		Offset: r.offset,
	}
	switch r.lang {
	case language.CS:
		result.AttributesToRetrieve = []string{"code", "cs.NAME"}
	case language.EN:
		result.AttributesToRetrieve = []string{"code", "en.NAME"}
	default:
		return result, fmt.Errorf("SearchRequest: unsupported language: %v", r.lang)
	}
	return result, nil
}
