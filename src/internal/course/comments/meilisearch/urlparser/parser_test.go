package urlparser

import (
	"testing"

	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch/params"
)

func TestQueryParsing(t *testing.T) {
	var (
		expectedResult = "academic_year IN [value1,value2] AND course_code IN [value3] AND study_year IN [value4,value5] AND course_code IN [value6]"
	)
	fp := FilterParser{
		ParamPrefix: "testparf",
		IDToParam:   params.IdToParam,
	}
	query := make(map[string][]string)
	query["testparf1"] = []string{"value1", "value2"}
	query["testparf2"] = []string{"value3"}
	query["testparf3"] = []string{"value4", "value5"}
	parseResult, err := fp.Parse(query)
	if err != nil {
		t.Error("Error parsing query:", err)
	}
	parseResult = parseResult.Add(params.CourseCode, "value6")
	result := parseResult.String()
	t.Log(result)
	if parseResult.String() != expectedResult {
		t.Error("Expected", expectedResult)
		t.Error("Actual", result)
	}
	conditionsCount := parseResult.ConditionsCount()
	if conditionsCount != 4 {
		t.Error("Unexpected number of conditions:", conditionsCount)
	}
}
