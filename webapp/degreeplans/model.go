package degreeplans

import (
	"fmt"

	"github.com/michalhercik/RecSIS/filters"
)

//================================================================================
// Constants
//================================================================================

const searchDegreePlanName = "search-dp-query"

const (
	CompareUrlParam = "cmp"
	dpCode          = "dpCode"
)

const unlimitedYear = 9999

const SearchIndex = "degree-plans"

const (
	facultyFacetID   = "faculty"
	sectionFacetID   = "section"
	studyTypeFacetID = "study_type"
	languageFacetID  = "teaching_lang"
	validityFacetID  = "validity"
	fieldFacetID     = "field.code"
)

const comparePrefix = "/compare/"

//================================================================================
// Data Types and Methods
//================================================================================

type degreePlanSearchPage struct {
	filters      map[string]filters.FacetIterator
	results      []degreePlanSearchResult
	searchQuery  string
	selectedPlan selectedPlan
}

type selectedPlan struct {
	isAnySelected bool
	code          string
}

type degreePlanSearchResult struct {
	code      string
	title     string
	studyType string
	validFrom dpYear
	validTo   dpYear
}

type dpYear int

func (ny dpYear) String() string {
	if int(ny) == unlimitedYear {
		return ""
	} else {
		return fmt.Sprintf("%d", int(ny))
	}
}
