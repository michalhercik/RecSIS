package filter

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	String() string
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

func (f FacetDistribution) ParamDistribution(p param) func(func(facet) bool) {
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
type param int

const (
	StartSemester param = iota
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

func (p *param) Scan(value interface{}) error {
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
	result = makeExpression(conditions)
	return result, nil
}

func Param(id int) (param, error) {
	result := param(id)
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
	result = makeCondition(par, v)
	return result, nil
}

type condition string

func makeCondition(p param, v []string) condition {
	return condition(fmt.Sprintf("%s IN [%s]", p, strings.Join(v, ",")))
}

type expression string

func (e expression) String() string {
	return string(e)
}

func makeExpression(conditions []condition) expression {
	var sb strings.Builder
	sb.WriteString(string(conditions[0]))
	for _, c := range conditions[1:] {
		sb.WriteString(" AND ")
		sb.WriteString(string(c))
	}
	return expression(sb.String())
}
