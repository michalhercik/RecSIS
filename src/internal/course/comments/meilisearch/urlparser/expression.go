package urlparser

import (
	"fmt"
	"strings"

	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch"
	"github.com/michalhercik/RecSIS/internal/course/comments/search"
)

type Condition struct {
	param  search.Parameter
	values []string
}

func (c Condition) String() string {
	return fmt.Sprintf("%s IN [%s]", c.param, strings.Join(c.values, ","))
}

type Expression []Condition

func (e Expression) String() string {
	var sb strings.Builder
	if len(e) == 0 {
		return ""
	}
	sb.WriteString(e[0].String())
	for _, c := range e[1:] {
		sb.WriteString(" AND ")
		sb.WriteString(c.String())
	}
	return sb.String()
}

func (e Expression) Add(param search.Filterable, values ...string) meilisearch.UrlQueryParserResult {
	// TODO: implement
	e = append(e, Condition{param, values})
	return e
}

func (e Expression) ConditionsCount() int {
	return len(e)
}
