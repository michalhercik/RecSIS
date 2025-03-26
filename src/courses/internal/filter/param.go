package filter

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	String() string
	Except() func(func(Parameter, string) bool)
	ConditionsCount() int
}

type ParamValue struct {
	ID    int    `db:"id"`
	Label string `db:"label"`
}

type FacetDistribution struct {
	facetDistribution map[string]map[int]int
	valueLabels       map[string][]ParamValue
}

func MakeFacetDistribution(facets map[string]map[int]int, valueLabels map[string][]ParamValue) FacetDistribution {
	return FacetDistribution{facetDistribution: facets, valueLabels: valueLabels}
}

func (f FacetDistribution) ParamDistribution(p Parameter) func(func(facet) bool) {
	ps := p.String()
	paramDist := f.facetDistribution[ps]
	paramVal := f.valueLabels[ps]
	return func(yield func(facet) bool) {
		for _, v := range paramVal {
			count, ok := paramDist[v.ID]
			if !ok {
				count = 0
			}
			f := facet{Label: v.Label, Code: v.ID, Count: count}
			if !yield(f) {
				return
			}
		}
	}
}

type facet struct {
	Label string
	Code  int
	Count int
}

const ParamPrefix = "par"

//go:generate enumer -type=param -transform=snake
type Parameter int

const (
	StartSemester Parameter = iota
	SemesterCount
	LectureRangeWinter
	SeminarRangeWinter
	LectureRangeSummer
	SeminarRangeSummer
	Credits
	FacultyGuarantor
	ExamType
	RangeUnit
	Taught
	TaughtLang
	Faculty
	Capacity
	MinNumber
)

func (p *Parameter) Scan(value interface{}) error {
	sParam, ok := value.(string)
	if !ok {
		return fmt.Errorf("param must be a string")
	}
	param, err := paramString(sParam)
	if err != nil {
		return err
	}
	*p = param
	return nil
}

func SliceOfParamStr() []string {
	keys := make([]string, len(_paramNameToValueMap))
	i := 0
	for k := range _paramNameToValueMap {
		keys[i] = k
		i++
	}
	return keys
}

func ParseFilters(query map[string][]string) (Expression, error) {
	var result expression
	conditions := make([]condition, 0, len(query))
	for k, v := range query {
		if strings.HasPrefix(k, ParamPrefix) {
			cond, err := parseParams(k, v)
			if err != nil {
				return result, err
			}
			conditions = append(conditions, cond)
		}
	}
	result = expression(conditions)
	return result, nil
}

func Param(id int) (Parameter, error) {
	result := Parameter(id)
	if !result.IsAparam() {
		return 0, fmt.Errorf("%d does not belong to param values", id)
	}
	return result, nil
}

func parseParams(k string, v []string) (condition, error) {
	var result condition
	parID, err := strconv.Atoi(k[3:])
	if err != nil {
		return result, err
	}
	par, err := Param(parID)
	if err != nil {
		return result, err
	}
	result = condition{par, v}
	return result, nil
}

type condition struct {
	param  Parameter
	values []string
}

func (c condition) String() string {
	return fmt.Sprintf("%s IN [%s]", c.param, strings.Join(c.values, ","))
}

type expression []condition

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

func (e expression) Except() func(func(Parameter, string) bool) {
	return e.except
}

func (e expression) except(yield func(Parameter, string) bool) {
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
