package params

import "github.com/michalhercik/RecSIS/internal/course/comments/search"

var parFac = makeParamFactory(serialID())

// For retrieving and filtering
var (
	AcademicYear  = parFac.makeParameter("academic_year") // for sorting as well
	CourseCode    = parFac.makeFilterableOnly("course_code")
	StudyYear     = parFac.makeFilterableOnly("study_year")
	Semester      = parFac.makeFilterableOnly("semester")
	StudyField    = parFac.makeFilterableOnly("study_field")
	TeacherCode   = parFac.makeFilterableOnly("teacher.KOD")
	StudyTypeCode = parFac.makeFilterableOnly("study_type.code")
	TargetType    = parFac.makeFilterableOnly("target_type")
)

// Only for retrieving data
var (
	Teacher         = parFac.makeParameter("teacher")
	StudyTypeNameCS = parFac.makeParameter("study_type.name_cs")
	StudyTypeNameEN = parFac.makeParameter("study_type.name_en")
	Content         = parFac.makeParameter("content")
	StudyTypeAbbr   = parFac.makeParameter("study_type.abbr")
)

// For parsing url query
var (
	StringToParam = parFac.StringToParam()
	IdToParam     = parFac.IDToParam()
)

var sortFac = makeParamFactory(serialID())
var (
	Asc  = sortFac.makeParameter("asc")
	Desc = sortFac.makeParameter("desc")
)

func serialID() func() int {
	counter := 0
	return func() int {
		counter++
		return counter
	}
}

type Parameter struct {
	id    int
	label string
}

func (p Parameter) ID() int {
	return p.id
}
func (p Parameter) String() string {
	return p.label
}

type ParamFactory struct {
	stringToParam map[string]search.Parameter
	idToParam     map[int]search.Parameter
	idgen         func() int
}

func makeParamFactory(idgen func() int) ParamFactory {
	return ParamFactory{
		stringToParam: make(map[string]search.Parameter),
		idToParam:     make(map[int]search.Parameter),
		idgen:         idgen,
	}
}
func (pf *ParamFactory) makeParameter(label string) search.Parameter {
	p := Parameter{id: pf.idgen(), label: label}
	pf.stringToParam[label] = p
	pf.idToParam[p.id] = p
	return p
}
func (pf *ParamFactory) makeFilterableOnly(label string) search.Filterable {
	p := Parameter{id: pf.idgen(), label: label}
	pf.stringToParam[label] = p
	pf.idToParam[p.id] = p
	return search.Filterable(p)
}
func (pf *ParamFactory) IDToParam() map[int]search.Parameter {
	return pf.idToParam
}
func (pf *ParamFactory) StringToParam() map[string]search.Parameter {
	return pf.stringToParam
}

// func (pf *ParamFactory) makeSortable(label string) search.Sortable {
// 	p := Parameter{id: pf.idgen(), label: label}
// 	pf.stringToParam[label] = p
// 	pf.idToParam[p.id] = p
// 	return search.Sortable(p)
// }
