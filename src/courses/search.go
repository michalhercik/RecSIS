package courses

import (
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/courses/internal/filter"
	"github.com/michalhercik/RecSIS/language"
)

type Request struct {
	userID      string
	query       string
	indexUID    string
	page        int
	hitsPerPage int
	lang        language.Language
	filter      filter.Expression
}

type Response struct {
	TotalHits         int
	TotalPages        int
	Courses           []string
	FacetDistribution map[string]map[int]int
}

type MultiResponse struct {
	Results []Response `json:"results"`
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var hit struct {
		TotalHits  int64 `json:"totalHits"`
		TotalPages int64 `json:"totalPages"`
		Hits       []struct {
			Code string `json:"code"`
		} `json:"Hits"`
		FacetDistribution map[string]map[int]int `json:"FacetDistribution"`
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

type QuickRequest struct {
	query    string
	indexUID string
	limit    int64
	offset   int64
	lang     language.Language
}

type QuickCourse struct {
	code string
	name string
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

type SearchEngine interface {
	Search(r Request) (Response, error)
	QuickSearch(r QuickRequest) (QuickResponse, error)
	FacetDistribution() (map[string]map[int]int, error)
}

type MeiliSearch struct {
	Client  meilisearch.ServiceManager
	Courses meilisearch.IndexConfig
}

func (s MeiliSearch) FacetDistribution() (map[string]map[int]int, error) {
	searchReq := &meilisearch.SearchRequest{
		Limit:  0, // TODO: not working, probably bug in meilisearch-go -> write own client...
		Facets: filter.SliceOfParamStr(),
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

func makeMultiSearchRequest(r Request, index meilisearch.IndexConfig) *meilisearch.MultiSearchRequest {
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
		Facets:               filter.SliceOfParamStr(),
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
			Facets:               []string{param.String()},
		})
	}
	return result
}

// TODO: write own meilisearch client
func (s MeiliSearch) Search(r Request) (Response, error) {
	var result Response
	searchReq := makeMultiSearchRequest(r, s.Courses)
	response, err := s.Client.MultiSearch(searchReq)
	if err != nil {
		return result, err
	}
	rawResponse, err := response.MarshalJSON()
	if err != nil {
		return result, err
	}
	multi := MultiResponse{}
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

func (s MeiliSearch) QuickSearch(r QuickRequest) (QuickResponse, error) {
	var result QuickResponse
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

func buildQuickSearchRequest(r QuickRequest) (*meilisearch.SearchRequest, error) {
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
