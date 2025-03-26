package filter

import (
	"fmt"
	"net/url"
	"testing"
)

func TestExcept(t *testing.T) {
	query := url.Values{}
	query.Add(fmt.Sprintf("%s%d", ParamPrefix, Credits), "5")
	query.Add(fmt.Sprintf("%s%d", ParamPrefix, SemesterCount), "1")
	query.Add(fmt.Sprintf("%s%d", ParamPrefix, Credits), "3")
	t.Log(query)
	exp, err := ParseFilters(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[Parameter]string{
		Credits:       "semester_count IN [1]",
		SemesterCount: "credits IN [5,3]"}
	i := 0
	for param, filter := range exp.Except() {
		t.Logf("param: %s, filter: %s", param, filter)
		expectedFilter, ok := expected[param]
		if !ok {
			t.Errorf("unexpected param: %s", param)
		}
		if expectedFilter != filter {
			t.Errorf("%s != %s", expectedFilter, filter)
		}
		i += 1
	}
}
