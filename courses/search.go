package courses

import (
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

type Language int

const (
	czech Language = iota
	english
)

type Request struct {
	query       string
	indexUID    string
	page        int64
	hitsPerPage int64
	lang        Language
	sortedBy    sortType
	semester    TeachingSemester
}

type Response struct {
	totalHits  int64
	totalPages int64
	courses    []Course
}

type QuickRequest struct {
	query    string
	indexUID string
	limit    int64
	offset   int64
	lang     Language
}

type QuickResponse struct {
	approxHits int64
	courses    []Course
}

type SearchEngine interface {
	Search(r *Request) (*Response, error)
	QuickSearch(r *QuickRequest) (*QuickResponse, error)
}

type MeiliSearch struct {
	Client meilisearch.ServiceManager
}

func (s MeiliSearch) Search(r *Request) (*Response, error) {
	index := s.Client.Index(r.indexUID)
	searchReq, err := buildSearchRequest(r)
	if err != nil {
		return nil, err
	}
	searchRes, err := index.Search(r.query, searchReq)
	if err != nil {
		return nil, err
	}
	courses, err := parseCourses(searchRes.Hits)
	if err != nil {
		return nil, err
	}
	res := Response{
		totalHits:  searchRes.TotalHits,
		totalPages: searchRes.TotalPages,
		courses:    courses,
	}
	return &res, nil
}

func (s MeiliSearch) QuickSearch(r *QuickRequest) (*QuickResponse, error) {
	index := s.Client.Index(r.indexUID)
	searchReq, err := buildQuickSearchRequest(r)
	if err != nil {
		return nil, err
	}
	searchRes, err := index.Search(r.query, searchReq)
	if err != nil {
		return nil, err
	}
	courses, err := parseCourses(searchRes.Hits)
	if err != nil {
		return nil, err
	}
	res := QuickResponse{
		approxHits: searchRes.EstimatedTotalHits,
		courses:    courses,
	}
	return &res, nil
}

func parseCourses(hits []interface{}) ([]Course, error) {
	courses := []Course{}
	for _, hit := range hits {
		course := Course{}
		payload, _ := json.Marshal(hit)
		err := json.Unmarshal(payload, &course)
		if err != nil {
			return courses, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func buildSearchRequest(r *Request) (*meilisearch.SearchRequest, error) {
	searchReq := &meilisearch.SearchRequest{
		Page:        r.page,
		HitsPerPage: r.hitsPerPage,
		AttributesToRetrieve: []string{
			"code",
			"start",
			"semesterCount",
			"lectureRange1",
			"seminarRange1",
			"lectureRange2",
			"seminarRange2",
			"examType",
			"credits",
			"teacher1Id",
			"teacher1Firstname",
			"teacher1Lastname",
			"teacher2Id",
			"teacher2Firstname",
			"teacher2Lastname",
			"teacher3Id",
			"teacher3Firstname",
			"teacher3Lastname",
			"rating",
		},
	}
	switch r.lang {
	case czech:
		searchReq.AttributesToRetrieve = append(searchReq.AttributesToRetrieve, []string{
			"nameCs",
			"annotationCs",
		}...)
	case english:
		searchReq.AttributesToRetrieve = append(searchReq.AttributesToRetrieve, []string{
			"nameEn",
			"annotationEn",
		}...)
	default:
		return searchReq, fmt.Errorf("SearchRequest: unsupported language: %d", r.lang)
	}
	return searchReq, nil
}

func buildQuickSearchRequest(r *QuickRequest) (*meilisearch.SearchRequest, error) {
	searchReq := &meilisearch.SearchRequest{
		Limit:  r.limit,
		Offset: r.offset,
		AttributesToRetrieve: []string{
			"code",
		},
	}
	switch r.lang {
	case czech:
		searchReq.AttributesToRetrieve = append(searchReq.AttributesToRetrieve, "nameCs")
	case english:
		searchReq.AttributesToRetrieve = append(searchReq.AttributesToRetrieve, "nameEn")
	default:
		return searchReq, fmt.Errorf("SearchRequest: unsupported language: %d", r.lang)
	}
	return searchReq, nil
}
