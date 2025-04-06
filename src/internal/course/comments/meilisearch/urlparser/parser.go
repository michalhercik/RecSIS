package urlparser

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch"
	"github.com/michalhercik/RecSIS/internal/course/comments/search"
)

type FilterParser struct {
	ParamPrefix string
	IDToParam   map[int]search.Parameter
}

func (fp FilterParser) Parse(query url.Values) (meilisearch.UrlQueryParserResult, error) {
	var result meilisearch.UrlQueryParserResult
	conditions := make(Expression, 0, len(query))
	for k, v := range query {
		if fp.isParam(k) {
			cond, err := fp.parseParam(k, v)
			if err != nil {
				return result, err
			}
			conditions = append(conditions, cond)
		}
	}
	result = meilisearch.UrlQueryParserResult(conditions)
	return result, nil
}

func (fp FilterParser) isParam(p string) bool {
	return strings.HasPrefix(p, fp.ParamPrefix)
}

func (fp FilterParser) parseParam(k string, v []string) (Condition, error) {
	var result Condition
	parID, err := fp.keyToID(k)
	if err != nil {
		return result, err
	}
	par, ok := fp.IDToParam[parID]
	if !ok {
		return result, err
	}
	result = Condition{par, v}
	return result, nil
}

func (fp FilterParser) keyToID(k string) (int, error) {
	return strconv.Atoi(k[len(fp.ParamPrefix):])
}
