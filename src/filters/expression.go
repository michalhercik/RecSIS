package filters

import (
	"fmt"
	"strings"
)

type expression []condition

func (e *expression) Append(param string, values ...string) {
	c := condition{
		param:  param,
		values: values,
	}
	*e = append(*e, c)
}

func (e expression) String() string {
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

func (e expression) Except() func(func(string, string) bool) {
	return e.except
}

func (e expression) except(yield func(string, string) bool) {
	if len(e) == 0 {
		return
	}
	if len(e) == 1 {
		yield(e[0].param, "")
		return
	}
	exceptExpr := make(expression, 0, len(e)-1)
	for i := range e.ConditionsCount() {
		exceptExpr = append(exceptExpr, e[:i]...)
		exceptExpr = append(exceptExpr, e[i+1:]...)
		if !yield(e[i].param, exceptExpr.String()) {
			return
		}
		exceptExpr = exceptExpr[:0]
	}
}

func (e expression) ConditionsCount() int {
	return len(e)
}

type condition struct {
	param  string
	values []string
}

func (c condition) String() string {
	result := fmt.Sprintf("%s IN [\"%s\"]", c.param, strings.Join(c.values, "\",\""))
	return result
}
