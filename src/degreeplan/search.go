package degreeplan

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/errorx"
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

func (s MeiliSearch) QuickSearch(r quickRequest, t text) (quickResponse, error) {
	var result quickResponse
	index := s.Client.Index(s.DegreePlans.Uid)
	searchReq := &meilisearch.SearchRequest{
		Limit:                r.limit,
		Offset:               0,
		AttributesToRetrieve: []string{"SPLAN", "NAZEV", "ZKRATKA"},
	}
	rawResponse, err := index.SearchRaw(r.query, searchReq)
	if err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("quick search failed: %w", err), errorx.P("query", r.query), errorx.P("limit", r.limit)),
			http.StatusInternalServerError,
			t.errFailedDPSearch,
		)
	}
	if err = json.Unmarshal(*rawResponse, &result); err != nil {
		return result, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to unmarshal quick search response: %w", err), errorx.P("query", r.query), errorx.P("limit", r.limit)),
			http.StatusInternalServerError,
			t.errFailedDPSearch,
		)
	}
	return result, nil
}
