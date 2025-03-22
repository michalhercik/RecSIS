package filter

import (
	"fmt"
	"net/url"
	"testing"
)

func TestQuery(t *testing.T) {
	u, err := url.Parse(fmt.Sprintf("http://localhost:8000/cs/courses?q=bla&par%d=7&par%d=1&par%d=3", Credits, SemesterCount, Credits))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Log(u)
	exp, err := ParseFilters(u.Query())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Log(exp)
	if exp != "credits IN [ocel,3] AND semester_count IN [1]" {
		t.Errorf("unexpected expression: %s", exp)
	}
}
