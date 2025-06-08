package degreeplan

import (
	"encoding/json"

	"github.com/meilisearch/meilisearch-go"
)

type quickRequest struct {
	query string
	limit int64
}

type quickDegreePlan struct {
	Code string `json:"SPLAN"`
	Name string `json:"NAZEV"`
	Type string `json:"ZKRATKA"`
}

type quickResponse struct {
	ApproxHits  int               `json:"approxHits"`
	DegreePlans []quickDegreePlan `json:"Hits"`
}

// type searchEngine interface {
// 	QuickSearch(r quickRequest) (quickResponse, error)
// }

type MeiliSearch struct {
	Client      meilisearch.ServiceManager
	DegreePlans meilisearch.IndexConfig
}

func (s MeiliSearch) QuickSearch(r quickRequest) (quickResponse, error) {
	var result quickResponse
	index := s.Client.Index(s.DegreePlans.Uid)
	searchReq := &meilisearch.SearchRequest{
		Limit:                r.limit,
		Offset:               0,
		AttributesToRetrieve: []string{"SPLAN", "NAZEV", "ZKRATKA"},
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
