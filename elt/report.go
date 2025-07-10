package main

import (
	"fmt"
	"time"
)

type Report struct {
	Failed  int
	Elapsed time.Duration
}

func makeReport(err error, elapsed time.Duration) Report {
	failed := countFailed(err)
	return Report{
		Failed:  failed,
		Elapsed: elapsed,
	}
}

func countFailed(err error) int {
	var failed = 0
	if err != nil {
		switch list := err.(type) {
		case listOfErrors:
			for _, e := range list {
				failed += countFailed(e)
			}
		default:
			failed = 1
		}
	}
	return failed
}

func (r Report) String() string {
	mark := "✅"
	if r.Failed > 0 {
		mark = "❌"
	}
	result := fmt.Sprintf("%s Finished with %d errors in %s\n", mark, r.Failed, r.Elapsed)
	return result
}
