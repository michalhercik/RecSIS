package meilisearch

import (
	"net/url"

	"github.com/michalhercik/RecSIS/internal/course/comments/search"
)

type UrlQueryParser interface {
	Parse(query url.Values) (UrlQueryParserResult, error)
}

type UrlQueryParserResult interface {
	Add(param search.Filterable, values ...string) UrlQueryParserResult
	ConditionsCount() int
	String() string
}
