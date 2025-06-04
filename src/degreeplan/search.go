package degreeplan

import (
	"encoding/json"

	"github.com/meilisearch/meilisearch-go"
)

type QuickRequest struct {
	query string
	limit int64
}

type QuickDegreePlan struct {
	Code string `json:"SPLAN"`
	Name string `json:"NAZEV"`
	Type string `json:"ZKRATKA"`
}

type QuickResponse struct {
	ApproxHits  int               `json:"approxHits"`
	DegreePlans []QuickDegreePlan `json:"Hits"`
}

type SearchEngine interface {
	QuickSearch(r QuickRequest) (QuickResponse, error)
}

type MeiliSearch struct {
	Client      meilisearch.ServiceManager
	DegreePlans meilisearch.IndexConfig
}

func (s MeiliSearch) QuickSearch(r QuickRequest) (QuickResponse, error) {
	var result QuickResponse
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
