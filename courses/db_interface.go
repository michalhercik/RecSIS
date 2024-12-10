package courses

import (
	"fmt"
	"strings"

	//"log"
)

type DataManager interface {
	Courses(query query) (coursesPage, error)
	AddCourseToBlueprint(user int, code string) ([]Assignment, error)
	RemoveCourseFromBlueprint(user int, code string) error
}

var db DataManager

func SetDataManager(newDB DataManager) {
	db = newDB
}

type coursesPage struct {
	courses    []Course
	startIndex int
	count      int
	total      int
	search     string
	sorted     sortType
}

type query struct {
	user       int
	startIndex int
	maxCount   int
	search     string
	sorted     sortType 
}

const (
	relevance    = iota
	recommended  = iota
	rating       = iota
	mostPopular  = iota
	newest       = iota
)

type sortType int

var sortTypeName = map[sortType]string{
	relevance:   "Relevant",
	recommended: "Recommended",
	rating:      "By rating",
	mostPopular: "Most popular",
	newest:		 "Newest",
}

func (st sortType) String() string {
	return sortTypeName[st]
}

type Teacher struct {
	sisId       int
	firstName   string
	lastName    string
	titleBefore string
	titleAfter  string
}

func (t Teacher) String() string {
	return fmt.Sprintf("%s %s %s %s",
		t.titleBefore, t.firstName, t.lastName, t.titleAfter)
}

type Semester int

const (
	winter Semester = iota + 1
	summer
	both
)

var semesterNameEn = map[Semester]string{
	winter: "Winter",
	summer: "Summer",
	both:   "Both",
}

func (s Semester) String() string {
	return semesterNameEn[s]
}

type Teachers []Teacher

func (t Teachers) string() string {
	names := []string{}
	for _, teacher := range t {
		names = append(names, teacher.String())
	}
	return strings.Join(names, "<br>")
}

func (t *Teachers) trim() {
	result := Teachers{}
	for _, teacher := range *t {
		if teacher.sisId != -1 {
			result = append(result, teacher)
		}
	}
	*t = result
}

type Assignment struct {
	year     int
	semester Semester
}

func (a Assignment) String() string {
	result := fmt.Sprintf("Year %d, semester %s", a.year, a.semester)
	if a.year == 0 {
		result = "Not assigned"
	}
	return result
}

type Course struct {
	code            	 string
	nameCs          	 string
	nameEn          	 string
	start           	 Semester
	semesterCount   	 int
	lectureRange1   	 int
	seminarRange1   	 int
	lectureRange2   	 int
	seminarRange2   	 int
	examType        	 string
	credits         	 int
	teachers        	 Teachers
	rating				 int
	blueprintAssignments []Assignment 
}

func newCourse() *Course {
	return &Course{
		teachers: []Teacher{{}, {}},
	}
}