package meilisearch

import (
	"encoding/json"

	"github.com/meilisearch/meilisearch-go"
	"github.com/michalhercik/RecSIS/internal/course/comments/search"
	"github.com/michalhercik/RecSIS/language"
)

type MeiliSearch struct {
	Client        meilisearch.ServiceManager // TODO: make interface not a concrete type
	CommentsIndex meilisearch.IndexConfig    // TODO: make interface not a concrete type
	UrlToFilter   UrlQueryParser
	TeacherParam  search.Filterable
	CourseParam   search.Filterable
}

func (m MeiliSearch) BuildRequest(lang language.Language) search.SearchRequestBuilder {
	return searchRequest{
		lang:         lang,
		filterParser: m.UrlToFilter,
		courseParam:  m.CourseParam,
		teacherParam: m.TeacherParam,
	}
}

func (m MeiliSearch) Comments(req search.SearchRequest) (search.SearchResult, error) {
	searchReq := &meilisearch.SearchRequest{
		Query:                req.Query(),
		Offset:               int64(req.Offset()),
		Limit:                int64(req.Limit()),
		Filter:               req.Filter(),
		Sort:                 req.Sort(),
		AttributesToRetrieve: req.Attributes(),
		Facets:               req.Facets(),
	}
	searchRes, err := m.search(req.Query(), searchReq)
	if err != nil {
		return nil, err
	}
	return search.SearchResult(searchRes), nil
}

func (m MeiliSearch) search(query string, searchReq *meilisearch.SearchRequest) (searchResponse, error) {
	var result searchResponse
	index := m.Client.Index(m.CommentsIndex.Uid)
	searchResRaw, err := index.SearchRaw(query, searchReq)
	if err != nil {
		return result, err
	}
	json.Unmarshal(*searchResRaw, &result)
	return result, nil
}
