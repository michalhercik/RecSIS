package filters

import (
	"fmt"
	"strings"
)

type expression []condition

func (e *expression) Append(param string, values ...string) {
	c := inCondition{
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
		yield(e[0].getParam(), "")
		return
	}
	exceptExpr := make(expression, 0, len(e)-1)
	for i := range e.ConditionsCount() {
		exceptExpr = append(exceptExpr, e[:i]...)
		exceptExpr = append(exceptExpr, e[i+1:]...)
		if !yield(e[i].getParam(), exceptExpr.String()) {
			return
		}
		exceptExpr = exceptExpr[:0]
	}
}

func (e expression) ConditionsCount() int {
	return len(e)
}

// condition interface
type condition interface {
	getParam() string
	String() string
}

// general IN condition
type inCondition struct {
	param  string
	values []string
}

func (c inCondition) getParam() string {
	return c.param
}

func (c inCondition) String() string {
	result := fmt.Sprintf("%s IN [\"%s\"]", c.param, strings.Join(c.values, "\",\""))
	return result
}

// custom condition
type customCondition struct {
	condition string
	param     string
	values    []string
}

func (c customCondition) getParam() string {
	return c.param
}

func (c customCondition) String() string {
	var conditions []string
	for _, v := range c.values {
		conditions = append(conditions, strings.ReplaceAll(c.condition, "{VAL}", v))
	}
	return fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
}
