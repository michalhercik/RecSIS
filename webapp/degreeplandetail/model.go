package degreeplandetail

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

//================================================================================
// Constants
//================================================================================

const dpCode = "dpCode"

const unlimitedYear = 9999

const (
	checkboxName = "selected-courses"
	maxYearParam = "maxYear"
)

//================================================================================
// Data Types and Methods
//================================================================================

type degreePlanPage struct {
	code            string
	title           string
	fieldCode       string
	fieldTitle      string
	validFrom       int
	validTo         int
	isUserPlan      bool
	blocs           []bloc
	recommendedPlan recommendedPlan
	searchEndpoint  string
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

func (dp *degreePlanPage) isValid() bool {
	return dp.validTo == unlimitedYear
}

// TODO: move this to ELT process
type node struct {
	Data struct {
		ID     string `json:"id"`
		Code   string `json:"code"`
		Label  string `json:"label"`
		InPlan string `json:"inPlan"`
	} `json:"data"`
}

type edge struct {
	Data struct {
		ID     string `json:"id"`
		Source string `json:"source"`
		Target string `json:"target"`
		Type   string `json:"type"`
	} `json:"data"`
}

// generate Cytoscape-compatible JSON data for requisites graph
// TODO: save this data to DB during DP page construction to avoid recomputation
func (dp *degreePlanPage) requisitesGraphData() string {
	nodeMap := make(map[string]node)
	edgeMap := make(map[string]edge)
	//edgeID := 0

	// Build nodes and edges from all courses
	for _, bloc := range dp.blocs {
		for _, course := range bloc.courses {
			if foundNode, exists := nodeMap[course.code]; !exists {
				n := node{}
				n.Data.ID = course.code
				n.Data.Code = course.code
				n.Data.Label = fmt.Sprintf("%s - %s", course.code, course.title) // TODO: maybe better label
				n.Data.InPlan = "true"
				nodeMap[course.code] = n
			} else {
				foundNode.Data.InPlan = "true"
				foundNode.Data.Label = fmt.Sprintf("%s - %s", course.code, course.title)
				nodeMap[course.code] = foundNode
			}
			// Add prerequisite nodes and edges
			//createEdges(course.prerequisites, "prerequisite", course.code, &nodeMap, edgeMap, &edgeID)
			// Add corequisite edges
			//createEdges(course.corequisites, "corequisite", course.code, &nodeMap, edgeMap, &edgeID)
			// Add incompatibility edges
			//createEdges(course.incompatibilities, "incompatibility", course.code, &nodeMap, edgeMap, &edgeID)
		}
	}

	// remove incompatibility edges to nodes not in DP
	edgesToRemove := make([]string, 0)
	for edgeKey, e := range edgeMap {
		if e.Data.Type == "incompatibility" {
			if node, exists := nodeMap[e.Data.Target]; !exists || node.Data.InPlan == "false" {
				edgesToRemove = append(edgesToRemove, edgeKey)
			}
		}
	}
	for _, edgeKey := range edgesToRemove {
		delete(edgeMap, edgeKey)
	}

	nodeSlice := make([]node, 0, len(nodeMap))
	for _, n := range nodeMap {
		nodeSlice = append(nodeSlice, n)
	}

	edges := make([]edge, 0, len(edgeMap))
	for _, e := range edgeMap {
		edges = append(edges, e)
	}

	result := struct {
		Nodes []node `json:"nodes"`
		Edges []edge `json:"edges"`
	}{nodeSlice, edges}

	jsonData, _ := json.Marshal(result)
	return string(jsonData)
}

func createEdges(courseCodes []string, edgeType, targetCode string, nodes *map[string]node, edges map[string]edge, edgeID *int) {
	for _, req := range courseCodes {
		if _, exists := (*nodes)[req]; !exists {
			n := node{}
			n.Data.ID = req
			n.Data.Code = req
			n.Data.Label = req
			n.Data.InPlan = "false"
			(*nodes)[req] = n
		}

		// Create a unique key for this edge to prevent duplicates
		edgeKey := fmt.Sprintf("%s->%s:%s", targetCode, req, edgeType)

		// Only add edge if it doesn't already exist
		if _, exists := edges[edgeKey]; !exists {
			e := edge{}
			e.Data.ID = fmt.Sprintf("e%d", *edgeID)
			e.Data.Source = targetCode
			e.Data.Target = req
			e.Data.Type = edgeType
			edges[edgeKey] = e
			(*edgeID)++
		}
	}
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

type recommendedPlan struct {
	years []year
}

func (rp *recommendedPlan) isEmpty() bool {
	return len(rp.years) == 0
}

func (rp *recommendedPlan) totalYears() int {
	return len(rp.years)
}

type year struct {
	winter []course
	summer []course
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
	isSupported        bool
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

func (c *course) creditsString() string {
	if !c.isSupported {
		return ""
	}
	return fmt.Sprintf("%d", c.credits)
}

func (c *course) winterString() string {
	if !c.isSupported {
		return "---"
	}
	winterText := ""
	if c.semester == teachingWinterOnly || c.semester == teachingBoth {
		winterText = fmt.Sprintf("%d/%d, %s", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64, c.examType)
	} else {
		winterText = "---"
	}
	return winterText
}

func (c *course) summerString() string {
	if !c.isSupported {
		return "---"
	}
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

func (c *course) semesterHoursExamString(t text) string {
	if !c.isSupported {
		return "---"
	}
	return fmt.Sprintf("%s %s, %s", c.semester.string(t), c.hoursString(), c.examType)
}

func (c *course) hoursString() string {
	if !c.isSupported {
		return "---"
	}
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
