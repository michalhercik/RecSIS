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
	lastPosition   = -1
	unassignedYear = -1
	ttDelay        = "600"
)

const (
	semesterUnassign string = "semester-unassign"
	selectedMove     string = "selected-move"
)

const (
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

const recordID string = "id"

//================================================================================
// Data Types and Methods
//================================================================================

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

type courseLocation struct {
	year     int
	semester semesterAssignment
	course   *course
}

func (bp *blueprintPage) courses() <-chan courseLocation {
	ch := make(chan courseLocation)
	go func() {
		for i := range bp.unassigned.courses {
			ch <- courseLocation{unassignedYear, assignmentNone, &bp.unassigned.courses[i]}
		}
		for _, year := range bp.years {
			for i := range year.winter.courses {
				ch <- courseLocation{year.position, assignmentWinter, &year.winter.courses[i]}
			}
			for i := range year.summer.courses {
				ch <- courseLocation{year.position, assignmentSummer, &year.summer.courses[i]}
			}
		}
		close(ch)
	}()
	return ch
}

type assignedYears []academicYear

func (ays *assignedYears) assignedCredits() int {
	total := 0
	for _, year := range *ays {
		total += year.credits()
	}
	return total
}

type academicYear struct {
	position int
	winter   semester
	summer   semester
}

func (ay academicYear) credits() int {
	return ay.winter.credits() + ay.summer.credits()
}

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
	if c.semester == teachingSummerOnly {
		summerText = fmt.Sprintf("%d/%d, %s", c.lectureRangeSummer.Int64, c.seminarRangeSummer.Int64, c.examType)
	} else if c.semester == teachingBoth {
		summerText = fmt.Sprintf("%d/%d, %s", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64, c.examType)
	} else {
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
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t teacher) string() string {
	firstRune, _ := utf8.DecodeRuneInString(t.firstName)
	return fmt.Sprintf("%c. %s", firstRune, t.lastName)
}

type semesterAssignment int

const (
	assignmentNone semesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

//================================================================================
// Helper Functions
//================================================================================

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
