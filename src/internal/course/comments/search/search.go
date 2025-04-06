package search

import (
	"net/url"

	"github.com/michalhercik/RecSIS/internal/interface/course"
	"github.com/michalhercik/RecSIS/language"
)

type SearchEngine interface {
	BuildRequest(lang language.Language) SearchRequestBuilder
	Comments(SearchRequest) (SearchResult, error)
}

type FacetDistribution map[string]map[string]int

type SearchResult interface {
	Comments() []course.Comment
	Facets() FacetDistribution
}

type SearchRequestBuilder interface {
	ParseURLQuery(query url.Values) (SearchRequestBuilder, error)
	SetQuery(query string) SearchRequestBuilder
	AddCourse(courseCode string) SearchRequestBuilder
	AddTeacher(teacherCode string) SearchRequestBuilder
	AddSort(param Sortable, how SortHow) SearchRequestBuilder
	SetOffset(offset int) SearchRequestBuilder
	SetLimit(limit int) SearchRequestBuilder
	Build() (SearchRequest, error)
}

type SearchRequest interface {
	Sort() []string
	Filter() string
	Query() string
	Attributes() []string
	Facets() []string
	Offset() int
	Limit() int
}

type Filterable interface {
	Parameter
}

type Sortable interface {
	Parameter
}

type SortHow interface {
	Parameter
}

type Parameter interface {
	ID() int
	String() string
}
