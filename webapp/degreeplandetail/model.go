package degreeplandetail

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode/utf8"
)

//================================================================================
// Constants
//================================================================================

const dpCode = "dpCode"

const (
	checkboxName = "selected-courses"
)

//================================================================================
// Data Types and Methods
//================================================================================

type degreePlanPage struct {
	degreePlanCode string
	isUserPlan     bool
	blocs          []bloc
}

func (dp *degreePlanPage) bpNumberOfSemesters() int {
	if len(dp.blocs) == 0 {
		return 0
	}
	if len(dp.blocs[0].courses) == 0 {
		return 0
	}
	return len(dp.blocs[0].courses[0].blueprintSemesters)
}

type bloc struct {
	name         string
	code         string
	limit        int
	isCompulsory bool
	isOptional   bool
	courses      []course
}

func (b *bloc) hasLimit() bool {
	return b.limit > -1
}

func (b *bloc) isAssigned() bool {
	if b.hasLimit() && b.assignedCredits() >= b.limit {
		return true
	}
	return false
}

func (b *bloc) assignedCredits() int {
	credits := 0
	for _, c := range b.courses {
		if c.isAssigned() {
			credits += c.credits
		}
	}
	return credits
}

func (b *bloc) isCompleted() bool {
	if b.hasLimit() && b.completedCredits() >= b.limit {
		return true
	}
	return false
}

func (b *bloc) completedCredits() int {
	credits := 0
	for _, c := range b.courses {
		// TODO: add course completion status -> change `false` to `course.Completed`
		if false {
			credits += c.credits
		}
	}
	return credits
}

func (b *bloc) blueprintCredits() int {
	credits := 0
	for _, c := range b.courses {
		if c.isInBlueprint() {
			credits += c.credits
		}
	}
	return credits
}

type course struct {
	code               string
	title              string
	credits            int
	semester           teachingSemester
	guarantors         teacherSlice
	lectureRangeWinter sql.NullInt64
	seminarRangeWinter sql.NullInt64
	lectureRangeSummer sql.NullInt64
	seminarRangeSummer sql.NullInt64
	examType           string
	blueprintSemesters []bool
}

func (c *course) isInBlueprint() bool {
	for _, isIn := range c.blueprintSemesters {
		if isIn {
			return true
		}
	}
	return false
}

func (c *course) isAssigned() bool {
	if len(c.blueprintSemesters) < 2 {
		return false
	}
	for _, isIn := range c.blueprintSemesters[1:] {
		if isIn {
			return true
		}
	}
	return false
}

func (c *course) isUnassigned() bool {
	if len(c.blueprintSemesters) < 1 {
		return false
	}
	return c.blueprintSemesters[0]
}

func (c *course) statusBackgroundColor() string {
	// TODO: add course completion status -> change `false` to `course.Completed`
	if false {
		return "bg-success"
	} else if c.isAssigned() {
		return "bg-blueprint"
	} else if c.isUnassigned() {
		return "bg-blueprint"
	} else {
		return "bg-danger"
	}
}

func (c *course) winterString() string {
	winterText := ""
	if c.semester == teachingWinterOnly || c.semester == teachingBoth {
		winterText = fmt.Sprintf("%d/%d, %s", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64, c.examType)
	} else {
		winterText = "---"
	}
	return winterText
}

func (c *course) summerString() string {
	summerText := ""
	switch c.semester {
	case teachingSummerOnly:
		summerText = fmt.Sprintf("%d/%d, %s", c.lectureRangeSummer.Int64, c.seminarRangeSummer.Int64, c.examType)
	case teachingBoth:
		summerText = fmt.Sprintf("%d/%d, %s", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64, c.examType)
	default:
		summerText = "---"
	}
	return summerText
}

func (c *course) hoursString() string {
	result := ""
	winter := c.lectureRangeWinter.Valid && c.seminarRangeWinter.Valid
	summer := c.lectureRangeSummer.Valid && c.seminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", c.lectureRangeSummer.Int64, c.seminarRangeSummer.Int64)
	}
	return result
}

type teachingSemester int

const (
	teachingWinterOnly teachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts teachingSemester) string(t text) string {
	switch ts {
	case teachingBoth:
		return t.both
	case teachingWinterOnly:
		return t.winter
	case teachingSummerOnly:
		return t.summer
	default:
		return ""
	}
}

type teacherSlice []teacher

func (t teacherSlice) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.string())
	}
	if len(names) == 0 {
		return "---"
	}
	return strings.Join(names, ", ")
}

type teacher struct {
	sisID       string
	lastName    string
	firstName   string
	titleBefore string
	titleAfter  string
}

func (t teacher) string() string {
	firstRune, _ := utf8.DecodeRuneInString(t.firstName)
	return fmt.Sprintf("%c. %s", firstRune, t.lastName)
}
