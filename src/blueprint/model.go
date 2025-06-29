package blueprint

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
	yearUnassign     string = "year-unassign"
	semesterUnassign string = "semester-unassign"
	selectedMove     string = "selected-move"
)

const (
	yearRemove     string = "year-remove"
	semesterRemove string = "semester-remove"
	selectedRemove string = "selected-remove"
)

const (
	yearParam     string = "year"
	semesterParam string = "semester"
	positionParam string = "position"
	typeParam     string = "type"
	unassignParam string = "unassign"
	foldedParam   string = "folded"
)

const (
	checkboxName string = "selected"
)

//================================================================================
// Data Types and Methods
//================================================================================

// all data needed for the blueprint page
type blueprintPage struct {
	unassigned semester
	years      assignedYears
}

func (bp *blueprintPage) totalCredits() int {
	total := bp.unassigned.credits()
	for _, year := range bp.years {
		total += year.credits()
	}
	return total
}

// wrapper for a slice of academic years
type assignedYears []academicYear

func (ays *assignedYears) assignedCredits() int {
	total := 0
	for _, year := range *ays {
		total += year.credits()
	}
	return total
}

// single academic year with its semesters
type academicYear struct {
	position int
	winter   semester
	summer   semester
}

func (ay academicYear) credits() int {
	return ay.winter.credits() + ay.summer.credits()
}

// semester has courses and a flag indicating whether it is folded
type semester struct {
	courses []course
	folded  bool
}

func (s semester) credits() int {
	sum := 0
	for _, course := range s.courses {
		sum += course.credits
	}
	return sum
}

// course representation
type course struct {
	id                 int
	code               string
	title              string
	semester           teachingSemester
	lectureRangeWinter sql.NullInt64
	seminarRangeWinter sql.NullInt64
	lectureRangeSummer sql.NullInt64
	seminarRangeSummer sql.NullInt64
	examType           string
	credits            int
	guarantors         teacherSlice
	warnings           []string
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
	case teachingWinterOnly:
		return t.winter
	case teachingSummerOnly:
		return t.summer
	case teachingBoth:
		return t.both
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
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t teacher) string() string {
	firstRune, _ := utf8.DecodeRuneInString(t.firstName)
	return fmt.Sprintf("%c. %s", firstRune, t.lastName)
}

// semester representation
type semesterAssignment int

const (
	assignmentNone semesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

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

func sortHxPatch(year, semester int, t text) string {
	return fmt.Sprintf(`sortHxPatch($item, %d, %d, $position, "%s")`, year, semester, t.language)
}

// hx-vals generator methods
func vType(typeStr string) string    { return fmt.Sprintf(`"%s": "%s"`, typeParam, typeStr) }
func vYear(year int) string          { return fmt.Sprintf(`"%s": %d`, yearParam, year) }
func vSem(semester int) string       { return fmt.Sprintf(`"%s": %d`, semesterParam, semester) }
func vPos() string                   { return fmt.Sprintf(`"%s": %d`, positionParam, lastPosition) }
func vUnassign(unassign bool) string { return fmt.Sprintf(`"%s": %t`, unassignParam, unassign) }
func vFolded(folded bool) string     { return fmt.Sprintf(`"%s": %t`, foldedParam, folded) }

func mergeVals(vals ...string) string { return strings.Join(vals, ", ") }
