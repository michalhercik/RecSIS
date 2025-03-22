package courses

import (
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/courses/internal/filter"
)

type Request struct {
	sessionID   string
	query       string
	indexUID    string
	page        int
	hitsPerPage int
	lang        Language
	filter      filter.Expression
}

type Response struct {
	TotalHits         int
	TotalPages        int
	Courses           []string
	FacetDistribution map[string]map[int]int
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
	lang     Language
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
}

type MeiliSearch struct {
	Client meilisearch.ServiceManager
}

func (s MeiliSearch) Search(r Request) (Response, error) {
	var result Response
	index := s.Client.Index(r.indexUID)
	searchReq := &meilisearch.SearchRequest{
		Page:                 int64(r.page),
		HitsPerPage:          int64(r.hitsPerPage),
		AttributesToRetrieve: []string{"code"},
		Filter:               r.filter,
		Facets:               filter.SliceOfParamStr(),
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
	case cs:
		result.AttributesToRetrieve = []string{"code", "cs.NAME"}
	case en:
		result.AttributesToRetrieve = []string{"code", "en.NAME"}
	default:
		return result, fmt.Errorf("SearchRequest: unsupported language: %v", r.lang)
	}
	return result, nil
}
