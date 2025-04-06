package meilisearch

import (
	"fmt"
	"net/url"

	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch/params"
	"github.com/michalhercik/RecSIS/internal/course/comments/search"
	"github.com/michalhercik/RecSIS/internal/interface/course"
	"github.com/michalhercik/RecSIS/language"
)

type searchResponse struct {
	Result       []course.Comment         `json:"hits"`
	ResultFacets search.FacetDistribution `json:"facetDistribution"`
}

func (s searchResponse) Comments() []course.Comment {
	return s.Result
}
func (s searchResponse) Facets() search.FacetDistribution {
	return s.ResultFacets
}

type searchRequest struct {
	filterParser UrlQueryParser
	courseParam  search.Filterable
	teacherParam search.Filterable
	query        string
	filter       UrlQueryParserResult
	sort         []string
	limit        int
	offset       int
	lang         language.Language
}

func (s searchRequest) ParseURLQuery(query url.Values) (search.SearchRequestBuilder, error) {
	var err error
	s.filter, err = s.filterParser.Parse(query)
	return s, err
}
func (s searchRequest) SetQuery(query string) search.SearchRequestBuilder {
	s.query = query
	return s
}
func (s searchRequest) AddCourse(code string) search.SearchRequestBuilder {
	s.filter = s.filter.Add(s.courseParam, code)
	return s
}
func (s searchRequest) AddTeacher(code string) search.SearchRequestBuilder {
	s.filter = s.filter.Add(s.teacherParam, code)
	return s
}
func (s searchRequest) AddSort(param search.Sortable, how search.SortHow) search.SearchRequestBuilder {
	s.sort = append(s.sort, fmt.Sprintf("%s:%s", param, how))
	return s
}
func (s searchRequest) SetOffset(offset int) search.SearchRequestBuilder {
	s.offset = offset
	return s
}
func (s searchRequest) SetLimit(limit int) search.SearchRequestBuilder {
	s.limit = limit
	return s
}
func (s searchRequest) Build() (search.SearchRequest, error) {
	return s, nil
}
func (s searchRequest) Sort() []string {
	return s.sort
}
func (s searchRequest) Filter() string {
	return s.filter.String()
}
func (s searchRequest) Query() string {
	return s.query
}
func (s searchRequest) Attributes() []string {
	var studyTypeName search.Parameter
	if s.lang == language.CS {
		studyTypeName = params.StudyTypeNameCS
	} else {
		studyTypeName = params.StudyTypeNameEN
	}
	attrs := []string{
		params.Content.String(),
		params.CourseCode.String(),
		params.StudyYear.String(),
		params.AcademicYear.String(),
		params.StudyField.String(),
		params.Teacher.String(),
		params.TargetType.String(),
		params.StudyTypeCode.String(),
		// params.StudyTypeAbbr.String(),
		studyTypeName.String(),
	}
	return attrs
}
func (s searchRequest) Facets() []string {
	var studyTypeName search.Parameter
	if s.lang == language.CS {
		studyTypeName = params.StudyTypeNameCS
	} else {
		studyTypeName = params.StudyTypeNameEN
	}
	facets := []string{
		params.StudyYear.String(),
		params.AcademicYear.String(),
		params.StudyField.String(),
		params.TeacherCode.String(),
		params.TargetType.String(),
		params.StudyTypeCode.String(),
		studyTypeName.String(),
	}
	return facets
}
func (s searchRequest) Offset() int {
	return s.offset
}
func (s searchRequest) Limit() int {
	return s.limit
}
