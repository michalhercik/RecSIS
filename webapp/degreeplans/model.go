package degreeplans

import (
	"fmt"

	"github.com/michalhercik/RecSIS/filters"
)

//================================================================================
// Constants
//================================================================================

const searchDegreePlanName = "search-dp-query"

const unlimitedYear = 9999

const SearchIndex = "degree-plans"

const (
	FacultyFacetID   = "faculty"
	SectionFacetID   = "section"
	StudyTypeFacetID = "study_type"
	LanguageFacetID  = "teaching_lang"
	ValidityFacetID  = "validity"
	FieldFacetID     = "field.code"
)

//================================================================================
// Data Types and Methods
//================================================================================

type degreePlanSearchPage struct {
	filters map[string]filters.FacetIterator
	results []degreePlanSearchResult
}

type degreePlanSearchResult struct {
	Code      string `db:"plan_code"`
	Title     string `db:"title"`
	StudyType string `db:"study_type"`
	ValidFrom dpYear `db:"valid_from"`
	ValidTo   dpYear `db:"valid_to"`
}

type dpYear int

func (ny dpYear) String() string {
	if int(ny) == unlimitedYear {
		return ""
	} else {
		return fmt.Sprintf("%d", int(ny))
	}
}
