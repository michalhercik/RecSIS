package degreeplan

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode/utf8"
)

//================================================================================
// Constants
//================================================================================

const (
	searchDegreePlanName  = "search-dp-query"
	searchDegreePlanYear  = "search-dp-year"
	saveDegreePlanYear    = "save-dp-year"
	searchDegreePlanLimit = 5 // TODO: change to a bigger number

	checkboxName = "selected-courses"
)

//================================================================================
// Data Types and Methods
//================================================================================

// all data needed for the degree plan page
type degreePlanPage struct {
	degreePlanCode string
	degreePlanYear int  // year when user started studying
	canSave        bool // if the user can save the degree plan
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

// bloc info with courses
type bloc struct {
	name         string
	code         int
	note         string
	limit        int
	isCompulsory bool
	courses      []course
}

func (b *bloc) hasLimit() bool {
	return b.limit > -1
}

func (b *bloc) isAssigned() bool {
	// ignores courses unassigned in the blueprint
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

func (b *bloc) isInBlueprint() bool {
	if b.hasLimit() && b.blueprintCredits() >= b.limit {
		return true
	}
	return false
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

// course representation
type course struct {
	code               string
	title              string
	note               string
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

// if is (un)assigned in the blueprint
func (c *course) isInBlueprint() bool {
	for _, isIn := range c.blueprintSemesters {
		if isIn {
			return true
		}
	}
	return false
}

// if is assigned in the blueprint
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

// if is unassigned in the blueprint
func (c *course) isUnassigned() bool {
	if len(c.blueprintSemesters) < 1 {
		return false
	}
	return c.blueprintSemesters[0]
}

func (c *course) statusBackgroundColor() string {
	// TODO: add course completion status -> change `false` to `course.Completed`
	if false {
		// most important is completion status
		return "bg-success"
	} else if c.isAssigned() {
		// if not completed, check if assigned
		return "bg-blueprint"
	} else if c.isUnassigned() {
		// if not even assigned, check if unassigned
		return "bg-blueprint" // same color, but the row has a warning icon
	} else {
		// if nothing else, then it is not completed
		return "bg-danger"
	}
}

// semester type - winter, summer, or both
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

// wrapper for teacher slice
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

// teacher type
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

//================================================================================
// Helper Functions
//================================================================================

func winterString(course *course) string {
	winterText := ""
	if course.semester == teachingWinterOnly || course.semester == teachingBoth {
		winterText = fmt.Sprintf("%d/%d, %s", course.lectureRangeWinter.Int64, course.seminarRangeWinter.Int64, course.examType)
	} else {
		winterText = "---"
	}
	return winterText
}

func summerString(course *course) string {
	summerText := ""
	if course.semester == teachingSummerOnly {
		summerText = fmt.Sprintf("%d/%d, %s", course.lectureRangeSummer.Int64, course.seminarRangeSummer.Int64, course.examType)
	} else if course.semester == teachingBoth {
		summerText = fmt.Sprintf("%d/%d, %s", course.lectureRangeWinter.Int64, course.seminarRangeWinter.Int64, course.examType)
	} else {
		summerText = "---"
	}
	return summerText
}

func hoursString(course *course, t text) string {
	result := ""
	winter := course.lectureRangeWinter.Valid && course.seminarRangeWinter.Valid
	summer := course.lectureRangeSummer.Valid && course.seminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", course.lectureRangeWinter.Int64, course.seminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", course.lectureRangeSummer.Int64, course.seminarRangeSummer.Int64)
	}
	return result
}
